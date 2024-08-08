package schema

import "time"

type GithubRepository struct {
	ID           int64     `json:"id"`
	GithubUserId int64     `json:"github_user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Private      bool      `json:"private"`
	HtmlUrl      string    `json:"html_url"`
	CreateAt     time.Time `json:"create_at"`
	UpdatedAt    time.Time `json:"updated-at"`
}
