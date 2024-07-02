package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

const (
	appID          = 934870                            // Replace with your App ID
	installationID = 52398709                          // Replace with your Installation ID
	privateKeyPath = "./issuetracker.private-key.pem"  // Path to your private key
)

type PullRequest struct {
	Number       int      `json:"number"`
	ChangedFiles []string `json:"changed_files"`
	DiffURL      string   `json:"diff_url"`
}

type Issue struct {
	Number int    `json:"number"`
	Body   string `json:"body"`
	Title  string `json:"title"`
	State  string `json:"state"`
}

type EventPayload struct {
	Action      string       `json:"action"`
	PullRequest *PullRequest `json:"pull_request"`
	Issue       *Issue       `json:"issue"`
}

func main() {
	// Load the private key
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	// Parse the private key
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}

	// Generate a JWT
	jwtToken, err := generateJWT(appID, parsedKey)
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}

	// Get the installation token
	installationToken, err := getInstallationToken(installationID, jwtToken)
	if err != nil {
		log.Fatalf("Error getting installation token: %v", err)
	}

	// Use the installation token to interact with the GitHub API
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: installationToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Read GitHub event data
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	eventData, err := ioutil.ReadFile(eventPath)
	if err != nil {
		log.Fatalf("Error reading event data: %v", err)
	}

	var eventPayload EventPayload
	if err := json.Unmarshal(eventData, &eventPayload); err != nil {
		log.Fatalf("Error parsing event data: %v", err)
	}

	// Route events
	switch eventName {
	case "issues":
		handleIssueEvent(ctx, client, "wenjielee1", "chatbot-test", eventPayload)
	case "pull_request":
		handlePullRequestEvent(ctx, client, "wenjielee1", "chatbot-test", eventPayload)
	default:
		fmt.Printf("Unhandled event: %s\n", eventName)
	}
}

func handleIssueEvent(ctx context.Context, client *github.Client, owner, repo string, eventPayload EventPayload) {
	if eventPayload.Issue == nil {
		log.Println("No issue data found in payload")
		return
	}

	issue := eventPayload.Issue
	fmt.Printf("Processing issue: %s\n", issue.Title)

	// Simple categorization example
	if strings.Contains(strings.ToLower(issue.Title), "bug") {
		addLabel(ctx, client, owner, repo, issue.Number, "bug")
		commentOnIssue(ctx, client, owner, repo, issue.Number, "OwO, It seems that we did a wittle fucky wucky. Sowwy!")
	} else if strings.Contains(strings.ToLower(issue.Title), "feature") {
		addLabel(ctx, client, owner, repo, issue.Number, "enhancement")
		commentOnIssue(ctx, client, owner, repo, issue.Number, "This issue seems to be a feature request. Thank you for your suggestion!")
	}

	// Request missing details
	requestMissingDetails(ctx, client, owner, repo, issue)
	// Suggest labels
	suggestLabelsForIssue(ctx, client, owner, repo, issue)
}

func handlePullRequestEvent(ctx context.Context, client *github.Client, owner, repo string, eventPayload EventPayload) {
	if eventPayload.PullRequest == nil {
		log.Println("No pull request data found in payload")
		return
	}

	pr := eventPayload.PullRequest
	fmt.Printf("Processing pull request: #%d\n", pr.Number)

	// Check if CHANGELOG is updated
	checkChangelogUpdated(ctx, client, owner, repo, pr)
	// Check for secret key leakage
	checkSecretKeyLeakage(ctx, client, owner, repo, pr)
	// Suggest labels
	suggestLabelsForPR(ctx, client, owner, repo, pr)
}



func checkChangelogUpdated(ctx context.Context, client *github.Client, owner, repo string, pr *PullRequest) {
	
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
		commentOnIssue(ctx, client, owner, repo, pr.Number, "Please update the CHANGELOG.md file.")
	}
}

func checkSecretKeyLeakage(ctx context.Context, client *github.Client, owner, repo string, pr *PullRequest) {
	opts := github.RawOptions{Type: github.Diff}
	diff, _, err := client.PullRequests.GetRaw(ctx, owner, repo, pr.Number, opts)
	if err != nil {
		log.Printf("Error fetching diff for PR #%d: %v", pr.Number, err)
		return
	}

	if strings.Contains(diff, "secret") {
		commentOnIssue(ctx, client, owner, repo, pr.Number, "Potential secret key leakage detected.")
	}
}

func requestMissingDetails(ctx context.Context, client *github.Client, owner, repo string, issue *Issue) {
	if len(issue.Body) < 10 {
		commentOnIssue(ctx, client, owner, repo, issue.Number, "Please provide more details for this issue.")
	}
}

func suggestLabelsForIssue(ctx context.Context, client *github.Client, owner, repo string, issue *Issue) {
	var labels []string
	if issue.State == "open" {
		labels = append(labels, "new issue")
		_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, issue.Number, labels)
		if err != nil {
			log.Printf("Error adding labels to issue #%d: %v", issue.Number, err)
		}
	}
}

func suggestLabelsForPR(ctx context.Context, client *github.Client, owner, repo string, pr *PullRequest) {
	var labels []string
	labels = append(labels, "new PR")
	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, pr.Number, labels)
	if err != nil {
		log.Printf("Error adding labels to PR #%d: %v", pr.Number, err)
	}
}

func addLabel(ctx context.Context, client *github.Client, owner, repo string, issueNumber int, label string) {
	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, issueNumber, []string{label})
	if err != nil {
		log.Printf("Error adding label to issue #%d: %v", issueNumber, err)
	}
}

func commentOnIssue(ctx context.Context, client *github.Client, owner, repo string, issueNumber int, comment string) {
	commentRequest := &github.IssueComment{Body: &comment}
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, issueNumber, commentRequest)
	if err != nil {
		log.Printf("Error commenting on issue #%d: %v", issueNumber, err)
	}
}

func generateJWT(appID int64, key *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute * 10).Unix(),
		Issuer:    fmt.Sprintf("%d", appID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func getInstallationToken(installationID int64, jwtToken string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: jwtToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", err
	}
	return token.GetToken(), nil
}
