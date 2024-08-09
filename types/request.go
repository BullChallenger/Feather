package types

type CreateUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateGithubUserReq struct {
	UserId   int64  `json:"user_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type CreateJenkinsUserReq struct {
	UserId   int64  `json:"user_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type CreateGithubRepositoryReq struct {
	UserId       int64  `json:"user_id" binding:"required"`
	GithubUserId int64  `json:"github_user_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Private      bool   `json:"private"`
}
