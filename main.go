package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ttae27/Food_Recipe_Back/initializers"
	"github.com/Ttae27/Food_Recipe_Back/models"
	"github.com/Ttae27/Food_Recipe_Back/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// init function to load environment variables
func init() {
	initializers.LoadEnvVariables()
}

// Function to retrieve environment variables and generate DSN
func getDatabaseDSN() (string, error) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Check if any of the necessary variables are empty, return error
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return "", fmt.Errorf("missing one or more environment variables: HOST=%s, PORT=%s, USER=%s, DB_NAME=%s", host, port, user, dbname)
	}

	// Generate and return the DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return dsn, nil
}

func main() {
	// Get the database DSN string
	dsn, err := getDatabaseDSN()
	if err != nil {
		log.Fatal("Database configuration error: ", err)
	}

	// Create a new logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // Log SQL queries that are slower than 1 second
			LogLevel:      logger.Info, // Set log level to Info
			Colorful:      true,        // Enable colorful output
		},
	)

	// Connect to the PostgreSQL database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to the database")
	}

	// AutoMigrate to ensure database schema matches models
	db.AutoMigrate(
		&models.Bookmark{},
		&models.Comment{},
		&models.Ingredient{},
		&models.IngredientCategory{},
		&models.Ingredient_IngredientCategory{},
		&models.Post_Like{},
		&models.Post_Comment{},
		&models.Post_Ingredient{},
		&models.Post{},
		&models.Category{},
		&models.Post_Category{},
		&models.User{},
		&models.User_Comment{},
	)

	// Create a new Fiber app
	app := fiber.New()

	// Set up routes for Posts
	routes.Routes_Post(db, app)
	routes.Routes_Search_Post(db, app)
	routes.Routes_User(db, app)

	// Get the server port from environment variable and start server
	serverPort := ":" + os.Getenv("SERVER_PORT")
	log.Fatal(app.Listen(serverPort))
}
