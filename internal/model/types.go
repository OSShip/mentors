package model

import "encoding/json"

type Application struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	Status     string          `json:"status"`
	GithubData json.RawMessage `json:"github_data,omitempty"`
	CreatedAt  string          `json:"created_at"`
}
