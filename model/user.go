package model

type User struct {
	Id          int64   `json:"id"`
	LastActive  int64   `json:"last_active_at"`
	Handle      string  `json:"handle"`
	Token       string  `json:"token"`
	PostIds     []int64 `json"post_ids"`
	VotePostIds []int64 `json"vote_post_ids"`
}

type Vote struct {
	PostId int64 `json:"post_id"`
	UserId int64 `json:"user_id"`
	Vector int8  `json:"vector"`
}
