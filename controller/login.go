package controller

import (
	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func LoginUser(db *gorm.DB, c *fiber.Ctx) error {
	user := new(models.User)
	input := new(models.User)

	if err := c.BodyParser(input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result := db.Where("Username = ?", input.Username).First(user)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	t, err := GenerateAccessToken(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"Message": "success"})
}
