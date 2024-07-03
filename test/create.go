package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

const (
	appID          = 934870                            // Replace with your App ID
	installationID = 52398709                          // Replace with your Installation ID
	privateKeyPath = "../issuetracker.private-key.pem"  // Path to your private key
	baseBranch     = "main"                            // Base branch for creating new branches
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

	// Create dummy branches and commits
	createDummyBranchesAndCommits(ctx, client, owner, repo)

	// Create dummy issues
	createDummyIssues(ctx, client, owner, repo)

	// Create dummy pull requests
	createDummyPullRequests(ctx, client, owner, repo)
}

func createDummyBranchesAndCommits(ctx context.Context, client *github.Client, owner, repo string) {
	baseRef, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		log.Fatalf("Error getting base branch reference: %v", err)
	}

	for i := 1; i <= 3; i++ {
		branchName := fmt.Sprintf("dummy-branch-%d", i)

		// Check if the branch already exists
		_, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branchName)
		if err == nil {
			log.Printf("Branch %s already exists, skipping creation", branchName)
		} else {
			// Create a new branch
			ref := &github.Reference{
				Ref:    github.String("refs/heads/" + branchName),
				Object: &github.GitObject{SHA: baseRef.Object.SHA},
			}
			_, _, err := client.Git.CreateRef(ctx, owner, repo, ref)
			if err != nil {
				log.Printf("Error creating branch %s: %v", branchName, err)
				continue
			} else {
				log.Printf("Created branch %s", branchName)
			}
		}

		// Create a dummy commit in the new branch
		commitMessage := fmt.Sprintf("Dummy commit for branch %s", branchName)
		filePath := fmt.Sprintf("dummy-file-%d.txt", i)
		fileContent := fmt.Sprintf("This is a dummy file for branch %s", branchName)
		createDummyCommit(ctx, client, owner, repo, branchName, filePath, fileContent, commitMessage)
	}
}

func createDummyCommit(ctx context.Context, client *github.Client, owner, repo, branchName, filePath, fileContent, commitMessage string) {
	// Get the current branch reference
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branchName)
	if err != nil {
		log.Printf("Error getting reference for branch %s: %v", branchName, err)
		return
	}

	// Get the commit object for the reference
	commit, _, err := client.Git.GetCommit(ctx, owner, repo, ref.Object.GetSHA())
	if err != nil {
		log.Printf("Error getting commit for branch %s: %v", branchName, err)
		return
	}

	// Create a new blob object for the file content
	blob := &github.Blob{
		Content:  github.String(fileContent),
		Encoding: github.String("utf-8"),
	}
	blob, _, err = client.Git.CreateBlob(ctx, owner, repo, blob)
	if err != nil {
		log.Printf("Error creating blob for branch %s: %v", branchName, err)
		return
	}

	// Create a new tree entry for the file
	entry := &github.TreeEntry{
		Path: github.String(filePath),
		Mode: github.String("100644"),
		Type: github.String("blob"),
		SHA:  blob.SHA,
	}

	// Get the current tree
	tree, _, err := client.Git.GetTree(ctx, owner, repo, commit.GetTree().GetSHA(), false)
	if err != nil {
		log.Printf("Error getting tree for branch %s: %v", branchName, err)
		return
	}

	// Create a new tree object with the new entry
	newTree, _, err := client.Git.CreateTree(ctx, owner, repo, *tree.SHA, []*github.TreeEntry{entry})
	if err != nil {
		log.Printf("Error creating tree for branch %s: %v", branchName, err)
		return
	}

	// Create a new commit object
	newCommit := &github.Commit{
		Message: &commitMessage,
		Tree:    newTree,
		Parents: []*github.Commit{commit},
	}
	newCommit, _, err = client.Git.CreateCommit(ctx, owner, repo, newCommit)
	if err != nil {
		log.Printf("Error creating commit for branch %s: %v", branchName, err)
		return
	}

	// Update the branch reference to point to the new commit
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, owner, repo, ref, false)
	if err != nil {
		log.Printf("Error updating reference for branch %s: %v", branchName, err)
	}
}

func createDummyIssues(ctx context.Context, client *github.Client, owner, repo string) {
	for i := 1; i <= 5; i++ {
		title := fmt.Sprintf("[BUG] Dummy issue #%d", i)
		body := fmt.Sprintf("This is the body for dummy issue #%d.", i)
		issueRequest := &github.IssueRequest{
			Title: &title,
			Body:  &body,
		}
		_, _, err := client.Issues.Create(ctx, owner, repo, issueRequest)
		if err != nil {
			log.Printf("Error creating issue #%d: %v", i, err)
		} else {
			log.Printf("Created dummy issue #%d", i)
		}
	}
}

func createDummyPullRequests(ctx context.Context, client *github.Client, owner, repo string) {
	for i := 1; i <= 3; i++ {
		title := fmt.Sprintf("Dummy pull request #%d", i)
		head := fmt.Sprintf("dummy-branch-%d", i)
		body := fmt.Sprintf("This is the body for dummy pull request #%d.", i)
		prRequest := &github.NewPullRequest{
			Title: &title,
			Head:  &head,
			Base:  github.String(baseBranch),
			Body:  &body,
		}
		_, _, err := client.PullRequests.Create(ctx, owner, repo, prRequest)
		if err != nil {
			log.Printf("Error creating pull request #%d: %v", i, err)
		} else {
			log.Printf("Created dummy pull request #%d", i)
		}
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
	return token.GetToken(),nil
}