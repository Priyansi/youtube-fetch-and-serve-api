package config

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

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
	ValidApiKey                    string
}

const (
	DEFAULT_MAX_TOKENS                         = 5
	DEFAULT_PER_PAGE_LIMIT                     = 5
	DEFAULT_FETCH_LATEST_VIDEOS_SECONDS        = 10
	DEFAULT_UPDATE_API_KEYS_EXPIRATION_MINUTES = 120
)

var config Config

// Initializes config with environment variables or default values if not specified
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
	if config.MaxVideosFetched > 50 || config.MaxVideosFetched < 1 {
		log.Infof("Config: Environment variable MAX_VIDEOS_FETCHED should be between 1 and 50. Please refer to README. Setting it to default value: %d", DEFAULT_MAX_TOKENS)
		config.MaxVideosFetched = DEFAULT_MAX_TOKENS
	}

	flag.Int64Var(&config.PerPageLimit, "perpagelimit", utils.GetEnvInt("PER_PAGE_LIMIT", DEFAULT_PER_PAGE_LIMIT), "Number of videos to be displayed per page")
	if config.PerPageLimit < 1 {
		log.Infof("Config: Environment variable PER_PAGE_LIMIT should be greater than 0. Please refer to README. Setting it to default value: %d", DEFAULT_PER_PAGE_LIMIT)
		config.PerPageLimit = DEFAULT_PER_PAGE_LIMIT
	}

	flag.Int64Var(&config.FetchLatestVideosSeconds, "fetchlatestvideosseconds", utils.GetEnvInt("FETCH_LATEST_VIDEOS_SECONDS", DEFAULT_FETCH_LATEST_VIDEOS_SECONDS), "Number of seconds after which latest videos are fetched from youtube and database is updated")
	if config.FetchLatestVideosSeconds < 1 {
		log.Infof("Config: Environment variable FETCH_LATEST_VIDEOS_SECONDS should be greater than 0. Please refer to README. Setting it to default value: %d", DEFAULT_FETCH_LATEST_VIDEOS_SECONDS)
		config.FetchLatestVideosSeconds = DEFAULT_FETCH_LATEST_VIDEOS_SECONDS
	}

	flag.Int64Var(&config.UpdateApiKeysExpirationMinutes, "updateapikeysexpirationminutes", utils.GetEnvInt("UPDATE_API_KEYS_EXPIRATION_MINUTES", DEFAULT_UPDATE_API_KEYS_EXPIRATION_MINUTES), "Number of minutes after which expired api keys whose quota has exceeded are checked for validity and updated")
	if config.UpdateApiKeysExpirationMinutes < 1 {
		log.Infof("Config: Environment variable UPDATE_API_KEYS_EXPIRATION_MINUTES should be greater than 0. Please refer to README. Setting it to default value: %d", DEFAULT_UPDATE_API_KEYS_EXPIRATION_MINUTES)
		config.UpdateApiKeysExpirationMinutes = DEFAULT_UPDATE_API_KEYS_EXPIRATION_MINUTES
	}

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

func GetValidApiKey() string {
	return config.ValidApiKey
}

func SetValidApiKey(apiKey string) {
	config.ValidApiKey = apiKey
}

func SetEtag(etag string) {
	config.Etag = etag
}
