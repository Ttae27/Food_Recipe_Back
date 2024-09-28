package routes

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	"os"
)

func requireAuth(c *fiber.Ctx) error{
	jwtSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	if jwtSecretKey == "" {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	
	cookie := c.Cookies("jwt")
		
	token,err := jwt.ParseWithClaims(cookie,&jwt.RegisteredClaims{},func(token *jwt.Token)(interface{},error){
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid{
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}