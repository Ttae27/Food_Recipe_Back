package controller

import (
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
	i, _ = strconv.ParseUint(form.Value["category"][0], 10, 64)
	post.CategoryID = uint(i)

	// image := form.File["image"][0]
	// destination := fmt.Sprintf("./uploads/%s", image.Filename)
	// if err := c.SaveFile(image, destination); err != nil {
	// 	return err
	// }
	// post.Picture = destination

	post.Picture = form.Value["image"][0]

	//uid
	userId, err := GetUserId(c)
	post.UserID = userId
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

	// Calculate total calories for the post
	totalCalories, err := CalculateCalories(db, ingredientSlice, quantitySlice)
	if err != nil {
		log.Fatal("Error calculating calories: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	post.Calories = uint(totalCalories)

	// Calculate total price for the post
	totalPrice, err := CalculatePrice(db, ingredientSlice, quantitySlice)
	if err != nil {
		log.Fatal("Error calculating price: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	post.Price = uint(totalPrice)

	// Save updated post with the calculated calories and price
	result = db.Save(post)
	if result.Error != nil {
		log.Fatal("Error saving post with calories and price: ", result.Error)
	}

	postIn := new(models.Post_Ingredient)
	for index, ingredients := range ingredientSlice {
		inWithQuan := new(models.IngredientWithQuantity)
		postIn.PostID = post.ID

		i, _ := strconv.ParseUint(ingredients, 10, 64)
		inWithQuan.IngredientID = uint(i)
		i, _ = strconv.ParseUint(quantitySlice[index], 10, 64)
		inWithQuan.Quantity = uint(i)
		result := db.Create(inWithQuan)
		db.Last(&inWithQuan)
		postIn.IngredientWithQuantityID = inWithQuan.ID
		result = db.Create(postIn)
		if result.Error != nil {
			log.Fatal("Error insert ingredient to post: ", result.Error)
		}
	}

	return c.SendString("Create Post Successful")
}

func GetPost(db *gorm.DB, c *fiber.Ctx) error {
	var post models.Post
	postId := c.Params("id")

	result := db.Preload("Ingredients").Preload("Ingredients.IngredientWithQuantity.Ingredient").Preload("Like").Preload("PostComments").Preload("PostComments.Comment").Preload("PostComments.User").Preload("Category").Preload("User").First(&post, postId)
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

func GetsPostWithUser(db *gorm.DB, c *fiber.Ctx) error {
	var posts []models.Post

	uid, _ := GetUserId(c)
	result := db.Where("user_id = ?", uid).Find(&posts)
	if result.Error != nil {
		log.Fatal("Error getting post: ", result.Error)
	}

	return c.JSON(posts)
}

func DeletePost(db *gorm.DB, c *fiber.Ctx) error {
	var post models.Post
	uidPost := new(models.Post)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	result := db.First(&uidPost, id)
	//check user
	if uid, _ := GetUserId(c); uid != uidPost.UserID {
		return c.SendString("User not writer of this post!")
	}

	result = db.Delete(&post, id)

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

	//check user
	if uid, _ := GetUserId(c); uid != oldPost.UserID {
		return c.SendString("User not writer of this post!")
	}

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

	postIn := new(models.Post_Ingredient)
	for index, ingredients := range ingredientSlice {
		inWithQuan := new(models.IngredientWithQuantity)
		postIn.PostID = uint(id)
		i, _ := strconv.ParseUint(ingredients, 10, 64)
		inWithQuan.IngredientID = uint(i)
		i, _ = strconv.ParseUint(quantitySlice[index], 10, 64)
		inWithQuan.Quantity = uint(i)
		result := db.Create(inWithQuan)
		db.Last(&inWithQuan)
		postIn.IngredientWithQuantityID = inWithQuan.ID
		result = db.Create(postIn)
		if result.Error != nil {
			log.Fatal("Error insert ingredient to post: ", result.Error)
		}
	}

	// Calculate total Calories
	totalCalories, err := CalculateCalories(db, ingredientSlice, quantitySlice)
	if err != nil {
		log.Println("Error in calculating calories:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error calculating calories")
	}

	// Calculate total Price
	totalPrice, err := CalculatePrice(db, ingredientSlice, quantitySlice)
	if err != nil {
		log.Println("Error in calculating price:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error calculating price")
	}

	oldPost.Calories = uint(totalCalories)
	oldPost.Price = uint(totalPrice)

	result = db.Save(oldPost)
	if result.Error != nil {
		log.Println("Error saving updated post:", result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving updated post")
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
	user_comment.UserID, _ = GetUserId(c)
	result = db.Create(user_comment)
	if result.Error != nil {
		log.Fatal("Error creating comment: ", result.Error)
	}

	post_comment := new(models.Post_Comment)
	post_comment.CommentID = comment.ID
	i, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	post_comment.PostID = uint(i)
	post_comment.UserID, _ = GetUserId(c)
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
	uid, _ := GetUserId(c)
	pid, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	id, err := strconv.ParseUint(form.Value["commentid"][0], 10, 64)
	//check writer , check comment's user
	post := new(models.Post)
	result := db.First(&post, pid)
	userComment := new(models.User_Comment)
	result = db.Where("comment_id = ?", id).First(&userComment)
	if userComment.UserID != uid && post.UserID != uid {
		return c.SendString("User not allow to delete comment")
	}

	if err != nil {
		return c.SendString("delete comment: can't convert commentid")
	}
	result = db.Delete(&comment, id)
	if result.Error != nil {
		log.Fatal("Error delete comment: ", result.Error)
	}

	result = db.Where("post_id = ? and comment_id = ?", pid, id).Delete(&models.Post_Comment{})
	if result.Error != nil {
		log.Fatal("Error delete comment: ", result.Error)
	}

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
	post_like.UserID, _ = GetUserId(c)
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
	uid, _ := GetUserId(c)

	result := db.Where("post_id = ? and user_id = ?", pid, uid).Delete(&models.Post_Like{})
	if result.Error != nil {
		log.Fatal("Error delete like: ", result.Error)
	}

	return c.JSON("delete like successful")
}

func GetsLike(db *gorm.DB, c *fiber.Ctx) error {
	var likes []models.Post_Like

	uid, _ := GetUserId(c)
	result := db.Where("user_id = ?", uid).Preload("Post").Find(&likes)
	if result.Error != nil {
		log.Fatal("Error getting like: ", result.Error)
	}

	return c.JSON(likes)
}

func AddBookmark(db *gorm.DB, c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	bookmark := new(models.Bookmark)
	i, _ := strconv.ParseUint(form.Value["postid"][0], 10, 64)
	bookmark.PostID = uint(i)
	bookmark.UserID, _ = GetUserId(c)
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
	uid, _ := GetUserId(c)
	result := db.Where("post_id = ? and user_id = ?", pid, uid).Delete(&models.Bookmark{})
	if result.Error != nil {
		log.Fatal("Error delete bookmark: ", result.Error)
	}

	return c.JSON("delete bookmark successful")
}

func GetsBookmark(db *gorm.DB, c *fiber.Ctx) error {
	var bookmarks []models.Bookmark

	uid, _ := GetUserId(c)
	result := db.Where("user_id = ?", uid).Preload("Post").Preload("Post.Category").Find(&bookmarks)
	if result.Error != nil {
		log.Fatal("Error getting post: ", result.Error)
	}

	return c.JSON(bookmarks)
}

func GetsIngredient(db *gorm.DB, c *fiber.Ctx) error {
	var ingredient models.IngredientCategory
	categoryId := c.Params("id")

	result := db.Preload("Ingredients").Preload("Ingredients.Ingredient").First(&ingredient, categoryId)
	if result.Error != nil {
		log.Fatal("Error getting ingredients: ", result.Error)
	}

	return c.JSON(ingredient)
}

func GetsAllIngredient(db *gorm.DB, c *fiber.Ctx) error {
	var ingredients []models.Ingredient

	result := db.Find(&ingredients)
	if result.Error != nil {
		log.Fatal("Error getting ingredients: ", result.Error)
	}

	return c.JSON(ingredients)
}
