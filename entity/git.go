package entity

type GitPayload struct {
	Paths         []string `json:"paths"`
	CommitMessage string   `json:"commit_message"`
}

type GitConfig struct {
	Url string `json:"url" bson:"url"`
}
