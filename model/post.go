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

func (e ModelEnv) GetPostChildIds(id int64) {
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
			return ErrInvalidParent
		}
		p.Depth = par.Depth + 1

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
		err := e.C.Get(&posts[i], strconv.FormatInt(pids[i], 10))
		if err != nil {
			return nil, err
		}
		posts[i].Score = scores[i]
	}
	return posts, nil
}
