package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/services"
)

// HandlePullRequestEvent processes GitHub pull request events by extracting pull request data from the event payload
// and delegating various checks and actions to the service layer.
func HandlePullRequestEvent(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, eventPayload models.EventPayload) {
	// Check if the pull request data is present in the event payload
	if eventPayload.PullRequest == nil {
		log.Println("No pull request data found in payload")
		return
	}

	// Extract the pull request data from the event payload
	pr := eventPayload.PullRequest

	// Log the pull request number being processed
	log.Printf("Processing pull request: #%d\n", pr.Number)

	// Cleanup of previous bot comments on a PR synchronize.
	if eventPayload.Action == "synchronize" {
		services.DeleteBotComments(ctx, client, jamaiClient, owner, repo, pr, "jambubot")
	}

	// Delegate various checks and actions to the services layer
	services.CheckChangelogUpdated(ctx, client, jamaiClient, owner, repo, pr)
	services.CheckSecretKeyLeakage(ctx, client, jamaiClient, owner, repo, pr)

	// services.SuggestLabelsForPR(ctx, client, owner, repo, pr)
}
