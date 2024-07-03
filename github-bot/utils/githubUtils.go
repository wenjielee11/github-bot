package utils

import (
    "context"
    "log"

    "github.com/google/go-github/v41/github"
)

func AddLabel(ctx context.Context, client *github.Client, owner, repo string, issueNumber int, label string) {
    _, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, issueNumber, []string{label})
    if err != nil {
        log.Printf("Error adding label to issue #%d: %v", issueNumber, err)
    }
}

func CommentOnIssue(ctx context.Context, client *github.Client, owner, repo string, issueNumber int, comment string) {
    commentRequest := &github.IssueComment{Body: &comment}
    _, _, err := client.Issues.CreateComment(ctx, owner, repo, issueNumber, commentRequest)
    if err != nil {
        log.Printf("Error commenting on issue #%d: %v", issueNumber, err)
    }
}
