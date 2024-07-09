package models

// EventPayload represents the payload of a GitHub event.
// It contains the action performed and optional pull request and issue data.
type EventPayload struct {
	Action      string       `json:"action"`       // The action that triggered the event (e.g., "opened", "closed").
	PullRequest *PullRequest `json:"pull_request"` // Pull request data, if applicable.
	Issue       *Issue       `json:"issue"`        // Issue data, if applicable.
}

// PullRequest represents the details of a GitHub pull request.
type PullRequest struct {
	Number       int    `json:"number"`        // The number of the pull request.
	ChangedFiles int    `json:"changed_files"` // The number of files changed in the pull request.
	DiffURL      string `json:"diff_url"`      // The URL to view the diff of the pull request.
}

// Issue represents the details of a GitHub issue.
type Issue struct {
	Number int    `json:"number"` // The number of the issue.
	Body   string `json:"body"`   // The body content of the issue.
	Title  string `json:"title"`  // The title of the issue.
	State  string `json:"state"`  // The state of the issue (e.g., "open", "closed").
}
