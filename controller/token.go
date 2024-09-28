package controller

import (
	"errors"
	"os"
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateAccessToken(id uint) (string,error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour*72).Unix()

	jwtSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	if jwtSecretKey == "" {
		return "",errors.New("Unauthorized")
	}

	// fmt.Println(jwtSecretKey)
	tokenString,err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		fmt.Println(err)
		return "",errors.New("error signing string")
	}
	return tokenString,nil
}

// func GenerateRefreshToken(id uint) (string,error){
// 	token := jwt.New(jwt.SigningMethodES256)
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["user_id"] = id
// 	claims["exp"] = time.Now().Add(time.Hour*72).Unix()

// 	jwtSecretKey := os.Getenv("ACCESS_SECRET_KEY")
// 	if jwtSecretKey == "" {
// 		fmt.Println("PIKA")
// 		return "",errors.New("Unauthorized")
// 	}

// 	tokenString,err := token.SignedString(jwtSecretKey)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "",errors.New("error signing string")
// 	}
// 	return tokenString,nil
// }