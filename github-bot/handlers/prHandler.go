package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/services"
)

func HandlePullRequestEvent(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, eventPayload models.EventPayload) {
	if eventPayload.PullRequest == nil {
		log.Println("No pull request data found in payload")
		return
	}

	pr := eventPayload.PullRequest
	log.Printf("Processing pull request: #%d\n", pr.Number)

	
	services.CheckChangelogUpdated(ctx, client, jamaiClient, owner, repo, pr)
	services.CheckSecretKeyLeakage(ctx, client, owner, repo, pr)
	services.SuggestLabelsForPR(ctx, client, owner, repo, pr)
}
