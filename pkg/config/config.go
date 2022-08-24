package config

import (
	"flag"
	"os"
)

type Config struct {
	AddrPort string
	ApiKey   string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.AddrPort, "addrport", os.Getenv("ADDR_PORT"), "Port where server is running")
	flag.StringVar(&conf.ApiKey, "apikey", os.Getenv("API_KEY"), "YouTube API Key")

	flag.Parse()

	return conf
}
