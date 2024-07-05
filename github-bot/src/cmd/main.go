package main

import (
	"log"
	"os"
	"strconv"

	"github.com/wenjielee1/github-bot/services"

	"github.com/wenjielee1/github-bot/utils"

	"github.com/wenjielee1/github-bot/handlers"
)

const (
	owner      = "EmbeddedLLM"
	repository = "JAM.ai.dev"
)

func main() {
	log.Println("Starting the GitHub bot")

	appIDStr := os.Getenv("GITHUB_BOT_APP_ID")
	installationIDStr := os.Getenv("GITHUB_BOT_INSTALLATION_ID")
	privateKeyBase64 := os.Getenv("GITHUB_BOT_PRIVATE_KEY")

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
