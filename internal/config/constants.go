package config

import (
	"os"
)

func getBaseURL() string {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	return os.Getenv("MLSERVICE_URL")
}

var BaseURL = getBaseURL()
