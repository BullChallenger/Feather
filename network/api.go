package network

import (
	"feather/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	api := &api{server: server}
	server.engine.POST("/api/users/create", api.createUser)
	server.engine.POST("/api/github_users/create", api.createGithubUser)
	server.engine.POST("/api/github_repo/create", api.createGithubRepository)
	server.engine.POST("/api/jenkins_users/create", api.createJenkinsUser)
}

func (api *api) createUser(ctx *gin.Context) {
	var req types.CreateUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := api.server.service.CreateUser(req.Email, req.Password); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (api *api) createGithubUser(ctx *gin.Context) {
	var req types.CreateGithubUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := api.server.service.CreateGithubUser(req.UserId, req.Nickname, req.Email, req.Token); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (api *api) createJenkinsUser(ctx *gin.Context) {
	var req types.CreateJenkinsUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := api.server.service.CreateJenkinsUser(req.UserId, req.Nickname, req.Token); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (api *api) createGithubRepository(ctx *gin.Context) {
	var req *types.CreateGithubRepositoryReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := api.server.service.CreateGithubRepository(req); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}
