package searchvideos

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/priyansi/fampay-backend-assignment/db/youtubevideoinfo"
)

const perPageLimit = 5

func Do(c *fiber.Ctx) error {
	searchQuery := c.Query("query", "")
	if searchQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query param is required",
		})
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "page query param must be an integer",
		})
	}

	return c.JSON(fiber.Map{
		"videos": youtubevideoinfo.SearchVideos(searchQuery, perPageLimit, int64(page)),
	})
}
