package services

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"golang.org/x/oauth2"
)

// GetInstallationToken retrieves an installation token for the GitHub App.
// It uses the provided installation ID and JWT token to authenticate.
func GetInstallationToken(installationID int64, jwtToken string) (string, error) {
	// Create a new context
	ctx := context.Background()

	// Initialize the OAuth2 token source and client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: jwtToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create an installation token using the GitHub client
	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", err
	}
	return token.GetToken(), nil
}

// GetJamAiHeader retrieves the JAM.AI authentication header.
// It fetches the JAM.AI key and project ID from environment variables.
func GetJamAiHeader() *models.JamaiAuth {
	// Retrieve the JAM.AI key from environment variables
	jamaiKey := os.Getenv("TRIAGE_BOT_JAMAI_KEY")
	if jamaiKey == "" {
		log.Fatalf("Error: JAMAI_KEY environment variable not set")
	}

	// Retrieve the JAM.AI project ID from environment variables
	projectId := os.Getenv("TRIAGE_BOT_JAMAI_PROJECT_ID")
	if projectId == "" {
		log.Fatalf("Error: JAMAI_PROJECT_ID environment variable not set")
	}

	// Return the JAM.AI authentication header
	return &models.JamaiAuth{
		Authorization: "Bearer " + jamaiKey,
		XProjectID:    projectId,
	}
}
