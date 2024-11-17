package models

import (
	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name     string
	Calories float64
	Price    float64
	Picture  string
	Unit     string
	Posts    []Post `gorm:"many2many:Post_Ingredient;"`
}

type IngredientCategory struct {
	gorm.Model
	Type        string
	Ingredients []Ingredient_IngredientCategory
}

type Ingredient_IngredientCategory struct {
	IngredientCategoryID uint
	IngredientID         uint
	Ingredient           Ingredient
}

type IngredientWithQuantity struct {
	gorm.Model
	IngredientID uint
	Ingredient   Ingredient
	Quantity     uint
}
