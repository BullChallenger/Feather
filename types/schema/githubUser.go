package schema

type GithubUser struct {
	ID       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
