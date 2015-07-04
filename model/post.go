package model

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
	ChildIds     []int64 `json:"child_ids"`
}

type PostScore struct {
	PostId int64 `json:"post_id"`
	Score  int64 `json:"score"`
}
