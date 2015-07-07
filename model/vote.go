package model

import (
	"errors"
	"strconv"
)

type Vote struct {
	PostId     int64  `json:"post_id"`
	UserHandle string `json:"user_handle"`
	Value      int8   `json:"vector"`
}

var ErrInvalidPost = errors.New("An invalid post id was specified")

func (e ModelEnv) CreateVote(v *Vote) error {
	_, err := e.GetPost(v.PostId)
	if err != nil {
		return ErrInvalidPost
	}

	delta := v.Value

	existing_vote, err := e.GetVoteByUserHandle(v.PostId, v.UserHandle)
	if err == nil {
		// Have existing vote
		delta = v.Value - existing_vote
	}
	e.C.ZIncrBy(TOP_POSTS_KEY, int64(delta), v.PostId)
	return e.C.Set(v.UserHandle+":"+strconv.FormatInt(v.PostId, 10), v.Value)
}

func (e ModelEnv) GetVoteByToken(pid int64, tok string) (int8, error) {
	active_user, err := e.GetActiveUser(tok)
	if err != nil {
		return 0, err
	}
	return e.GetVoteByUserHandle(pid, active_user)
}
func (e ModelEnv) GetVoteByUserHandle(pid int64, user string) (int8, error) {
	// FIXME hacks
	var v int8
	//FIXME seperator
	err := e.C.Get(&v, user+":"+strconv.FormatInt(pid, 10))
	return v, err
}

func (e ModelEnv) PopulateUserVote(p *Post, tok string) error {
	v, err := e.GetVoteByToken(p.Id, tok)
	if err != nil {
		return nil
	}
	p.CurrentUserVote = v
	return nil
}
