package apikeys

import "go.mongodb.org/mongo-driver/mongo"

var collection *mongo.Collection

func SetCollection(client *mongo.Client) {
	collection = client.Database("youtube-fetch-search").Collection("api-keys")
}
