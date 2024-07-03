package models

type EventPayload struct {
    Action      string       `json:"action"`
    PullRequest *PullRequest `json:"pull_request"`
    Issue       *Issue       `json:"issue"`
}

type PullRequest struct {
    Number       int    `json:"number"`
    ChangedFiles int    `json:"changed_files"` // Changed to int
    DiffURL      string `json:"diff_url"`
}

type Issue struct {
    Number int    `json:"number"`
    Body   string `json:"body"`
    Title  string `json:"title"`
    State  string `json:"state"`
}
