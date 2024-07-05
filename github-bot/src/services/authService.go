package services

import (
	"context"
	"os"
    "log"
	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"golang.org/x/oauth2"
)

func GetInstallationToken(installationID int64, jwtToken string) (string, error) {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: jwtToken},
    )
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
    if err != nil {
        return "", err
    }
    return token.GetToken(), nil
}

func GetJamAiHeader() *models.JamaiAuth{
    jamaiKey := os.Getenv("TRIAGE_BOT_JAMAI_KEY")
	if jamaiKey == "" {
		log.Fatalf("Error: JAMAI_KEY environment variable not set")
	}
    projectId :=os.Getenv("TRIAGE_BOT_JAMAI_PROJECT_ID")
    if projectId == "" {
		log.Fatalf("Error: JAMAI_PROJECT_ID environment variable not set")
	}
	return &models.JamaiAuth{
		Authorization: "Bearer " + jamaiKey,
		XProjectID: projectId,
	}  
  
}
