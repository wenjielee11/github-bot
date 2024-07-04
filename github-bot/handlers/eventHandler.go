package handlers

import (
	"context"
	"encoding/json"

	"github.com/wenjielee1/github-bot/models"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func HandleGitHubEvents(owner, repo, token string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	log.Printf("Event Name: %s", eventName)
	log.Printf("Event Path: %s", eventPath)

	eventData, err := ioutil.ReadFile(eventPath)
	if err != nil {
		log.Fatalf("Error reading event data: %v", err)
	}

	var eventPayload models.EventPayload
	if err := json.Unmarshal(eventData, &eventPayload); err != nil {
		log.Fatalf("Error parsing event data: %v", err)
	}

	switch eventName {
	case "issues":
		
		HandleIssueEvent(ctx, client, owner, repo, eventPayload)
	case "pull_request":
		HandlePullRequestEvent(ctx, client, owner, repo, eventPayload)
	default:
		log.Printf("Unhandled event: %s", eventName)
	}
}
