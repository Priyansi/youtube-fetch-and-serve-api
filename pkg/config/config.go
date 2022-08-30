package config

import (
	"flag"
	"log"
	"os"

	"github.com/priyansi/fampay-backend-assignment/pkg/utils"
)

type Config struct {
	AddrPort                       string
	Etag                           string
	MaxVideosFetched               int64
	PerPageLimit                   int64
	FetchLatestVideosSeconds       int64
	UpdateApiKeysExpirationMinutes int64
	Query                          string
	MongoDbURI                     string
}

const (
	DEFAULT_MAX_TOKENS                         = 5
	DEFAULT_PER_PAGE_LIMIT                     = 5
	DEFAULT_FETCH_LATEST_VIDEOS_SECONDS        = 10
	DEFAULT_UPDATE_API_KEYS_EXPIRATION_MINUTES = 120
)

var config Config

func InitConfig() {
	flag.StringVar(&config.AddrPort, "addrport", os.Getenv("ADDR_PORT"), "Port where server is running")
	if config.AddrPort == "" {
		log.Fatalf("Config: Environment variable ADDR_PORT not found. Please refer to README to find how to set it.")
	}
	flag.StringVar(&config.MongoDbURI, "mongodburi", os.Getenv("MONGODB_URI"), "MongoDB URI for connection")
	if config.MongoDbURI == "" {
		log.Fatalf("Config: Environment variable MONGODB_URI not found. Please refer to README to find how to set it.")
	}
	flag.StringVar(&config.Query, "query", os.Getenv("QUERY"), "Predefined search query")
	if config.Query == "" {
		log.Fatalf("Config: Environment variable QUERY not found. Please refer to README to find how to set it.")
	}
	flag.Int64Var(&config.MaxVideosFetched, "maxvideosfetched", utils.GetEnvInt("MAX_VIDEOS_FETCHED", DEFAULT_MAX_TOKENS), "Max videos that can be fetched in a single API call")
	flag.Int64Var(&config.PerPageLimit, "perpagelimit", utils.GetEnvInt("PER_PAGE_LIMIT", DEFAULT_PER_PAGE_LIMIT), "Number of videos to be displayed per page")
	flag.Int64Var(&config.FetchLatestVideosSeconds, "fetchlatestvideosseconds", utils.GetEnvInt("FETCH_LATEST_VIDEOS_SECONDS", DEFAULT_FETCH_LATEST_VIDEOS_SECONDS), "Number of seconds after which latest videos are fetched from youtube and database is updated")
	flag.Int64Var(&config.UpdateApiKeysExpirationMinutes, "updateapikeysexpirationminutes", utils.GetEnvInt("UPDATE_API_KEYS_EXPIRATION_MINUTES", DEFAULT_UPDATE_API_KEYS_EXPIRATION_MINUTES), "Number of minutes after which expired api keys whose quota has exceeded are checked for validity and updated")

	flag.Parse()

	// etag will be empty for the first API call
	config.Etag = ""
}

func GetAddrPort() string {
	return config.AddrPort
}

func GetQuery() string {
	return config.Query
}

func GetMaxVideosFetched() int64 {
	return config.MaxVideosFetched
}

func GetPerPageLimit() int64 {
	return config.PerPageLimit
}

func GetFetchLatestVideosSeconds() int64 {
	return config.FetchLatestVideosSeconds
}

func GetUpdateApiKeysExpirationMinutes() int64 {
	return config.UpdateApiKeysExpirationMinutes
}

func GetMongoDbURI() string {
	return config.MongoDbURI
}

func GetEtag() string {
	return config.Etag
}

func SetEtag(etag string) {
	config.Etag = etag
}
