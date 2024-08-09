package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"feather/types"
	"feather/types/dto"
	"log"
	"net/http"
)

const (
	springBootTemplateURL = "https://api.github.com/repos/BullChallenger/spring-boot-template/generate"
	jenkinsWebhookURL     = "https://jks.dev-in-wonderland.pro/github-webhook/"
)

// createRepositoryInGithub GitHub API를 사용하여 리포지토리를 생성합니다.
func (service *Service) createRepositoryInGithub(req *types.CreateGithubRepositoryReq) (*dto.GithubRepositoryRes, error) {
	u, err := service.repository.GithubUser(req.GithubUserId)
	if err != nil {
		log.Println("깃허브 사용자에 대한 정보를 불러오는데 실패했습니다. : ", "err", err.Error())
		return nil, err
	}

	r := dto.GithubRepository{
		Name:        req.Name,
		Description: req.Description,
		Private:     req.Private,
	}

	jsonR, err := json.Marshal(r)
	if err != nil {
		log.Println("깃허브 리포지토리 생성에 대한 요청 정보를 Json 형태로 변경하지 못했습니다. : ", err)
		return nil, err
	}

	githubReq, err := http.NewRequest("POST", springBootTemplateURL, bytes.NewBuffer(jsonR))
	if err != nil {
		log.Println("깃허브 Rest API 요청 과정에서 에러가 발생했습니다. : ", err)
		return nil, err
	}

	githubReq.Header.Set("Content-Type", "application/json")
	githubReq.Header.Set("Authorization", "Bearer "+u.Token)

	client := &http.Client{}
	resp, err := client.Do(githubReq)
	if err != nil {
		log.Println("깃허브 Rest API 요청 과정에서 에러가 발생했습니다. : ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("깃허브 Rest API 요청을 통해 리포지토리를 생성하는데 실패했습니다. : %s\n", resp.Status)
		return nil, errors.New("깃허브 Rest API 요청을 통해 리포지토리를 생성하는데 실패했습니다")
	}

	var repo dto.GithubRepositoryRes
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		log.Println("Json 데이터를 디코딩하는 과정에서 에러가 발생했습니다. : ", err)
		return nil, err
	}

	err = service.createWebhook(&repo, u.Token)
	if err != nil {
		log.Println("젠킨스 웹훅 생성에 실패했습니다. : ", err)
		return nil, err
	}

	log.Println("깃허브 리포지토리가 성공적으로 생성되었습니다!")
	return &repo, nil
}

// createWebhook GitHub 리포지토리에 Webhook을 생성합니다.
func (service *Service) createWebhook(githubRepo *dto.GithubRepositoryRes, token string) error {
	type config struct {
		Url         string `json:"url"`
		ContentType string `json:"content_type"`
		InsecureSsl string `json:"insecure_ssl"`
	}

	type webhookReq struct {
		Name   string   `json:"name"`
		Active bool     `json:"active"`
		Event  []string `json:"event"`
		Config config   `json:"config"`
	}

	w := &webhookReq{
		Name:   "web",
		Active: true,
		Event:  []string{"push"},
		Config: config{
			Url:         jenkinsWebhookURL,
			ContentType: "json",
			InsecureSsl: "0",
		},
	}

	jsonW, err := json.Marshal(w)
	if err != nil {
		log.Println("젠킨스 웹훅에 대한 요청 정보를 Json 형태로 변경하지 못했습니다. : ", err)
		return err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/repos/"+githubRepo.FullName+"/hooks", bytes.NewBuffer(jsonW))
	if err != nil {
		log.Println("젠킨스 Rest API 요청 과정에서 에러가 발생했습니다. :", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("젠킨스 Rest API 요청 과정에서 에러가 발생했습니다. : ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("젠킨스 웹훅을 생성하는데 실패했습니다. : %s\n", resp.Status)
		return errors.New("젠킨스 웹훅을 생성하는데 실패했습니다")
	}

	log.Println("젠킨스 웹훅이 성공적으로 생성되었습니다!")
	return nil
}
