package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"feather/repository"
	"feather/types"
	"feather/types/dto"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Service struct {
	repository *repository.Repository
}

const (
	jenkinsHostURL        = "https://jks.dev-in-wonderland.pro"
	jenkinsWebhookURL     = "https://jks.dev-in-wonderland.pro/github-webhook/"
	jenkinsXMLFilePath    = "config/feather_jenkins_job.xml"
	jenkinsUser           = "cheshire-cat"
	jenkinsToken          = "11e94137b6bde1652311303fe9d57745e3"
	springBootTemplateURL = "https://api.github.com/repos/BullChallenger/spring-boot-template/generate"
)

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}
	return s
}

func (service *Service) CreateUser(email string, password string) error {
	if err := service.repository.CreateUser(email, password); err != nil {
		log.Println("Failed to Create User", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (service *Service) CreateJenkinsUser(userId int64, nickname string, token string) error {
	if err := service.repository.CreateJenkinsUser(userId, nickname, token); err != nil {
		log.Println("Failed to Create Github User", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (service *Service) CreateGithubUser(userId int64, nickname string, email string, token string) error {
	if err := service.repository.CreateGithubUser(userId, nickname, email, token); err != nil {
		log.Println("Failed to Create Github User", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (service *Service) CreateGithubRepository(githubRepoDTO *types.CreateGithubRepositoryReq) error {
	var repo *dto.GithubRepositoryRes

	if resp, err := service.createRepositoryInGithub(githubRepoDTO); err != nil {
		log.Println("Failed to Create Github Repository with Github API", "err", err.Error())
		return err
	} else {
		repo = resp
	}

	if err := service.repository.CreateGithubRepository(githubRepoDTO.GithubUserId, githubRepoDTO.Name, githubRepoDTO.Description, githubRepoDTO.Private); err != nil {
		log.Println("Failed to Create Github Repository", "err", err.Error())
		return err
	}

	if err := service.createJenkinsJob(repo.Name, repo.Description, repo.HtmlUrl); err != nil {
		log.Println("Failed to Create Jenkins Job", "err", err.Error())
		return err
	}
	return nil
}

func (service *Service) createRepositoryInGithub(githubRepoDTO *types.CreateGithubRepositoryReq) (*dto.GithubRepositoryRes, error) {
	if u, err := service.repository.GithubUser(githubRepoDTO.GithubUserId); err != nil {
		log.Println("Failed to Get Github User", "err", err.Error())
		return nil, err
	} else {
		r := dto.GithubRepository{
			Name:        githubRepoDTO.Name,
			Description: githubRepoDTO.Description,
			Private:     githubRepoDTO.Private,
		}

		jsonR, err := json.Marshal(r)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)
			return nil, err
		}
		req, err := http.NewRequest("POST", springBootTemplateURL, bytes.NewBuffer(jsonR))
		if err != nil {
			fmt.Println("Error creating request: ", err)
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+u.Token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Failed to create repository: %s\n", resp.Status)
			return nil, errors.New("Failed to create repository")
		}

		var repo *dto.GithubRepositoryRes
		if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
		}

		service.createWebhook(repo, u.Token)

		fmt.Println("Repository created successfully!")
		return repo, nil
	}
}

func (service *Service) createWebhook(githubRepository *dto.GithubRepositoryRes, token string) error {
	/**
	{
		"name":"web",  default: web
		"active":true, default: true
		"events":["push","pull_request"], default: push
		"config":{
			"url":"https://example.com/webhook",
			"content_type":"json", default: form
			"insecure_ssl":"0" default: 0 => verification is performed
		}
	}
	*/
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
		log.Println("Error marshalling JSON: ", err)
		return err
	}
	req, err := http.NewRequest("POST", "https://api.github.com/repos/"+githubRepository.FullName+"/hooks", bytes.NewBuffer(jsonW))
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Failed to create webhook: %s\n", resp.Status)
		return errors.New("Failed to create webhook")
	}

	fmt.Println("Webhook created successfully!")
	return nil
}

func (service *Service) createJenkinsJob(jobName string, jobDescription string, githubURL string) error {

	type GithubProjectProperty struct {
		Plugin      string `xml:"plugin,attr"`
		ProjectUrl  string `xml:"projectUrl"`
		DisplayName string `xml:"displayName"`
	}

	type DurabilityHintJobProperty struct {
		Hint string `xml:"hint"`
	}

	type GitHubPushTrigger struct {
		Plugin string `xml:"plugin,attr"`
		Spec   string `xml:"spec"`
	}

	type Triggers struct {
		GitHubPushTrigger GitHubPushTrigger `xml:"com.cloudbees.jenkins.GitHubPushTrigger"`
	}

	type PipelineTriggersJobProperty struct {
		Triggers Triggers `xml:"triggers"`
	}

	type DisableConcurrentBuildsJobProperty struct {
		AbortPrevious bool `xml:"abortPrevious"`
	}

	type Properties struct {
		DisableConcurrentBuildsJobProperty DisableConcurrentBuildsJobProperty `xml:"org.jenkinsci.plugins.workflow.job.properties.DisableConcurrentBuildsJobProperty"`
		GithubProjectProperty              GithubProjectProperty              `xml:"com.coravy.hudson.plugins.github.GithubProjectProperty"`
		DurabilityHintJobProperty          DurabilityHintJobProperty          `xml:"org.jenkinsci.plugins.workflow.job.properties.DurabilityHintJobProperty"`
		PipelineTriggersJobProperty        PipelineTriggersJobProperty        `xml:"org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty"`
	}

	type BranchSpec struct {
		Name string `xml:"name"`
	}

	type Branches struct {
		BranchSpec BranchSpec `xml:"hudson.plugins.git.BranchSpec"`
	}

	type SubmoduleCfg struct {
		Class string `xml:"class,attr"`
	}

	type UserRemoteConfig struct {
		URL           string `xml:"url"`
		CredentialsId string `xml:"credentialsId"`
	}

	type UserRemoteConfigs struct {
		UserRemoteConfig UserRemoteConfig `xml:"hudson.plugins.git.UserRemoteConfig"`
	}

	type Scm struct {
		Class                             string            `xml:"class,attr"`
		Plugin                            string            `xml:"plugin,attr"`
		ConfigVersion                     int               `xml:"configVersion"`
		UserRemoteConfigs                 UserRemoteConfigs `xml:"userRemoteConfigs"`
		Branches                          Branches          `xml:"branches"`
		DoGenerateSubmoduleConfigurations bool              `xml:"doGenerateSubmoduleConfigurations"`
		SubmoduleCfg                      SubmoduleCfg      `xml:"submoduleCfg"`
	}

	type Definition struct {
		Class       string `xml:"class,attr"`
		Plugin      string `xml:"plugin,attr"`
		Scm         Scm    `xml:"scm"`
		ScriptPath  string `xml:"scriptPath"`
		Lightweight bool   `xml:"lightweight"`
	}

	type FlowDefinition struct {
		XMLName          xml.Name   `xml:"flow-definition"`
		Plugin           string     `xml:"plugin,attr"`
		Description      string     `xml:"description"`
		KeepDependencies bool       `xml:"keepDependencies"`
		Properties       Properties `xml:"properties"`
		Definition       Definition `xml:"definition"`
		Disabled         bool       `xml:"disabled"`
	}

	xmlTemplate, err := os.ReadFile(jenkinsXMLFilePath)
	if err != nil {
		log.Fatalf("XML 파일을 읽을 수 없습니다: %v", err)
		return err
	}

	var config FlowDefinition
	err = xml.Unmarshal(xmlTemplate, &config)
	if err != nil {
		fmt.Printf("XML 파싱 중 오류 발생: %v\n", err)
		return err
	}

	config.Description = jobDescription
	config.Properties.GithubProjectProperty.ProjectUrl = githubURL
	config.Definition.Scm.UserRemoteConfigs.UserRemoteConfig.URL = githubURL

	modifiedXML, err := xml.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Printf("XML 변환 중 오류 발생: %v\n", err)
		return err
	}

	client := &http.Client{}

	createJobURL := fmt.Sprintf("%s/createItem?name=%s", jenkinsHostURL, jobName)
	req, err := http.NewRequest("POST", createJobURL, bytes.NewBuffer(modifiedXML))
	if err != nil {
		log.Fatalf("요청 생성 중 오류: %v", err)
		return err
	}
	req.SetBasicAuth(jenkinsUser, jenkinsToken)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("요청 전송 중 오류: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		fmt.Println("Job이 성공적으로 생성되었습니다!")
		return nil
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("Job 생성 실패: %s\n%s", resp.Status, string(bodyBytes))
		return err
	}
}
