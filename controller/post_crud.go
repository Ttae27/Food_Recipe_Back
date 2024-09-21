package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	post := new(models.Post)
	//add based post
	err = json.Unmarshal([]byte(form.Value["post"][0]), &post)
	if err != nil {
		return err
	}

	image := form.File["image"][0]
	destination := fmt.Sprintf("./uploads/%s", image.Filename)
	if err := c.SaveFile(image, destination); err != nil {
		return err
	}
	post.Picture = destination

	result := db.Create(post)
	if result.Error != nil {
		log.Fatal("Error creating post: %v", result.Error)
	}

	//add category
	postCat := new(models.Post_Category)
	postCat.PostID = post.ID
	i, _ := strconv.ParseUint(form.Value["category"][0], 10, 64)
	postCat.CategoryID = uint(i)
	result = db.Create(postCat)
	if result.Error != nil {
		log.Fatal("Error insert ingredient to post: %v", result.Error)
	}

	//manage form to add to table Post_ingredient
	ingredient := strings.Trim(form.Value["ingredient"][0], "[]")
	ingredientSlice := strings.Split(ingredient, ",")
	quantity := strings.Trim(form.Value["quantity"][0], "[]")
	quantitySlice := strings.Split(quantity, ",")

	for index, ingredients := range ingredientSlice {
		postIn := new(models.Post_Ingredient)
		postIn.PostID = post.ID
		i, _ := strconv.ParseUint(ingredients, 10, 64)
		postIn.IngredientID = uint(i)
		i, _ = strconv.ParseUint(quantitySlice[index], 10, 64)
		postIn.Quantity = uint(i)
		result := db.Create(postIn)
		if result.Error != nil {
			log.Fatal("Error insert ingredient to post: %v", result.Error)
		}
	}

	return c.SendString("Create Post Successful")
}

func GetPost(db *gorm.DB, c *fiber.Ctx) error {
	var post models.Post
	postId := c.Params("id")

	result := db.Preload("Ingredients").Preload("PostComments").Preload("Like").First(&post, postId)
	if result.Error != nil {
		log.Fatal("Error getting post: %v", result.Error)
	}

	return c.JSON(post)
}

func GetsPost(db *gorm.DB, c *fiber.Ctx) error {
	var posts []models.Post

	result := db.Preload("Ingredients").Preload("PostComments").Preload("Like").Find(&posts)
	if result.Error != nil {
		log.Fatal("Error getting post: %v", result.Error)
	}

	return c.JSON(posts)
}

func DeletePost(db *gorm.DB, c *fiber.Ctx) error {
	var post models.Post
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	result := db.Delete(&post, id)

	if result.Error != nil {
		log.Fatal("Error delete post: %v", result.Error)
	}

	return c.SendString("Delete successfully")
}
