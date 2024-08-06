package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

const (
	appID          = 934870                               // Replace with your App ID
	installationID = 52704794                             // Replace with your Installation ID
	privateKeyPath = "../issuetracker.private-key.pem" // Path to your private key
	baseBranch     = "main"                               // Base branch for creating new branches
)

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

	owner := "wenjielee1"  // Replace with your GitHub username
	repo := "chatbot-test" // Replace with your GitHub repository

	// Delete dummy branches
	deleteDummyBranches(ctx, client, owner, repo)

	// Close dummy issues
	closeDummyIssues(ctx, client, owner, repo)

	// Close dummy pull requests
	closeDummyPullRequests(ctx, client, owner, repo)
}

func deleteDummyBranches(ctx context.Context, client *github.Client, owner, repo string) {
	for i := 1; i <= 3; i++ {
		branchName := fmt.Sprintf("dummy-branch-%d", i)

		// Delete the branch
		_, err := client.Git.DeleteRef(ctx, owner, repo, "refs/heads/"+branchName)
		if err != nil {
			log.Printf("Error deleting branch %s: %v", branchName, err)
		} else {
			log.Printf("Deleted branch %s", branchName)
		}
	}
}

func closeDummyIssues(ctx context.Context, client *github.Client, owner, repo string) {
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		log.Fatalf("Error listing issues: %v", err)
	}

	for _, issue := range issues {
			log.Printf("Attempting to close issue #%d with title: %s", issue.GetNumber(), issue.GetTitle())
			state := "closed"
			issueRequest := &github.IssueRequest{
				State: &state,
			}
			_, _, err := client.Issues.Edit(ctx, owner, repo, issue.GetNumber(), issueRequest)
			if err != nil {
				log.Printf("Error closing issue #%d: %v", issue.GetNumber(), err)
			} else {
				log.Printf("Closed issue #%d", issue.GetNumber())
			}

			// Verify the issue is closed
			updatedIssue, _, err := client.Issues.Get(ctx, owner, repo, issue.GetNumber())
			if err != nil {
				log.Printf("Error getting updated issue #%d: %v", issue.GetNumber(), err)
			} else if updatedIssue.GetState() == "closed" {
				log.Printf("Issue #%d successfully closed", issue.GetNumber())
			} else {
				log.Printf("Issue #%d not closed, current state: %s", issue.GetNumber(), updatedIssue.GetState())
			}
	}
}



func closeDummyPullRequests(ctx context.Context, client *github.Client, owner, repo string) {
	prs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		log.Fatalf("Error listing pull requests: %v", err)
	}

	for _, pr := range prs {
		if pr.GetTitle() != "" && pr.GetBody() != "" && isDummyPullRequest(pr.GetTitle(), pr.GetBody()) {
			log.Printf("Attempting to close pull request #%d with title: %s", pr.GetNumber(), pr.GetTitle())
			state := "closed"
			_, _, err := client.PullRequests.Edit(ctx, owner, repo, pr.GetNumber(), &github.PullRequest{
				State: &state,
			})
			if err != nil {
				log.Printf("Error closing pull request #%d: %v", pr.GetNumber(), err)
			} else {
				log.Printf("Closed pull request #%d", pr.GetNumber())
			}

			// Verify the pull request is closed
			updatedPR, _, err := client.PullRequests.Get(ctx, owner, repo, pr.GetNumber())
			if err != nil {
				log.Printf("Error getting updated pull request #%d: %v", pr.GetNumber(), err)
			} else if updatedPR.GetState() == "closed" {
				log.Printf("Pull request #%d successfully closed", pr.GetNumber())
			} else {
				log.Printf("Pull request #%d not closed, current state: %s", pr.GetNumber(), updatedPR.GetState())
			}
		}
	}
}

func isDummyPullRequest(title, body string) bool {
	return strings.HasPrefix(title, "Dummy pull request") && strings.HasPrefix(body, "This is the body for")
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