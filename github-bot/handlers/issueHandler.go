package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wenjielee1/github-bot/services"


	"github.com/wenjielee1/github-bot/models"

	"github.com/google/go-github/v41/github"
)

func HandleIssueEvent(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, eventPayload models.EventPayload) {
	if eventPayload.Issue == nil {
		log.Println("No issue data found in payload")
		return
	}

	issue := eventPayload.Issue

	log.Printf("Processing issue: %s", issue.Title)

	services.ProcessIssue(ctx, client, jamaiClient, fmt.Sprintf("%s_%s", owner, repo), owner, repo, issue)
	
}
