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
)

// createRepositoryInGithub는 GitHub API를 사용하여 리포지토리를 생성합니다.
func (service *Service) createRepositoryInGithub(githubRepoDTO *types.CreateGithubRepositoryReq) (*dto.GithubRepositoryRes, error) {
	u, err := service.repository.GithubUser(githubRepoDTO.GithubUserId)
	if err != nil {
		log.Println("깃허브 사용자에 대한 정보를 불러오는데 실패했습니다. : ", "err", err.Error())
		return nil, err
	}

	r := dto.GithubRepository{
		Name:        githubRepoDTO.Name,
		Description: githubRepoDTO.Description,
		Private:     githubRepoDTO.Private,
	}

	jsonR, err := json.Marshal(r)
	if err != nil {
		log.Println("깃허브 리포지토리 생성에 대한 요청 정보를 Json 형태로 변경하지 못했습니다. : ", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", springBootTemplateURL, bytes.NewBuffer(jsonR))
	if err != nil {
		log.Println("깃허브 Rest API 요청 과정에서 에러가 발생했습니다. : ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+u.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
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
