package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/priyansi/fampay-backend-assignment/api/router"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	"github.com/priyansi/fampay-backend-assignment/pkg/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("main: failed to load environment variables")
	}
	config.InitConfig()
	srv := server.
		Get().
		SetAddrPort(config.GetAddrPort()).
		SetRouter(router.Get()).
		SetLogger(logger.Error)

	go func() {
		logger.Info.Printf("Starting server at %s", config.GetAddrPort())
		if err := srv.Start(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	switch <-ch {
	case os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT:
		logger.Info.Println("Interrupt signal received, exiting.")
		logger.Info.Println("Last page token: ", config.GetPageToken())
		logger.Info.Println("Server shutdown")
		os.Exit(0)
	default:
		logger.Info.Println("Something went wrong.")
		logger.Info.Println("Last page token: ", config.GetPageToken())
		logger.Info.Println("Server shutdown")
		os.Exit(1)
	}
}
