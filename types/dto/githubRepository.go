package dto

type GithubRepository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type GithubRepositoryRes struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HtmlUrl     string `json:"html_url"`
}
