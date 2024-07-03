package services

import (
	"context"
	"github-bot/models"
	"github-bot/utils"
	"log"
	"strings"

	"github.com/google/go-github/v41/github"
)


func CheckChangelogUpdated(ctx context.Context, client *github.Client, owner, repo string, pr *models.PullRequest) {
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pr.Number, nil)
	if err != nil {
		log.Printf("Error listing files for PR #%d: %v", pr.Number, err)
		return
	}

	updated := false
	for _, file := range files {
		if file.GetFilename() == "CHANGELOG.md" {
			updated = true
			break
		}
	}

	if !updated {
		utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, "Please update the CHANGELOG.md file.")
	}
}

func CheckSecretKeyLeakage(ctx context.Context, client *github.Client, owner, repo string, pr *models.PullRequest) {
	opts := github.RawOptions{Type: github.Diff}
	diff, _, err := client.PullRequests.GetRaw(ctx, owner, repo, pr.Number, opts)
	if err != nil {
		log.Printf("Error fetching diff for PR #%d: %v", pr.Number, err)
		return
	}

	if strings.Contains(diff, "secret") {
		utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, "Potential secret key leakage detected.")
	}
}

func SuggestLabelsForPR(ctx context.Context, client *github.Client, owner, repo string, pr *models.PullRequest) {
	var labels []string
	labels = append(labels, "new PR")
	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, pr.Number, labels)
	if err != nil {
		log.Printf("Error adding labels to PR #%d: %v", pr.Number, err)
	}
}
