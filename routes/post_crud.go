package routes

import (
	"github.com/Ttae27/Food_Recipe_Back/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes_Post(db *gorm.DB, app *fiber.App) {
	//Post endpoints
	app.Post("/post", func(c *fiber.Ctx) error {
		return controller.CreatePost(db, c)
	})
	app.Get("/post", func(c *fiber.Ctx) error {
		return controller.GetsPost(db, c)
	})
	app.Get("/post/:id", func(c *fiber.Ctx) error {
		return controller.GetPost(db, c)
	})
	// app.Put("/post/:id", func(c *fiber.Ctx) error {
	// 	return UpdatePost(db, c)
	// })
	// app.Delete("/post/:id", func(c *fiber.Ctx) error {
	// 	return DeleteBook(db, c)
	// })
}
