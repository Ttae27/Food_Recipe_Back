package controller

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func LogoutUser(c *fiber.Ctx) error {
	intId,err:= strconv.Atoi(c.Params("id"))
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

	claims := jwt.MapClaims{}
	token,err := jwt.ParseWithClaims(cookie,claims,func(token *jwt.Token)(interface{},error){
		return []byte(jwtSecretKey), nil
	})
	// fmt.Println(token.Valid)
	if err != nil || !token.Valid{
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userIDFloat := claims["user_id"].(float64)
	userID := uint(userIDFloat)
	fmt.Println(id,userID)
	if id != userID{
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.ClearCookie("jwt")
	return c.JSON(fiber.Map{"Message":"Logout Successfully"})

}