package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"encoding/json"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/utils"
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
		prompt = fmt.Sprintf("Remind the user to update their CHANGELOG.md file. Please provide suggestions for the changelog based on the following changes:\n\n%s", changes.String())
	}
	message := map[string]string{
		"PullReqBody": prompt,
	}
	resp, err := AddRow(jamaiClient, models.ActionTable, fmt.Sprintf("%s_%s", owner, repo), message)
	if err != nil {
		log.Fatalf("Error getting changelog suggestions from LLM: %v", err)
	}
	suggestions, err := readAndCollectContent(resp, "PullReqResponse")
	if err != nil {
		log.Fatalf("Error processing PR %d:\n%v", pr.Number, err)
	}
	utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, suggestions)
}

// Function to fetch the diff of a specific commit
func getCommitDiff(ctx context.Context, client *github.Client, owner, repo, sha string) (string, error) {
	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, sha, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching commit %s: %w", sha, err)
	}

	var diff strings.Builder
	for _, file := range commit.Files {
		if file.GetPatch() != "" {
			diff.WriteString(fmt.Sprintf("File: %s\nPatch:\n%s\n\n", file.GetFilename(), file.GetPatch()))
		}
	}

	return diff.String(), nil
}

// Function to check for potential secret key leakage using LLM across all commits in a PR
func CheckSecretKeyLeakage(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, pr *models.PullRequest) {
	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, pr.Number, nil)
	if err != nil {
		log.Printf("Error listing commits for PR #%d: %v", pr.Number, err)
		return
	}

	var changes strings.Builder

	for _, commit := range commits {

		diff, err := getCommitDiff(ctx, client, owner, repo, commit.GetSHA())
		if err != nil {
			log.Printf("Error fetching diff for commit %s: %v", commit.GetSHA(), err)
			continue
		}
		changes.WriteString(fmt.Sprintf("Commit: %s\n", commit.GetSHA()))
		changes.WriteString(fmt.Sprintf("Diff: %s\n\n", diff))
		prompt := changes.String()

		message := map[string]string{
			"PullReqSecretsBody": prompt,
		}

		resp, err := AddRow(jamaiClient, models.ActionTable, fmt.Sprintf("%s_%s", owner, repo), message)
		if err != nil {
			log.Fatalf("Error getting secret key leakage suggestions from LLM: %v", err)
		}
		result, err := readAndCollectContent(resp, "PullReqSecretsResponse")
		if err != nil {
			log.Fatalf("Error processing PR %d:\n%v", pr.Number, err)
		}
		suggestions := parseCreatePrSecretResponse(result)
		if suggestions.Leak {
			response := fmt.Sprintf("Commit %s:\n%s", suggestions.Commit, suggestions.Response)
			utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, response)
		}
	
	}
}

// Function to parse the response into CreateIssueResponse
func parseCreatePrSecretResponse(content string) (models.CreatePullReqSecretResponse) {
	var prSecretResponse models.CreatePullReqSecretResponse
	if err := json.Unmarshal([]byte(content), &prSecretResponse); err != nil {
		log.Fatalf("Error unmarshaling response into CreatePrSecretResponse")
	}
	return prSecretResponse
}


// func SuggestLabelsForPR(ctx context.Context, client *github.Client, owner, repo string, pr *models.PullRequest) {
// 	var labels []string
// 	labels = append(labels, "new PR")
// 	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, pr.Number, labels)
// 	if err != nil {
// 		log.Printf("Error adding labels to PR #%d: %v", pr.Number, err)
// 	}
// }
