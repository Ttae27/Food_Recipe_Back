package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SearchPostByNameAndFilters searches for posts by their title and applies additional filters.
func SearchPostByNameAndFilters(db *gorm.DB, c *fiber.Ctx) error {
	searchName := c.Query("title")
	minPrice := c.Query("min_price", "0")
	maxPrice := c.Query("max_price", "1000000") // Set a high default max price
	categoryType := c.Query("type")

	// Parse price range
	min, err := strconv.ParseUint(minPrice, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid min_price"})
	}

	max, err := strconv.ParseUint(maxPrice, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid max_price"})
	}

	// Normalize category type
	var categoryID uint
	if categoryType != "" {
		normalizedCategoryType := capitalizeCase(categoryType)
		var category models.Category
		if err := db.Where("type = ?", normalizedCategoryType).First(&category).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Category not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve category",
			})
		}
		categoryID = category.ID
	}

	var posts []models.Post
	query := db.Model(&models.Post{})

	// Apply title filter if provided
	if searchName != "" {
		query = query.Where("title LIKE ?", "%"+searchName+"%")
	}

	// Apply price range filter
	query = query.Where("price BETWEEN ? AND ?", min, max)

	// Apply category filter if provided
	if categoryType != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	result := query.Preload("Ingredients").
		Preload("PostComments").
		Preload("Like").
		Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error searching for posts: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("No posts found matching the search.")
	}

	return c.JSON(posts)
}

// capitalizeCase converts the first letter to uppercase and the rest to lowercase.
func capitalizeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// CalculateCalories calculates the total calories for a list of ingredients and their quantities.
func CalculateCalories(db *gorm.DB, ingredientSlice []string, quantitySlice []string) (float64, error) {
	var totalCalories float64

	for index, ingredientIDStr := range ingredientSlice {
		ingredientID, err := strconv.ParseUint(ingredientIDStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid ingredient ID: %s", ingredientIDStr)
		}
		quantityForPost, err := strconv.ParseUint(quantitySlice[index], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid quantity: %s", quantitySlice[index])
		}

		var ingredient models.Ingredient
		if err := db.First(&ingredient, ingredientID).Error; err != nil {
			return 0, fmt.Errorf("error fetching ingredient with ID: %d", ingredientID)
		}

		caloriesForIngredient := ingredient.Calories * float64(quantityForPost)
		totalCalories += caloriesForIngredient
	}

	return totalCalories, nil
}

// CalculatePrice calculates the total price for a list of ingredients and their quantities.
func CalculatePrice(db *gorm.DB, ingredientSlice []string, quantitySlice []string) (float64, error) {
	var totalPrice float64

	for index, ingredientIDStr := range ingredientSlice {
		ingredientID, err := strconv.ParseUint(ingredientIDStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid ingredient ID: %s", ingredientIDStr)
		}
		quantityForPost, err := strconv.ParseUint(quantitySlice[index], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid quantity: %s", quantitySlice[index])
		}

		var ingredient models.Ingredient
		if err := db.First(&ingredient, ingredientID).Error; err != nil {
			return 0, fmt.Errorf("error fetching ingredient with ID: %d", ingredientID)
		}

		priceForIngredient := ingredient.Price * float64(quantityForPost)
		totalPrice += priceForIngredient
	}

	return totalPrice, nil
}
