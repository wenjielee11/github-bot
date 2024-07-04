package handlers

import (
	"context"
	"log"

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

	jamaiClient := services.NewJamaiClient(services.GetJamAiHeader())
	actionTableId := owner+"_"+repo 
	issueResponseMessage := utils.GetColumnMessage("IssueResponse")
	prResponseMessage := utils.GetColumnMessage("PullReqResponse")
	agents:= []models.Agent{
		{ColumnID: "IssueBody", Messages: nil},
		{ColumnID: "PullReqBody", Messages: nil},
		{ColumnID: "IssueResponse", Messages: issueResponseMessage},
		{ColumnID: "PullReqResponse", Messages: prResponseMessage},
	}

	services.CreateTable(jamaiClient, models.ActionTable, actionTableId, agents)
	services.ProcessIssue(ctx, client, jamaiClient, actionTableId, owner, repo, issue)
	
}
