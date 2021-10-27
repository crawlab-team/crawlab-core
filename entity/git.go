package entity

type GitPayload struct {
	Paths         []string `json:"paths"`
	CommitMessage string   `json:"commit_message"`
}
