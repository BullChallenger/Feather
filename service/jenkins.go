package service

import (
	"bytes"
	"encoding/xml"
	"errors"
	"feather/types"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	jenkinsHostURL     = "https://jks.dev-in-wonderland.pro"
	jenkinsXMLFilePath = "config/feather_jenkins_job.xml"
	jenkinsUser        = "cheshire-cat"
	jenkinsToken       = "11e94137b6bde1652311303fe9d57745e3"
)

// createJenkinsJob Jenkins에 새로운 Job을 생성합니다.
func (service *Service) createJenkinsJob(jobName string, jobDescription string, githubURL string, jenkinsUser string, token string) error {
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

	req.SetBasicAuth(jenkinsUser, token)
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
