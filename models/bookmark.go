package models

type Bookmark struct {
	UserID uint
	User   User
	PostID uint
	Post   Post
}
