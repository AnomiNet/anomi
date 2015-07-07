package model

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

type User struct {
	Id          int64   `json:"id"`
	LastActive  int64   `json:"last_active_at"`
	Handle      string  `json:"handle"`
	Token       string  `json:"token"`
	PostIds     []int64 `json:"post_ids"`
	VotePostIds []int64 `json:"vote_post_ids"`
}

//FIXME seperator
const TOKEN_LEN = 16
const NEXT_USER_ID_KEY = "counter:next.user.id"
const ACTIVE_USERS_KEY = "active.users"

var ErrUserExists = errors.New("User with this handle already exists")
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (u *User) GenerateToken() {
	s := make([]rune, TOKEN_LEN)
	max := big.NewInt(int64(len(letters)))
	for i := range s {
		j, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("Couldn't generate token")
		}
		s[i] = letters[j.Int64()]
	}
	u.Token = string(s)
}

func (e ModelEnv) GetUserByHandle(handle string) (*User, error) {
	u := User{}
	err := e.C.Get(&u, handle)
	if err != nil {
		return nil, err
	} else {
		return &u, err
	}
}

func (e ModelEnv) CreateUser(u *User) error {
	if ok, _ := e.GetUserByHandle(u.Handle); ok != nil {
		return ErrUserExists
	}
	var err error
	u.Id, err = e.C.Incr(NEXT_USER_ID_KEY)
	if err != nil {
		return err
	}
	u.GenerateToken()
	u.Touch()
	u.PostIds = make([]int64, 0)
	u.VotePostIds = make([]int64, 0)

	return e.C.Set(u.Handle, u)
}

func (e ModelEnv) SetActiveUser(u *User) error {
	return e.C.Set(ACTIVE_USERS_KEY+"."+u.Token, u.Handle)
}

func (e ModelEnv) GetActiveUser(tok string) (string, error) {
	var current_user string
	err := e.C.Get(&current_user, ACTIVE_USERS_KEY+"."+tok)
	return current_user, err
}

func (u *User) Touch() {
	u.LastActive = time.Now().Unix()
}
