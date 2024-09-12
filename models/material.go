package models

import (
	"gorm.io/gorm"
)

type Material struct {
	gorm.Model
	Name     string
	Calories uint
	Price    uint
	Picture  string
}
