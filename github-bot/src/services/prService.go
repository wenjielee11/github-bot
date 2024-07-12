package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/utils"
)

// CheckChangelogUpdated checks if the CHANGELOG.md file is updated in the pull request and provides suggestions if not.
func CheckChangelogUpdated(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, pr *models.PullRequest) {
	// List the files changed in the pull request
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pr.Number, nil)

	if err != nil {
		log.Printf("Error listing files for PR #%d: %v", pr.Number, err)
		return
	}

	updated := false
	var changelogContent string
	var changes strings.Builder

	// Check if the CHANGELOG.md file is updated and collect the changes
	for _, file := range files {
		log.Printf("Processing PR file "+ file.GetFilename())
		if file.GetFilename() == "CHANGELOG.md" {
			updated = true
			changelogContent = file.GetPatch()
		}
		changes.WriteString(fmt.Sprintf("File: %s\n", file.GetFilename()))
		changes.WriteString(fmt.Sprintf("Changes: %s\n\n", file.GetPatch()))
	}

	// Prepare the prompt for suggestions
	var prompt string
	if updated {
		prompt = fmt.Sprintf("Here is the current changelog:\n\n%s\n\nPlease suggest improvements to it based on the following changes:\n\n%s", changelogContent, changes.String())
	} else {
		prompt = fmt.Sprintf("Remind the user to update their CHANGELOG.md file. Please provide suggestions for the changelog based on the following changes:\n\n%s", changes.String())
	}

	// Send the prompt to the LLM for suggestions
	message := map[string]string{
		"PullReqBody": prompt,
	}
	resp, err := AddRow(jamaiClient, models.ActionTable, fmt.Sprintf("%s_%s", owner, repo), message)
	if err != nil {
		log.Fatalf("Error getting changelog suggestions from LLM: %v", err)
	}

	// Read and collect the suggestions from the response
	suggestions, err := readAndCollectContent(resp, "PullReqResponse")
	if err != nil {
		log.Fatalf("Error processing PR %d:\n%v", pr.Number, err)
	}

	// Comment on the pull request with the suggestions
	utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, suggestions)
}

// getCommitDiff fetches the diff of a specific commit.
func getCommitDiff(ctx context.Context, client *github.Client, owner, repo, sha string) (string, error) {
	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, sha, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching commit %s: %w", sha, err)
	}

	var diff strings.Builder
	for _, file := range commit.Files {
		log.Printf("Getting diff for %s", *file.Filename)
		if file.GetPatch() != "" {
			diff.WriteString(fmt.Sprintf("File: %s\n%s\n", file.GetFilename(), file.GetPatch()))
		}
	}

	return diff.String(), nil
}

// CheckSecretKeyLeakage checks for potential secret key leakage using LLM across all commits in a pull request.
func CheckSecretKeyLeakage(ctx context.Context, client *github.Client, jamaiClient *http.Client, owner, repo string, pr *models.PullRequest) {
	// List the commits in the pull request
	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, pr.Number, nil)
	if err != nil {
		log.Printf("Error listing commits for PR #%d: %v", pr.Number, err)
		return
	}

	

	// Check each commit for potential secret key leakage
	for _, commit := range commits {
		var changes strings.Builder
		log.Print("Processing Commit SHA "+ commit.GetSHA())
		diff, err := getCommitDiff(ctx, client, owner, repo, commit.GetSHA())
		if err != nil {
			log.Printf("Error fetching diff for commit %s: %v", commit.GetSHA(), err)
			
		}
		changes.WriteString(fmt.Sprintf("Commit: %s\n", commit.GetSHA()))
		changes.WriteString(fmt.Sprintf("Diff:\n %s", diff))
		log.Printf("Diff of commit %s:\n %s", commit.GetSHA(), diff)
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
		suggestions, err := parseCreatePrSecretResponse(result)
		if err!=nil{
			log.Printf("Error unmarshaling secret response:\n%v", err)
			utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, fmt.Sprintf("Jambo! I had issues checking commit %s for secret leaks. Please contact my developers for more assistance!", commit.GetSHA()))
			continue
		}
		if suggestions.Leak {
			response := fmt.Sprintf("Commit %s:\n%s", suggestions.Commit, suggestions.Response)
			utils.CommentOnIssue(ctx, client, owner, repo, pr.Number, response)
		}
	}
}

// parseCreatePrSecretResponse parses the response into CreatePullReqSecretResponse.
func parseCreatePrSecretResponse(content string) (models.CreatePullReqSecretResponse, error) {
	var prSecretResponse models.CreatePullReqSecretResponse
	if err := json.Unmarshal([]byte(content), &prSecretResponse); err != nil {
		return prSecretResponse, err
	}
	return prSecretResponse, nil
}

// SuggestLabelsForPR suggests labels for a pull request.
// func SuggestLabelsForPR(ctx context.Context, client *github.Client, owner, repo string, pr *models.PullRequest) {
// 	var labels []string
// 	labels = append(labels, "new PR")
// 	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, pr.Number, labels)
// 	if err != nil {
// 		log.Printf("Error adding labels to PR #%d: %v", pr.Number, err)
// 	}
// }
