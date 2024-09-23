package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/Ttae27/Food_Recipe_Back/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect ot database")
	}

	db.AutoMigrate(&models.Bookmark{}, &models.Comment{}, &models.Ingredient{}, &models.IngredientCategory{}, &models.Ingredient_IngredientCategory{}, &models.Post_Like{}, &models.Post_Comment{}, &models.Post_Ingredient{}, &models.Post{}, &models.Category{}, &models.Post_Category{}, &models.User{}, &models.User_Comment{})

	app := fiber.New()

	routes.Routes_Post(db, app)
	log.Fatal(app.Listen(":8080"))
}
