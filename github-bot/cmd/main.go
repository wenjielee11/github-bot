package main

import (
	"github-bot/handlers"
	"github-bot/services"
	"github-bot/utils"
	"log"
	"os"
	"strconv"
)

const (
	owner      = "wenjielee1"
	repository = "chatbot-test"
)

func main() {
	log.Println("Starting the GitHub bot")

	appIDStr := os.Getenv("APP_ID")
	installationIDStr := os.Getenv("INSTALLATION_ID")
	privateKeyBase64 := os.Getenv("PRIVATE_KEY")

	log.Printf("APP_ID: %s", appIDStr)
	log.Printf("INSTALLATION_ID: %s", installationIDStr)

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error converting APP_ID to int64: %v", err)
	}

	installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error converting INSTALLATION_ID to int64: %v", err)
	}

	privateKey, err := utils.DecodePrivateKey(privateKeyBase64)
	if err != nil {
		log.Fatalf("Error decoding private key: %v", err)
	}

	jwtToken, err := utils.GenerateJWT(appID, privateKey)
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}

	installationToken, err := services.GetInstallationToken(installationID, jwtToken)
	if err != nil {
		log.Fatalf("Error getting installation token: %v", err)
	}

	handlers.HandleGitHubEvents(owner, repository, installationToken)
}
