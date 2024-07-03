package services

import (
	"context"
	"log"

	"github.com/wenjielee1/github-bot/utils"

	"github.com/wenjielee1/github-bot/models"

	"github.com/google/go-github/v41/github"
)

func RequestMissingDetails(ctx context.Context, client *github.Client, owner, repo string, issue *models.Issue) {
	if len(issue.Body) < 10 {
		utils.CommentOnIssue(ctx, client, owner, repo, issue.Number, "Please provide more details for this issue.")
	}
}

func SuggestLabelsForIssue(ctx context.Context, client *github.Client, owner, repo string, issue *models.Issue) {
	var labels []string
	if issue.State == "open" {
		labels = append(labels, "new issue")
		_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, issue.Number, labels)
		if err != nil {
			log.Printf("Error adding labels to issue #%d: %v", issue.Number, err)
		}
	}
}
