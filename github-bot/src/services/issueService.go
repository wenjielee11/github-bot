package services

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/utils"
)

// ProcessIssue processes a GitHub issue by adding its details to a table,
// reading the response, and updating the issue with labels and comments.
func ProcessIssue(ctx context.Context, client *github.Client, jamaiClient *http.Client, tableId string, owner, repo string, issue *models.Issue) {
	// Create a message map with the issue title and body
	message := map[string]string{
		"IssueBody": issue.Title + "\n" + issue.Body,
	}
	
	// Add the issue details to the table and get the response
	resp, err := AddRow(jamaiClient, models.ActionTable, tableId, message)
	if err != nil {
		log.Fatalf("Error processing issue %d %s:\n%v", issue.Number, issue.Title, err)
	}

	// Read and collect the response content for "IssueResponse"
	respString, err := readAndCollectContent(resp, "IssueResponse")
	if err != nil {
		log.Fatalf("Error processing issue %d %s:\n%v", issue.Number, issue.Title, err)
	}

	// Parse the create issue response
	result, err := parseCreateIssueResponse(respString)
	if err != nil {
		log.Fatalf("Error parsing create issue response: %v", err)
	}
	
	// Append priority label to the result labels
	labels := append(result.Labels, "priority: "+result.Priority)

	LabelIssue(ctx, client, jamaiClient, tableId, owner, repo, issue, labels)

	// Comment on the issue with the response. Disabled for now as of 16/7
	// utils.CommentOnIssue(ctx, client, owner, repo, issue.Number, result.Response)
}

func LabelIssue(ctx context.Context, client *github.Client, jamaiClient *http.Client, tableId string, owner, repo string, issue *models.Issue, labels []string){
	currentLabels, _, err:= client.Issues.ListLabelsByIssue(ctx, owner, repo, issue.Number, nil) 
	if err != nil {
        log.Fatalf("Error retrieving labels: %v", err)
    }
	for _, label:= range currentLabels{
		
		if strings.Contains(*label.Name, "priority") || strings.Contains(*label.Name, "status") {
            log.Printf("Found label with prefix priority or status, skipping label: "+ *label.Name)
            return
        }
	}
	// // Create priority labels in the repository if they do not exist
	// utils.CreatePriorityLabels(ctx, client, owner, repo)

	// Add labels to the issue
	utils.AddLabels(ctx, client, owner, repo, issue.Number, labels)
}
