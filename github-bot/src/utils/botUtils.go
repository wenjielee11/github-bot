package utils

import (
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
		return defaultValue
	}
	return value
}

func GetRepoName(defaultValue string) string {
	value, exists := os.LookupEnv("REPO_NAME")
	if !exists {
		return defaultValue
	}
	return value
}
