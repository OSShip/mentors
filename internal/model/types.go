package model

import "encoding/json"

type Application struct {
	ID                      string          `json:"id"`
	UserID                  string          `json:"user_id"`
	Status                  string          `json:"status"`
	GithubData              json.RawMessage `json:"github_data,omitempty"`
	CreatedAt               string          `json:"created_at"`
	ApplicantEmail          string          `json:"applicant_email,omitempty"`
	ApplicantDisplayName    string          `json:"applicant_display_name,omitempty"`
	ApplicantGithubUsername string          `json:"applicant_github_username,omitempty"`
}
