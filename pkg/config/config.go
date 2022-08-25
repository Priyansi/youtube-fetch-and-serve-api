package config

import (
	"flag"
	"os"

	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
	"github.com/priyansi/fampay-backend-assignment/pkg/utils"
)

type Config struct {
	AddrPort   string
	ApiKey     string
	PageToken  string
	MaxResults int64
	Query      string
	MongoDbURI string
}

var config Config

func InitConfig() {
	flag.StringVar(&config.AddrPort, "addrport", os.Getenv("ADDR_PORT"), "Port where server is running")
	if config.AddrPort == "" {
		logger.Error.Fatalf("Config: Environment variable ADDR_PORT not found. Please refer to README to find how to set it.")
	}
	flag.StringVar(&config.MongoDbURI, "mongodburi", os.Getenv("MONGODB_URI"), "MongoDB URI for connection")
	if config.MongoDbURI == "" {
		logger.Error.Fatalf("Config: Environment variable MONGODB_URI not found. Please refer to README to find how to set it.")
	}
	flag.StringVar(&config.Query, "query", os.Getenv("QUERY"), "Predefined search query")
	if config.Query == "" {
		logger.Error.Fatalf("Config: Environment variable QUERY not found. Please refer to README to find how to set it.")
	}
	flag.StringVar(&config.ApiKey, "apikey", os.Getenv("API_KEY"), "YouTube API Key")
	flag.StringVar(&config.PageToken, "pagetoken", os.Getenv("PAGE_TOKEN"), "Used for pagination when fetching responses subsequently")
	flag.Int64Var(&config.MaxResults, "maxresults", utils.GetEnvInt("MAX_RESULTS"), "Max results that can be fetched in a single API call")

	flag.Parse()
}

func GetAddrPort() string {
	return config.AddrPort
}

func GetApiKey() string {
	return config.ApiKey
}

func GetPageToken() string {
	return config.PageToken
}

func GetQuery() string {
	return config.Query
}

func GetMaxResults() int64 {
	return config.MaxResults
}

func GetMongoDbURI() string {
	return config.MongoDbURI
}

func SetPageToken(pageToken string) {
	config.PageToken = pageToken
}
