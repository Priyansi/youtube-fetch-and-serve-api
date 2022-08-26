package apikeys

import (
	"context"
	"errors"

	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var collection *mongo.Collection

func SetCollection(client *mongo.Client) {
	collection = client.Database("youtube-fetch-search").Collection("api-keys")
}

func IsKeyValid(key string) bool {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		return false
	}

	call := youtubeService.Channels.List([]string{"id"}).ForUsername("Youtube")
	_, err = call.Do()
	return err == nil
}

func InsertKey(key string) error {
	_, err := collection.InsertOne(context.TODO(), bson.M{"key": key, "isExpired": false})
	if err != nil {
		logger.Error.Printf("InsertKey: Error inserting key: %v", err)
		return err
	}
	logger.Info.Println("InsertKey: Inserted new API key")
	return nil
}

func GetValidKey() (string, error) {
	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{"isExpired": false})
	if err != nil {
		logger.Error.Printf("GetValidKey: Error finding valid key: %v", err)
		return "", errors.New("error finding valid key")
	}
	defer cursor.Close(ctx)

	var key bson.M

	for cursor.Next(ctx) {
		err := cursor.Decode(&key)
		if err != nil {
			logger.Error.Printf("GetValidKey: Error decoding key: %v", err)
			continue
		}
		return key["key"].(string), nil
	}

	return "", errors.New("no valid key found")
}

func CheckValidityOfKeys() {
	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{"isExpired": false})
	if err != nil {
		logger.Error.Printf("CheckValidityOfKeys: Error finding valid key: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var key bson.M

	for cursor.Next(ctx) {
		err := cursor.Decode(&key)
		if err != nil {
			logger.Error.Printf("CheckValidityOfKeys: Error decoding key: %v", err)
			continue
		}
		if !IsKeyValid(key["key"].(string)) {
			collection.UpdateOne(context.TODO(), bson.M{"key": key["key"]}, bson.M{"$set": bson.M{"isExpired": true}})
		}
	}
}
