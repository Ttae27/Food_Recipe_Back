package controller

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func LogoutUser(c *fiber.Ctx) error {
	intId, err := strconv.Atoi(c.Params("id"))
	// fmt.Println(c)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id := uint(intId)

	jwtSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	if jwtSecretKey == "" {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// fmt.Println(jwtSecretKey)
	cookie := c.Cookies("jwt")
	fmt.Println(cookie)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	// fmt.Println(token.Valid)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "token error"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token claims"})
	}
	userID := uint(userIDFloat)
	fmt.Println(id, userID)
	if id != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Message": "user error"})
	}

	c.ClearCookie("jwt")
	return c.JSON(fiber.Map{"Message": "Logout Successfully"})

}

func GetUserId(c *fiber.Ctx) (uint, error) {
	jwtSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	cookie := c.Cookies("jwt")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, c.SendStatus(fiber.StatusUnauthorized)
	}

	userIDFloat := claims["user_id"].(float64)
	userID := uint(userIDFloat)

	return userID, nil
}
