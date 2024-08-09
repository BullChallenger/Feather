package service

import (
	"feather/repository"
	"feather/types"
	"log"
)

type Service struct {
	repository *repository.Repository
}

// NewService Repository를 주입받아 새로운 Service를 생성합니다.
func NewService(repository *repository.Repository) *Service {
	return &Service{repository: repository}
}

// CreateUser 서비스에 대한 신규 사용자를 생성합니다.
func (service *Service) CreateUser(email string, password string) error {
	err := service.repository.CreateUser(email, password)
	if err != nil {
		log.Println("회원 생성에 실패했습니다. : ", "err", err.Error())
		return err
	}
	return nil
}

// CreateJenkinsUser Jenkins 사용자를 생성합니다.
func (service *Service) CreateJenkinsUser(userId int64, nickname string, token string) error {
	err := service.repository.CreateJenkinsUser(userId, nickname, token)
	if err != nil {
		log.Println("젠킨스 사용자 등록에 실패했습니다. : ", "err", err.Error())
		return err
	}
	return nil
}

// CreateGithubUser GitHub 사용자를 생성합니다.
func (service *Service) CreateGithubUser(userId int64, nickname string, email string, token string) error {
	err := service.repository.CreateGithubUser(userId, nickname, email, token)
	if err != nil {
		log.Println("깃허브 사용자 등록에 실패했습니다. : ", "err", err.Error())
		return err
	}
	return nil
}

// CreateGithubRepository GitHub 리포지토리를 생성하고 생성된 리포지토리에 Jenkins Job을 설정합니다.
func (service *Service) CreateGithubRepository(githubRepoDTO *types.CreateGithubRepositoryReq) error {
	repo, err := service.createRepositoryInGithub(githubRepoDTO)
	if err != nil {
		log.Println("깃허브 리포지토리를 생성하는데 실패했습니다. : ", "err", err.Error())
		return err
	}

	err = service.repository.CreateGithubRepository(githubRepoDTO.GithubUserId, githubRepoDTO.Name, githubRepoDTO.Description, githubRepoDTO.Private)
	if err != nil {
		log.Println("깃허브 리포지토리에 대한 정보를 저장하는데 실패했습니다. : ", "err", err.Error())
		return err
	}

	err = service.createJenkinsJob(repo.Name, repo.Description, repo.HtmlUrl)
	if err != nil {
		log.Println("젠킨스 잡을 생성하는데 실패했습니다. : ", "err", err.Error())
		return err
	}
	return nil
}
