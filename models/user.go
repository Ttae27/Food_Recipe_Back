package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	Password       string
	FirstName      string
	LastName       string
	PictureProfile string
	Bookmarks      []Post `gorm:"many2many:Bookmark;"`
	Likes          []Post `gorm:"many2many:Post_Like;"`
	Comments       []Comment
}

type User_Comment struct {
	UserID    uint
	CommentID uint
	Comment   Comment
}

