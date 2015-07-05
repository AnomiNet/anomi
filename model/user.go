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

type Vote struct {
	PostId int64 `json:"post_id"`
	UserId int64 `json:"user_id"`
	Vector int8  `json:"vector"`
}

const TOKEN_LEN = 16
const NEXT_USER_ID_KEY = "counter:next.user.id"

var ErrUserExists = errors.New("anomi/model: user already exists")
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken() (string, error) {
	s := make([]rune, TOKEN_LEN)
	max := big.NewInt(int64(len(s)))
	for i := range s {
		j, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		s[i] = letters[j.Int64()]
	}
	return string(s), nil
}

func (e ModelEnv) GetUserByHandle(handle string) *User {
	u := User{}
	err := e.C.Get(&u, handle)
	if err != nil {
		return nil
	} else {
		return &u
	}
}

func (e ModelEnv) CreateUser(u *User) error {
	if e.GetUserByHandle(u.Handle) != nil {
		return ErrUserExists
	}
	var err error
	u.Id, err = e.C.Incr(NEXT_USER_ID_KEY)
	if err != nil {
		return err
	}
	u.Token, err = GenerateToken()
	if err != nil {
		return err
	}
	u.Touch()
	u.PostIds = make([]int64, 0)
	u.VotePostIds = make([]int64, 0)

	return e.C.Set(u.Handle, u)
}

func (u *User) Touch() {
	u.LastActive = time.Now().Unix()
}
