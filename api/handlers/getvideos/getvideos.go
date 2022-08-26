package getvideos

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/priyansi/fampay-backend-assignment/db/youtubevideoinfo"
)

func Do(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "page query param must be an integer",
		})
	}

	videos := youtubevideoinfo.GetVideos(int64(page))

	return c.JSON(fiber.Map{
		"videos": videos,
	})
}
