package models

type Bookmark struct {
	UserId uint
	User   User
	PostId uint
	Post   Post
}
