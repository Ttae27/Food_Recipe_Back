package controller

import (
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
	post.Title = form.Value["title"][0]
	post.Detail = form.Value["detail"][0]
	post.Recipe = form.Value["recipe"][0]
	i, _ := strconv.ParseUint(form.Value["timetocook"][0], 10, 64)
	post.TimeToCook = uint(i)

	image := form.File["image"][0]
	destination := fmt.Sprintf("./uploads/%s", image.Filename)
	if err := c.SaveFile(image, destination); err != nil {
		return err
	}
	post.Picture = destination

	result := db.Create(post)
	if result.Error != nil {
		log.Fatal("Error creating post: ", result.Error)
	}

	//add category
	postCat := new(models.Post_Category)
	postCat.PostID = post.ID
	i, _ = strconv.ParseUint(form.Value["category"][0], 10, 64)
	postCat.CategoryID = uint(i)
	result = db.Create(postCat)
	if result.Error != nil {
		log.Fatal("Error insert ingredient to post: ", result.Error)
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
			log.Fatal("Error insert ingredient to post: ", result.Error)
		}
	}

	return c.SendString("Create Post Successful")
}

func GetPost(db *gorm.DB, c *fiber.Ctx) error {
	var post models.Post
	postId := c.Params("id")

	result := db.Preload("Ingredients").Preload("Like").Preload("PostComments").Preload("PostComments.Comment").First(&post, postId)
	if result.Error != nil {
		log.Fatal("Error getting post: ", result.Error)
	}

	return c.JSON(post)
}

func GetsPost(db *gorm.DB, c *fiber.Ctx) error {
	var posts []models.Post

	result := db.Preload("Ingredients").Preload("PostComments").Preload("Like").Find(&posts)
	if result.Error != nil {
		log.Fatal("Error getting post: ", result.Error)
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
		log.Fatal("Error delete post: ", result.Error)
	}

	return c.SendString("Delete successfully")
}

func UpdatePost(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	id, _ := strconv.Atoi(c.Params("id"))
	oldPost := new(models.Post)
	newPost := new(models.Post)

	result := db.First(&oldPost, id)
	if result.Error != nil {
		log.Fatal("Error getting post to update: ", result.Error)
	}
	newPost.Title = form.Value["title"][0]
	newPost.Detail = form.Value["detail"][0]
	newPost.Recipe = form.Value["recipe"][0]
	i, _ := strconv.ParseUint(form.Value["timetocook"][0], 10, 64)
	newPost.TimeToCook = uint(i)
	db.Model(&oldPost).Association("Ingredients").Clear() // delete association
	db.Model(&oldPost).Updates(newPost)

	ingredient := strings.Trim(form.Value["ingredient"][0], "[]")
	ingredientSlice := strings.Split(ingredient, ",")
	quantity := strings.Trim(form.Value["quantity"][0], "[]")
	quantitySlice := strings.Split(quantity, ",")

	for index, ingredients := range ingredientSlice {
		postIn := new(models.Post_Ingredient)
		postIn.PostID = uint(id)
		i, _ := strconv.ParseUint(ingredients, 10, 64)
		postIn.IngredientID = uint(i)
		i, _ = strconv.ParseUint(quantitySlice[index], 10, 64)
		postIn.Quantity = uint(i)
		result := db.Create(postIn)
		if result.Error != nil {
			log.Fatal("Error insert(update) ingredient to post: ", result.Error)
		}
	}

	return c.SendString("update post successful")
}

func AddComment(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	comment := new(models.Comment)
	comment.Comment = form.Value["comment"][0]
	result := db.Create(comment)
	if result.Error != nil {
		log.Fatal("Error creating comment: ", result.Error)
	}

	user_comment := new(models.User_Comment)
	user_comment.CommentID = comment.ID
	i, _ := strconv.ParseUint(form.Value["token"][0], 10, 64)
	user_comment.UserID = uint(i)
	result = db.Create(user_comment)
	if result.Error != nil {
		log.Fatal("Error creating comment: ", result.Error)
	}

	post_comment := new(models.Post_Comment)
	post_comment.CommentID = comment.ID
	i, _ = strconv.ParseUint(form.Value["postid"][0], 10, 64)
	post_comment.PostID = uint(i)
	result = db.Create(post_comment)
	if result.Error != nil {
		log.Fatal("Error creating comment: ", result.Error)
	}

	return c.SendString("Add comment successful")
}

func DeleteComment(db *gorm.DB, c *fiber.Ctx) error {
	var comment models.Comment
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	id, err := strconv.ParseUint(form.Value["commentid"][0], 10, 64)
	if err != nil {
		return c.SendString("193")
	}
	result := db.Delete(&comment, id)
	if result.Error != nil {
		log.Fatal("Error delete comment: ", result.Error)
	}

	cid, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	result = db.Where("post_id = ? and comment_id = ?", cid, id).Delete(&models.Post_Comment{})
	if result.Error != nil {
		log.Fatal("Error delete comment: ", result.Error)
	}

	uid, _ := strconv.ParseUint(form.Value["token"][0], 10, 64)
	result = db.Where("user_id = ? and comment_id = ?", uid, id).Delete(&models.User_Comment{})
	if result.Error != nil {
		log.Fatal("Error delete comment: ", result.Error)
	}

	return c.SendString("delete comment successful")
}

func AddLike(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	post_like := new(models.Post_Like)
	i, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	post_like.PostID = uint(i)
	i, _ = strconv.ParseUint(form.Value["token"][0], 10, 64)
	post_like.UserID = uint(i)
	result := db.Create(post_like)
	if result.Error != nil {
		log.Fatal("Error add Like: ", result.Error)
	}

	return c.JSON("add like successful")
}

func DeleteLike(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	pid, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	uid, _ := strconv.ParseUint(form.Value["token"][0], 10, 64)
	result := db.Where("post_id = ? and user_id = ?", pid, uid).Delete(&models.Post_Like{})
	if result.Error != nil {
		log.Fatal("Error delete like: ", result.Error)
	}

	return c.JSON("delete like successful")
}

func AddBookmark(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	bookmark := new(models.Bookmark)
	i, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	bookmark.PostID = uint(i)
	i, _ = strconv.ParseUint(form.Value["token"][0], 10, 64)
	bookmark.UserID = uint(i)
	result := db.Create(bookmark)
	if result.Error != nil {
		log.Fatal("Error add bookmark: ", result.Error)
	}

	return c.JSON("add bookmark successful")
}

func DeleteBookmark(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	pid, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	uid, _ := strconv.ParseUint(form.Value["token"][0], 10, 64)
	result := db.Where("post_id = ? and user_id = ?", pid, uid).Delete(&models.Bookmark{})
	if result.Error != nil {
		log.Fatal("Error delete bookmark: ", result.Error)
	}

	return c.JSON("delete bookmark successful")
}

func GetsIngredient(db *gorm.DB, c *fiber.Ctx) error {
	var ingredients []models.IngredientCategory

	result := db.Preload("Ingredients").Preload("Ingredients.Ingredient").Find(&ingredients)
	if result.Error != nil {
		log.Fatal("Error getting ingredients: ", result.Error)
	}

	return c.JSON(ingredients)
}
