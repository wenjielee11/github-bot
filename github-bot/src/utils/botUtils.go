package utils

import (
	"log"
	"os"
)

const (
	BotVersion = "v0.5.0"
)

func GetBotVersion() string {
	return BotVersion
}

func GetRepoOwner(defaultValue string) string {
	value, exists := os.LookupEnv("REPO_OWNER")
	if !exists {
		log.Printf("Repo owner not found. Defaulting to %s", defaultValue)
		return defaultValue
	}
	return value
}

func GetRepoName(defaultValue string) string {
	value, exists := os.LookupEnv("REPO_NAME")
	if !exists {
		log.Printf("Repo name not found. Defaulting to %s", defaultValue)
		return defaultValue
	}
	return value
}
