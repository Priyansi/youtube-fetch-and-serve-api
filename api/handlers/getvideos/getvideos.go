package getvideos

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/priyansi/fampay-backend-assignment/pkg/config"
	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	ytServiceOrderBy        = "date"
	ytServiceType           = "video"
	ytServicePublishedAfter = "2022-01-01T00:00:00Z"
)

type Video struct {
	Title       string
	Description string
	PublishedAt string
}

func Do() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.Background()
		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(config.GetApiKey()))
		if err != nil {
			logger.Error.Fatalf("GetVideos: Error creating new service: %v", err)
		}
		call := youtubeService.Search.List([]string{"id,snippet"}).
			Q(config.GetQuery()).
			MaxResults(config.GetMaxResults()).
			Order(ytServiceOrderBy).
			Type(ytServiceType).
			PublishedAfter(ytServicePublishedAfter).
			PageToken(config.GetPageToken())

		response, err := call.Do()
		if err != nil {
			logger.Error.Printf("GetVideos: Error fetching response: %v", err)
		}

		logger.Info.Println("GetVideos: Fetched response from youtube")
		for _, item := range response.Items {
			video := Video{
				Title:       item.Snippet.Title,
				Description: item.Snippet.Description,
				PublishedAt: item.Snippet.PublishedAt,
			}
			fmt.Println(video)
		}
		fmt.Printf("GetVideos: page tokens %v %v", response.NextPageToken, response.PrevPageToken)
		config.SetPageToken(response.NextPageToken)
	}
}
