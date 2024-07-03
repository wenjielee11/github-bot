package handlers
import (
	"context"
	"github-bot/models"
	"github-bot/services"
	"log"
	"github.com/google/go-github/v41/github"
)


func HandlePullRequestEvent(ctx context.Context, client *github.Client, owner, repo string, eventPayload models.EventPayload) {
	if eventPayload.PullRequest == nil {
		log.Println("No pull request data found in payload")
		return
	}

	pr := eventPayload.PullRequest
	log.Printf("Processing pull request: #%d\n", pr.Number)

	services.CheckChangelogUpdated(ctx, client, owner, repo, pr)
	services.CheckSecretKeyLeakage(ctx, client, owner, repo, pr)
	services.SuggestLabelsForPR(ctx, client, owner, repo, pr)
}