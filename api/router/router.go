package router

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/priyansi/fampay-backend-assignment/api/handlers/getvideos"
	"github.com/priyansi/fampay-backend-assignment/api/handlers/postapikey"
	"github.com/priyansi/fampay-backend-assignment/api/handlers/searchvideos"
)

func SetRoutes(app *fiber.App) {
	app.Get("/get_videos", func(c *fiber.Ctx) error {
		return getvideos.Do(c)
	})

	app.Get("/search_videos", func(c *fiber.Ctx) error {
		return searchvideos.Do(c)
	})

	app.Post("/post_api_key", func(c *fiber.Ctx) error {
		return postapikey.Do(c)
	})
}
