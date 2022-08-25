package db

import (
	"context"
	"time"

	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.GetMongoDbURI()))
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
	defer client.Disconnect(ctx)
}
