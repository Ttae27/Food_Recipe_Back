package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title        string
	Detail       string
	UserID       uint
	User         User
	CategoryID   uint
	Category     Category
	Ingredients  []Post_Ingredient
	Recipe       string
	Picture      string
	Like         []User `gorm:"many2many:Post_Like;"`
	PostComments []Post_Comment
	Calories     uint
	Price        uint
	TimeToCook   uint
	Bookmarks    []User `gorm:"many2many:Bookmark;"`
}

type Category struct {
	gorm.Model
	Type           string
	PostCategories []Post_Category
}

type Post_Category struct {
	CategoryID uint
	PostID     uint
	Post       Post
}
