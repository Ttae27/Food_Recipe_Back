package controller

import (
	"time"
	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"fmt"
)

func LoginUser(db *gorm.DB, c *fiber.Ctx) error{
	user := new(models.User)
	input := new(models.User)

	if err:=c.BodyParser(input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result := db.Where("Username = ?",input.Username).First(user)
	if result.Error != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message":"No username"})
	}
	
	fmt.Println(user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(input.Password));err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message":"Password incorrect"})
	}

	t,err := GenerateAccessToken(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name: "jwt",
		Value: t,
		Expires: time.Now().Add(time.Hour*72),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"Message":"success"})
}