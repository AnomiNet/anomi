package model

import (
	"errors"
	"github.com/anominet/anomi/cache"
	"strconv"
	"time"
)

type Post struct {
	Id           int64   `json:"id"`
	CreationTime int64   `json:"created_at"`
	UserHandle   string  `json:"user_handle"`
	UserId       int64   `json:"user_id"`
	ParentId     int64   `json:"parent_id"`
	RootId       int64   `json:"root_id"`
	Url          string  `json:"url"`
	Body         string  `json:"body"`
	Tldr         string  `json:"tldr"`
	Depth        int64   `json:"depth"`
	ChildIds     []int64 `json:"child_ids"` // Kept denormalized in redis
	Score        int64   `json:"score"`
}

type PostScore struct {
	PostId int64 `json:"post_id"`
	Score  int64 `json:"score"`
}

const NEXT_POST_ID_KEY = "counter:next.post.id"
const TOP_POSTS_KEY = "zset:top.posts"

var ErrInvalidParent = errors.New("No such parent post")

func (p *Post) GenerateTldr() {
	if len(p.Body) > 80 {
		p.Tldr = p.Body[:76] + "..."
	} else {
		p.Tldr = p.Body
	}
}

func (e ModelEnv) GetPost(id int64) (*Post, error) {
	p := Post{}
	err := e.C.Get(&p, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	} else {
		return &p, err
	}
}

func (e ModelEnv) GetPostNormalized(id int64) (*Post, error) {
	p := Post{}
	err := e.C.Get(&p, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	p.Score, err = e.C.ZScore(TOP_POSTS_KEY, id)
	if err != nil {
		e.Log.Error(err)
		return nil, err
	}
	p.ChildIds, err = e.GetPostChildIds(id)
	if err != nil {
		return nil, err
	}
	return &p, err
}

func (e ModelEnv) GetPostChildIds(id int64) ([]int64, error) {
	list := []int64{}
	err := e.C.GetList(&list, "child.ids:"+strconv.FormatInt(id, 10))
	return list, err
}

func (e ModelEnv) AppendPostChildId(pid, cid int64) error {
	return e.C.Append("child.ids:"+strconv.FormatInt(pid, 10), cid)
}

func (e ModelEnv) CreatePost(p *Post) error {
	var err error
	if p.Tldr == "" {
		p.GenerateTldr()
	}

	p.CreationTime = time.Now().Unix()

	if p.ParentId != 0 {
		par, err := e.GetPost(p.ParentId)
		if err != nil {
			//if err == redis.ErrNil {
				return ErrInvalidParent
			//} else {
				//return err
			//}
		}
		p.Depth = par.Depth + 1
		p.RootId = par.RootId

		p.Id, err = e.C.Incr(NEXT_POST_ID_KEY)
		if err != nil {
			return err
		}

		err = e.AppendPostChildId(p.ParentId, p.Id)
		if err != nil {
			return err
		}
	} else {
		p.Id, err = e.C.Incr(NEXT_POST_ID_KEY)
		if err != nil {
			return err
		}
		p.Depth = 0
	}

	// FIXME using id as score for testing
	e.C.ZAdd(TOP_POSTS_KEY, p.Id, p.Id)

	if p.ParentId == 0 {
		p.RootId = p.Id
	}

	return e.C.Set(strconv.FormatInt(p.Id, 10), p)
}

func (e ModelEnv) GetTopPosts(limit int64) ([]Post, error) {
	pids := []int64{}
	scores, err := e.C.ZRangeByScore(&pids, TOP_POSTS_KEY, cache.HIGH_TO_LOW, limit)
	if err != nil {
		return nil, err
	}
	posts := make([]Post, len(scores))
	for i := range posts {
		p, err := e.GetPost(pids[i])
		if err != nil {
			return nil, err
		}
		posts[i] = *p
		posts[i].Score = scores[i]
		posts[i].ChildIds, err = e.GetPostChildIds(pids[i])
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (e ModelEnv) GetPostInContext(id int64) ([]Post, error) {
	post, err := e.GetPostNormalized(id)
	if err != nil {
		return nil, err
	}
	posts_len := 1

	if post.ParentId != 0 {
		posts_len += 1
	}

	posts_len += len(post.ChildIds)
	posts := make([]Post, posts_len)


	posts[0] = *post

	i := 1

	if post.ParentId != 0 {
		p, err := e.GetPostNormalized(post.ParentId)
		if err != nil {
			return nil, err
		}
		posts[i] = *p
		i=i+1
	}

	// FIXME recurse
	for j := range post.ChildIds {
		p, err := e.GetPostNormalized(post.ChildIds[j])
		if err != nil {
			return nil, err
		}
		posts[i] = *p
		i=i+1
	}

	return posts, nil
}
