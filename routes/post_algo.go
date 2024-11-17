package routes

import (
	"github.com/Ttae27/Food_Recipe_Back/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes_Search_Post(db *gorm.DB, app *fiber.App) {
	app.Get("/posts/search", func(c *fiber.Ctx) error {
		return controller.SearchPostByNameAndFilters(db, c)
	})
}
