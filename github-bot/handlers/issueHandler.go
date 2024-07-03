package handlers

import (
	"context"
	"log"
	"strings"

	"github.com/wenjielee1/github-bot/services"

	"github.com/wenjielee1/github-bot/utils"

	"github.com/wenjielee1/github-bot/models"

	"github.com/google/go-github/v41/github"
)

func HandleIssueEvent(ctx context.Context, client *github.Client, owner, repo string, eventPayload models.EventPayload) {
	if eventPayload.Issue == nil {
		log.Println("No issue data found in payload")
		return
	}

	issue := eventPayload.Issue
	log.Printf("Processing issue: %s", issue.Title)

	if strings.Contains(strings.ToLower(issue.Title), "bug") {
		utils.AddLabel(ctx, client, owner, repo, issue.Number, "bug")
		utils.CommentOnIssue(ctx, client, owner, repo, issue.Number, "OwO, It seems that we did a wittle fucky wucky. Sowwy!")
	} else if strings.Contains(strings.ToLower(issue.Title), "feature") {
		utils.AddLabel(ctx, client, owner, repo, issue.Number, "enhancement")
		utils.CommentOnIssue(ctx, client, owner, repo, issue.Number, "This issue seems to be a feature request. Thank you for your suggestion!")
	}

	services.RequestMissingDetails(ctx, client, owner, repo, issue)
	services.SuggestLabelsForIssue(ctx, client, owner, repo, issue)
}
