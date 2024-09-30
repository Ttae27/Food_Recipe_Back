package routes

import (
	"github.com/Ttae27/Food_Recipe_Back/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes_Search_Post(db *gorm.DB, app *fiber.App) {
	app.Get("/posts/search", func(c *fiber.Ctx) error {
		return controller.SearchPostByName(db, c)
	})

	app.Get("/posts/filter/price", func(c *fiber.Ctx) error {
		return controller.GetPostsByPriceRange(db, c)
	})

	app.Get("/posts/filter/category", func(c *fiber.Ctx) error {
		return controller.GetPostsByCategoryType(db, c)
	})
}
