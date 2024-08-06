package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v41/github"
	"github.com/wenjielee1/github-bot/models"
	"github.com/wenjielee1/github-bot/services"
	"github.com/wenjielee1/github-bot/utils"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// HandleGitHubEvents processes GitHub events by reading event data,
// initializing necessary services, and delegating to specific event handlers.
func HandleGitHubEvents(owner, repo, token string) {
	// Create a new context
	ctx := context.Background()

	// Initialize the OAuth2 token source and client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the GitHub event name and path from environment variables
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	log.Printf("Event Name: %s", eventName)
	log.Printf("Event Path: %s", eventPath)

	// Read the event data from the file
	eventData, err := ioutil.ReadFile(eventPath)
	if err != nil {
		log.Fatalf("Error reading event data: %v", err)
	}

	// Unmarshal the event data into an EventPayload struct
	var eventPayload models.EventPayload
	if err := json.Unmarshal(eventData, &eventPayload); err != nil {
		log.Fatalf("Error parsing event data: %v", err)
	}

	// Initialize the JAM.AI client and prepare messages for different event types
	repoLabels := utils.GetLabels(ctx, client, owner, repo)
	var labels []string
	for _, label := range repoLabels {
		if !strings.Contains(*label.Name, "priority") {
			labels = append(labels, *label.Name)
		}
	}
	// Join labels into a single string
	labelsStr := strings.Join(labels, ", ")

	jamaiClient := services.NewJamaiClient(services.GetJamAiHeader())
	actionTableId := owner + "_" + repo +"_"+ utils.GetBotVersion()
	issueResponseMessage := utils.GetColumnMessage("IssueResponse", labelsStr)
	prResponseMessage := utils.GetColumnMessage("PullReqResponse", labelsStr)
	prSecretsMessage := utils.GetColumnMessage("PullReqSecretsResponse", labelsStr)
	secretsJSONMessage := utils.GetColumnMessage("SecretsJSONResponse", labelsStr)
	// Define the agents and their respective messages
	agents := []models.Agent{
		{ColumnID: "IssueBody", Messages: nil},
		{ColumnID: "PullReqBody", Messages: nil},
		{ColumnID: "PullReqSecretsBody", Messages: nil},
		{ColumnID: "IssueResponse", Messages: issueResponseMessage},
		{ColumnID: "PullReqResponse", Messages: prResponseMessage},
		{ColumnID: "PullReqSecretsResponse", Messages: prSecretsMessage},
		{ColumnID: "SecretsJSONResponse", Messages: secretsJSONMessage},
	}
	
	// Create a table in the JAM.AI client with the defined agents
	services.CreateTable(jamaiClient, models.ActionTable, actionTableId, agents)

	// Handle specific GitHub events
	switch eventName {
	case "issues":
		HandleIssueEvent(ctx, client, jamaiClient, owner, repo, eventPayload)
	case "pull_request":
		HandlePullRequestEvent(ctx, client, jamaiClient, owner, repo, eventPayload)
	default:
		log.Printf("Unhandled event: %s", eventName)
	}
}
