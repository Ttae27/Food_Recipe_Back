package routes

import (
	"fmt"
	"github.com/Ttae27/Food_Recipe_Back/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes_User(db *gorm.DB, app *fiber.App) {
	app.Post("/user", func(c *fiber.Ctx) error {
		return controller.CreateUser(db, c)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return controller.LoginUser(db, c)
	})
	app.Get("/user/:id", requireAuth, func(c *fiber.Ctx) error {
		return controller.GetUser(db, c)
	})
	app.Put("/user/:id", func(c *fiber.Ctx) error {
		return controller.UpdateUser(db, c)
	})
	app.Delete("/user/:id", func(c *fiber.Ctx) error {
		return controller.DeleteUser(db, c)
	})
	app.Post("/logout/:id", func(c *fiber.Ctx) error {
		fmt.Println("TestLogout")
		return controller.LogoutUser(c)
	})
}
