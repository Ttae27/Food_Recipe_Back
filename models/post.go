package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title      string
	Category   string
	Detail     string
	Material   []Material `gorm:"many2many:Post_Material;"`
	Recipe     string
	Picture    string
	Like       []User    `gorm:"many2many:Post_Like;"`
	Comment    []Comment `gorm:"many2many:Post_Comment;"`
	Calories   uint
	Price      uint
	TimeToCook uint
}
