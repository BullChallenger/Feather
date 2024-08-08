package schema

type JenkinsUser struct {
	ID       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}
