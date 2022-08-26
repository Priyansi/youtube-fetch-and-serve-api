package utils

import (
	"os"
	"strconv"

	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
)

func GetEnvInt(key string, defaultVal int64) int64 {
	s := os.Getenv(key)
	if s == "" {
		logger.Info.Printf("Utils: environment variable %v not found, using default value %v.", key, defaultVal)
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		logger.Info.Printf("Utils: environment variable %v not found, using default value %v.", key, defaultVal)
		return defaultVal
	}
	return int64(v)
}
