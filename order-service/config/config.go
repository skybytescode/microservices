package config

import (
	"os"
	"strconv"
)

func GetEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

func GetDataSourceURL() string {
	return os.Getenv("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	port, err := strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	if err != nil {
		return 9000
	}
	return port
}
