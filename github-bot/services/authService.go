package services

import (
    "context"
    "github.com/google/go-github/v41/github"
    "golang.org/x/oauth2"
)

func GetInstallationToken(installationID int64, jwtToken string) (string, error) {
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
