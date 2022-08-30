package main

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/priyansi/fampay-backend-assignment/api/router"
	"github.com/priyansi/fampay-backend-assignment/db"
	"github.com/priyansi/fampay-backend-assignment/db/apikeys"
	"github.com/priyansi/fampay-backend-assignment/db/youtubevideoinfo"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		log.Fatal("main: failed to load environment variables")
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
					log.Errorf("main: error fetching new videos and updating db: %v", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()

	go func() {
		ticker := time.NewTicker(time.Duration(config.GetUpdateApiKeysExpirationMinutes()) * time.Minute)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				apikeys.UpdateExpirationOfExpiredKeys()
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
