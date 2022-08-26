package db

import (
	"context"
	"time"

	"github.com/priyansi/fampay-backend-assignment/db/apikeys"
	"github.com/priyansi/fampay-backend-assignment/db/youtubevideoinfo"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDb() {
	client := ConnectToMongoDb()
	youtubevideoinfo.SetCollection(client)
	apikeys.SetCollection(client)
}

func ConnectToMongoDb() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI((config.GetMongoDbURI())))
	if err != nil {
		logger.Error.Fatalf("ConnectMongo: Error connecting to mongo db: %v", err)
	}
	return client
}
