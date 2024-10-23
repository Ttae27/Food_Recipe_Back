package controller

import (
	"log"
	"strconv"

	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB,c *fiber.Ctx) error {
	input := new(models.User)

	if err:= c.BodyParser(input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var existingUser models.User
	if err := db.Where("username ILIKE ?", input.Username).First(&existingUser).Error; err == nil {
		// If no error, that means a user with this username already exists
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message":"username already exist"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password),bcrypt.DefaultCost)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	input.Password = string(hashedPassword)

	result := db.Create(input)

	if result.Error != nil{
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(input)
}

func GetUser(db *gorm.DB, c *fiber.Ctx) error {
	var user models.User
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result := db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
            // Return 404 if user is not found
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "User not found",
            })
        }
        // Return 500 if there is some other database error
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to retrieve user",
        })
	}

	return c.JSON(user)
}

func UpdateUser(db *gorm.DB,c *fiber.Ctx) error {
	user := new(models.User)
	updatedUser := new(models.User)

	id,err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message":err})
	}
	
	result := db.First(&user,id)
	if result.Error != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message":result.Error})
	}

	if err := c.BodyParser(updatedUser);err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message":err})
	}

	var existingUser models.User
	if err := db.Where("username ILIKE ?", updatedUser.Username).First(&existingUser).Error; err == nil {
		// If no error, that means a user with this username already exists
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message":"username already exist"})
	}

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password),bcrypt.DefaultCost)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		updatedUser.Password = string(hashedPassword)
	} 

	user.ID = uint(id)

	updateResult := db.Model(&user).Updates(updatedUser)
	if updateResult.Error != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message":result.Error})
	}
	return c.JSON(updatedUser)
}

func DeleteUser(db *gorm.DB,c *fiber.Ctx) error {
	var user models.User;
	id,err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result := db.Delete(&user,id)
	if result.Error != nil{
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON("Deleted Successfully")
}

