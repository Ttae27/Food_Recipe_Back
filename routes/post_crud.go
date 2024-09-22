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
	app.Put("/post/:id", func(c *fiber.Ctx) error {
		return controller.UpdatePost(db, c)
	})
	app.Delete("/post/:id", func(c *fiber.Ctx) error {
		return controller.DeletePost(db, c)
	})

	//Comment endpoints
	app.Post("/comment", func(c *fiber.Ctx) error {
		return controller.AddComment(db, c)
	})
	app.Delete("/comment", func(c *fiber.Ctx) error {
		return controller.DeleteComment(db, c)
	})

	//Like endpoints
	app.Post("/like", func(c *fiber.Ctx) error {
		return controller.AddLike(db, c)
	})
	app.Delete("/like", func(c *fiber.Ctx) error {
		return controller.DeleteLike(db, c)
	})

	//Bookmark endpoints
	app.Post("/bookmark", func(c *fiber.Ctx) error {
		return controller.AddBookmark(db, c)
	})
	app.Delete("/bookmark", func(c *fiber.Ctx) error {
		return controller.DeleteBookmark(db, c)
	})
}
