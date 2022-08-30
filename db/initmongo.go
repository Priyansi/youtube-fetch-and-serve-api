package db

import (
	"context"
	"log"
	"time"

	"github.com/priyansi/fampay-backend-assignment/db/apikeys"
	"github.com/priyansi/fampay-backend-assignment/db/youtubevideoinfo"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDb() {
	client := ConnectToMongoDb()
	youtubevideoinfo.SetCollection(client)
	youtubevideoinfo.CreateTitleAndDescriptionIndex()
	apikeys.SetCollection(client)
}

func ConnectToMongoDb() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI((config.GetMongoDbURI())).SetServerAPIOptions(serverAPIOptions))
	if err != nil {
		log.Fatalf("ConnectMongo: Error connecting to mongo db: %v", err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("ConnectMongo: Error pinging mongo db: %v", err)
	}
	return client
}
