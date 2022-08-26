package main

import (
	"log"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/priyansi/fampay-backend-assignment/api/router"
	"github.com/priyansi/fampay-backend-assignment/db"
	"github.com/priyansi/fampay-backend-assignment/db/apikeys"
	"github.com/priyansi/fampay-backend-assignment/db/youtubevideoinfo"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("main: failed to load environment variables")
	}

	config.InitConfig()
	db.InitMongoDb()

	go func() {
		ticker := time.NewTicker(time.Duration(config.GetFetchLatestVideosSeconds()) * time.Second)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				err := youtubevideoinfo.FetchNewVideosAndUpdateDb()
				if err != nil {
					logger.Error.Printf("main: error fetching new videos and updating db: %v", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()

	go func() {
		ticker := time.NewTicker(time.Duration(config.GetCheckApiKeysValidityMinutes()) * time.Minute)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				apikeys.CheckValidityOfKeys()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	app := fiber.New()
	router.SetRoutes(app)

	app.Listen(config.GetAddrPort())
}
