package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"feather/types"
	"feather/types/dto"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	jenkinsHostURL     = "https://jks.dev-in-wonderland.pro"
	jenkinsWebhookURL  = "https://jks.dev-in-wonderland.pro/github-webhook/"
	jenkinsXMLFilePath = "config/feather_jenkins_job.xml"
	jenkinsUser        = "cheshire-cat"
	jenkinsToken       = "11e94137b6bde1652311303fe9d57745e3"
)

// createWebhook는 GitHub 리포지토리에 Webhook을 생성합니다.
func (service *Service) createWebhook(githubRepository *dto.GithubRepositoryRes, token string) error {
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

	req, err := http.NewRequest("POST", "https://api.github.com/repos/"+githubRepository.FullName+"/hooks", bytes.NewBuffer(jsonW))
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

// createJenkinsJob은 Jenkins에 새로운 Job을 생성합니다.
func (service *Service) createJenkinsJob(jobName string, jobDescription string, githubURL string) error {
	xmlTemplate, err := os.ReadFile(jenkinsXMLFilePath)
	if err != nil {
		log.Println("XML 파일을 읽을 수 없습니다: ", err)
		return err
	}

	var config types.FlowDefinition
	err = xml.Unmarshal(xmlTemplate, &config)
	if err != nil {
		log.Println("XML 파싱 중 오류 발생: ", err)
		return err
	}

	config.Description = jobDescription
	config.Properties.GithubProjectProperty.ProjectUrl = githubURL
	config.Definition.Scm.UserRemoteConfigs.UserRemoteConfig.URL = githubURL

	modifiedXML, err := xml.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Println("XML 변환 중 오류 발생: ", err)
		return err
	}

	client := &http.Client{}
	createJobURL := fmt.Sprintf("%s/createItem?name=%s", jenkinsHostURL, jobName)
	req, err := http.NewRequest("POST", createJobURL, bytes.NewBuffer(modifiedXML))
	if err != nil {
		log.Println("요청 생성 중 오류: ", err)
		return err
	}

	req.SetBasicAuth(jenkinsUser, jenkinsToken)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("요청 전송 중 오류: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Job 생성 실패: %s\n%s", resp.Status, string(bodyBytes))
		return errors.New("failed to create Jenkins job")
	}

	log.Println("Job이 성공적으로 생성되었습니다!")
	return nil
}
