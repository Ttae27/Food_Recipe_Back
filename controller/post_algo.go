package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SearchPostByName searches for posts by their title in any language.
func SearchPostByName(db *gorm.DB, c *fiber.Ctx) error {
	searchName := c.Query("title")

	if searchName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Search query not provided")
	}

	var posts []models.Post

	result := db.Where("title LIKE ?", "%"+searchName+"%").
		Preload("Ingredients").
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

// Calculate Calories
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

// Calculate Price
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

// Price Range
func GetPostsByPriceRange(db *gorm.DB, c *fiber.Ctx) error {
	minPrice := c.Query("min_price", "0")
	maxPrice := c.Query("max_price", "1000000") // Set a high default max price

	min, err := strconv.ParseUint(minPrice, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid min_price"})
	}

	max, err := strconv.ParseUint(maxPrice, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid max_price"})
	}

	var posts []models.Post
	result := db.Where("price BETWEEN ? AND ?", min, max).Find(&posts)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
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

// Category Type 
func GetPostsByCategoryType(db *gorm.DB, c *fiber.Ctx) error {
	categoryType := c.Query("type")

	if categoryType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category type is required",
		})
	}

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

	var posts []models.Post
	if err := db.Where("category_id = ?", category.ID).
		Preload("Category").
		Preload("Ingredients").
		Preload("Like").
		Preload("PostComments").
		Preload("Bookmarks").
		Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(posts)
}
