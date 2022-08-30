package youtubevideoinfo

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/priyansi/fampay-backend-assignment/db/apikeys"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var collection *mongo.Collection

const (
	ytServiceOrderBy        = "date"
	ytServiceType           = "video"
	ytServicePublishedAfter = "2022-01-01T00:00:00Z"
)

func SetCollection(client *mongo.Client) {
	collection = client.Database("youtube-fetch-search").Collection("youtube-video-info")
}

func CreateTitleAndDescriptionIndex() {
	model := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "description", Value: "text"},
		},
		Options: options.Index().SetWeights(bson.D{
			{Key: "title", Value: 1},
			{Key: "description", Value: 1},
		}).SetCollation(&options.Collation{
			Locale: "simple"}),
	}

	options := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := collection.Indexes().CreateOne(context.TODO(), model, options)
	if err != nil {
		log.Fatalf("SetCollection: Error creating index: %v", err)
	}
}

// insert many entries into the youtube-video-info collection
func bulkInsert(videos []types.Video) error {
	models := make([]mongo.WriteModel, 0)
	for _, video := range videos {
		videoBson := bson.M{
			"title":       video.Title,
			"description": video.Description,
			"publishedAt": video.PublishedAt,
		}
		query := bson.M{
			"$setOnInsert": videoBson,
			"$set":         bson.M{},
		}

		models = append(models, mongo.NewUpdateOneModel().SetUpsert(true).SetUpdate(query).SetFilter(bson.M{"uniqueId": video.UniqueId}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	res, err := collection.BulkWrite(context.Background(), models, opts)
	if err != nil {
		log.Errorf("BulkInsert: Error inserting many: %v", err)
		return err
	}
	log.Infof("BulkInsert: Inserted %v documents into collection", res.UpsertedCount)
	return nil
}

func FetchNewVideosAndUpdateDb() error {
	ctx := context.Background()
	key, err := apikeys.GetValidKey()
	if err != nil {
		log.Errorf("FetchNewVideosAndUpdateDb: Error fetching valid key: %v", err)
		log.Errorf("FetchNewVideosAndUpdateDb: Please post new api keys as given in README")
		return err
	}
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Errorf("FetchNewVideosAndUpdateDb: Error creating new service: %v", err)
	}

	call := youtubeService.Search.List([]string{"id,snippet"}).
		Q(config.GetQuery()).
		MaxResults(config.GetMaxVideosFetched()).
		Order(ytServiceOrderBy).
		Type(ytServiceType).
		PublishedAfter(ytServicePublishedAfter)

	if config.GetEtag() != "" {
		call = call.IfNoneMatch(config.GetEtag())
	}

	response, err := call.Do()
	if err != nil {
		if err.Error() == "googleapi: got HTTP response code 304 with body: " {
			log.Info("FetchNewVideosAndUpdateDb: Etag has not changed. Skipping update.")
			return nil
		}
		if strings.Contains(err.Error(), "quotaExceeded") {
			log.Info("FetchNewVideosAndUpdateDb: Quota exceeded. Setting key to expired. Skipping update.")
			err := apikeys.SetKeyToExpired(key)
			if err != nil {
				log.Errorf("FetchNewVideosAndUpdateDb: Error setting quota exceeded key to expired: %v", err)
			}
			return err
		}
		log.Errorf("FetchNewVideosAndUpdateDb: Error fetching response: %v", err)
		return err
	}

	log.Info("FetchNewVideosAndUpdateDb: Fetched response from youtube")
	videos := make([]types.Video, 0)
	for _, item := range response.Items {
		video := types.Video{
			UniqueId:    item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
		}
		videos = append(videos, video)
	}

	log.Infof("FetchNewVideosAndUpdateDb: Fetched %v videos. Updating the database.", len(videos))
	err = bulkInsert(videos)
	if err != nil {
		log.Errorf("FetchNewVideosAndUpdateDb: Error inserting into db: %v", err)
		return err
	}

	config.SetEtag(response.Etag)
	return nil
}

func getFindOptions(perPageLimit int64, currPage int64) *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetSkip((currPage - 1) * perPageLimit)
	findOptions.SetLimit(perPageLimit)
	return findOptions
}

func GetVideos(currPage int64) []types.Video {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{}, getFindOptions(config.GetPerPageLimit(), currPage))
	if err != nil {
		log.Errorf("GetVideos: Error fetching videos: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	videos := make([]types.Video, 0)
	for cursor.Next(ctx) {
		var video types.Video
		err := cursor.Decode(&video)
		if err != nil {
			log.Errorf("GetVideos: Error decoding video: %v", err)
			continue
		}
		videos = append(videos, video)
	}
	return videos
}

func SearchVideos(query string, currPage int64) []types.Video {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Infof("SearchVideos: Searching for %v", query)

	perPageLimit := config.GetPerPageLimit()
	// search only in title and description
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	firstMatchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	addFieldsStage := bson.D{
		{Key: "$addFields", Value: bson.M{
			"score": bson.M{
				"$meta": "textScore",
			},
		}},
	}
	secondMatchStage := bson.D{
		{Key: "$match", Value: bson.M{
			"score": bson.M{
				"$gte": 1,
			},
		}},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.M{
			"score": -1,
		}},
	}
	setSkip := bson.D{
		{Key: "$skip", Value: (currPage - 1) * perPageLimit},
	}
	setLimit := bson.D{
		{Key: "$limit", Value: perPageLimit},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{firstMatchStage, addFieldsStage, secondMatchStage, sortStage, setSkip, setLimit})

	if err != nil {
		log.Errorf("SearchVideos: Error searching videos: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var results []bson.M
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Errorf("SearchVideos: Error decoding videos: %v", err)
		return nil
	}

	videos := make([]types.Video, 0)
	for _, result := range results {
		video := types.Video{
			UniqueId:    result["uniqueId"].(string),
			Title:       result["title"].(string),
			Description: result["description"].(string),
			PublishedAt: result["publishedAt"].(string),
		}
		videos = append(videos, video)
	}
	return videos
}
