package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wenjielee1/github-bot/utils"

	"github.com/wenjielee1/github-bot/models"

	"github.com/google/go-github/v41/github"
)


func CheckChangelogUpdated(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, pr *models.PullRequest) {
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pr.Number, nil)
	if err != nil {
		log.Printf("Error listing files for PR #%d: %v", pr.Number, err)
		return
	}

	updated := false
	var changelogContent string
	var changes strings.Builder

	for _, file := range files {
		if file.GetFilename() == "CHANGELOG.md" {
			updated = true
			changelogContent = file.GetPatch()
		}
		changes.WriteString(fmt.Sprintf("File: %s\n", file.GetFilename()))
		changes.WriteString(fmt.Sprintf("Changes: %s\n\n", file.GetPatch()))
	}

	var prompt string
	if updated {
		prompt = fmt.Sprintf("Here is the current changelog:\n\n%s\n\nPlease suggest improvements to it based on the following changes:\n\n%s", changelogContent, changes.String())
	} else {
		prompt = fmt.Sprintf("Please provide suggestions for the changelog based on the following changes:\n\n%s", changes.String())
		utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, "Please update the CHANGELOG.md file.")
	}
	message:= map[string]string{
		"PullReqBody": prompt,
	}
	resp, err := AddRow(jamaiClient, models.ActionTable, fmt.Sprintf("%s_%s", owner,repo), message)
	if err != nil {
		log.Fatalf("Error getting changelog suggestions from LLM: %v", err)
	}
	suggestions, err:= readAndCollectContent(resp, "PullReqResponse")
	if err != nil {
		log.Fatalf("Error processing PR %d:\n%v", pr.Number, err)
	}
	utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, suggestions)
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
