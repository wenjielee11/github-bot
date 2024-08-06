package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/services"
	"github.com/wenjielee1/github-bot/utils"
)

// HandleIssueEvent processes GitHub issue events by extracting issue data from the event payload
// and delegating the processing to the service layer.
func HandleIssueEvent(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, eventPayload models.EventPayload) {
	// Check if the issue data is present in the event payload
	if eventPayload.Issue == nil {
		log.Println("No issue data found in payload")
		return
	}

	// Extract the issue data from the event payload
	issue := eventPayload.Issue

	// Log the title of the issue being processed
	log.Printf("Processing issue: %s", issue.Title)

	// Delegate the processing of the issue to the services layer
	services.ProcessIssue(ctx, client, jamaiClient, fmt.Sprintf("%s_%s_%s", owner, repo, utils.GetBotVersion()), owner, repo, issue)
}
