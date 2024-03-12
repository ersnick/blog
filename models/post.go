package models

type Post struct {
	Id                   int64
	Body                 string
	Username             string
	Likes                int64
	IsLikedByCurrentUser bool
	Editable             bool
	Deletable            bool
}
