package services

import (
	"context"
	"log"
	"net/http"

	"github.com/wenjielee1/github-bot/utils"

	"github.com/wenjielee1/github-bot/models"

	"github.com/google/go-github/v41/github"
)

func ProcessIssue(ctx context.Context, client *github.Client, jamaiClient *http.Client, tableId string, owner, repo string, issue *models.Issue) {

	message := map[string]string{
		"IssueBody": issue.Title + "\n" + issue.Body,
	}
	resp, err := AddRow(jamaiClient, models.ActionTable, tableId, message)
	if err != nil {
		log.Fatalf("Error processing issue %d %s:\n%v", issue.Number, issue.Title, err)
	}
	respString, err := readAndCollectContent(resp, "IssueResponse")
	if err != nil {
		log.Fatalf("Error processing issue %d %s:\n%v", issue.Number, issue.Title, err)
	}
	result, err := parseCreateIssueResponse(respString)
	if err != nil {
		log.Fatalf("Error parsing create issue response: %v", err)
	}

	labels := append(result.Labels, "priority: "+result.Priority)
	utils.CreatePriorityLabels(ctx, client, owner, repo)
	utils.AddLabels(ctx, client, owner, repo, issue.Number, labels)
	utils.CommentOnIssue(ctx, client, owner, repo, issue.Number, result.Response)

}
