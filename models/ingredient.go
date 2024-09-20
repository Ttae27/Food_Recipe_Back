package models

import (
	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name     string
	Calories uint
	Price    uint
	Picture  string
	Unit     string
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
