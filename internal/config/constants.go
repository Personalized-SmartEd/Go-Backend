package config

import (
	"os"
)

func getBaseURL() string {
	return os.Getenv("MLSERVICE_URL")
}

var BaseURL = getBaseURL()
