package utils

import (
	"os"
	"strconv"

	"github.com/priyansi/fampay-backend-assignment/pkg/logger"
)

const DEFAULT_MAXTOKENS_VAL = 5

func GetEnvInt(key string) int64 {
	s := os.Getenv(key)
	if s == "" {
		logger.Info.Printf("Utils: environment variable %v not found, using default value %v.", key, DEFAULT_MAXTOKENS_VAL)
		return DEFAULT_MAXTOKENS_VAL
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		logger.Info.Printf("Utils: environment variable %v not found, using default value %v.", key, DEFAULT_MAXTOKENS_VAL)
		return DEFAULT_MAXTOKENS_VAL
	}
	return int64(v)
}
