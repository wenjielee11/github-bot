package services

import (
	"context"
	"log"
	"net/http"

	"encoding/json"
	"github.com/wenjielee1/github-bot/utils"

	"github.com/wenjielee1/github-bot/models"

	"github.com/google/go-github/v41/github"
)


func ProcessIssue(ctx context.Context, client *github.Client, jamaiClient *http.Client, tableId string,owner, repo string, issue *models.Issue){
	
	message:= map[string]string{
		"IssueBody": issue.Title+"\n"+issue.Body,
	}
	result, err := AddRow(jamaiClient, models.ActionTable, tableId, message)
	if err != nil {
		log.Printf("Error processing issue %d %s:\n%v", issue.Number, issue.Title, err)
	}

	var response models.CreateIssueResponse
	if err := json.Unmarshal(result, &response); err != nil {
		log.Printf("error unmarshaling JSON: %v", err)
	}

	labels := append(response.Labels, "priority: "+ response.Priority)
	utils.CreatePriorityLabels(ctx, client, owner,repo)
	utils.AddLabels(ctx, client, owner, repo, issue.Number, labels)
	utils.CommentOnIssue(ctx, client, owner, repo, issue.Number, response.Response)
}



