package main

import (
	"log"
	"os"
	"strconv"

	"github.com/wenjielee1/github-bot/handlers"
	"github.com/wenjielee1/github-bot/services"
	"github.com/wenjielee1/github-bot/utils"
)

// Constants for the repository owner and name.
const (
	owner      = "EmbeddedLLM"
	repository = "JAM.ai.dev"
)

// main is the entry point of the GitHub bot application.
// It retrieves necessary credentials from environment variables, generates tokens,
// and starts handling GitHub events.
func main() {
	log.Println("Starting the GitHub bot")
	// Retrieve the GitHub App ID from environment variables.
	// The GitHub App ID is a unique identifier for the GitHub App. It is assigned by GitHub when the app is created.
	appIDStr := os.Getenv("TRIAGE_BOT_APP_ID")
	log.Printf("APP_ID: %s", appIDStr)

	// Retrieve the Installation ID from environment variables.
	// The Installation ID is a unique identifier for the installation of the GitHub App on a specific repository or organization.
	installationIDStr := os.Getenv("TRIAGE_BOT_INSTALLATION_ID")
	log.Printf("INSTALLATION_ID: %s", installationIDStr)

	// Retrieve the private key from environment variables and decode it from base64 format.
	// The private key is used to sign JWT tokens for authenticating as the GitHub App.
	privateKeyBase64 := os.Getenv("TRIAGE_BOT_PRIVATE_KEY")

	// Convert appIDStr to int64.
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error converting APP_ID to int64: %v", err)
	}

	// Convert installationIDStr to int64.
	installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error converting INSTALLATION_ID to int64: %v", err)
	}

	// Decode the private key from base64 format.
	privateKey, err := utils.DecodePrivateKey(privateKeyBase64)
	if err != nil {
		log.Fatalf("Error decoding private key: %v", err)
	}

	// Generate a JWT token for authentication.
	// The JWT token is used to authenticate as the GitHub App and is required to perform actions on behalf of the app.
	jwtToken, err := utils.GenerateJWT(appID, privateKey)
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}

	// Retrieve the installation token using the JWT token.
	// The installation token is used to authenticate API requests for a specific installation of the GitHub App.
	installationToken, err := services.GetInstallationToken(installationID, jwtToken)
	if err != nil {
		log.Fatalf("Error getting installation token: %v", err)
	}

	// Handle GitHub events using the installation token.
	handlers.HandleGitHubEvents(utils.GetRepoOwner(owner), utils.GetRepoName(repository), installationToken)
}
