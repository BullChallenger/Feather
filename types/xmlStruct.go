package types

import "encoding/xml"

type githubProjectProperty struct {
	Plugin      string `xml:"plugin,attr"`
	ProjectUrl  string `xml:"projectUrl"`
	DisplayName string `xml:"displayName"`
}

type durabilityHintJobProperty struct {
	Hint string `xml:"hint"`
}

type gitHubPushTrigger struct {
	Plugin string `xml:"plugin,attr"`
	Spec   string `xml:"spec"`
}

type triggers struct {
	GitHubPushTrigger gitHubPushTrigger `xml:"com.cloudbees.jenkins.GitHubPushTrigger"`
}

type pipelineTriggersJobProperty struct {
	Triggers triggers `xml:"triggers"`
}

type disableConcurrentBuildsJobProperty struct {
	AbortPrevious bool `xml:"abortPrevious"`
}

type properties struct {
	DisableConcurrentBuildsJobProperty disableConcurrentBuildsJobProperty `xml:"org.jenkinsci.plugins.workflow.job.properties.DisableConcurrentBuildsJobProperty"`
	GithubProjectProperty              githubProjectProperty              `xml:"com.coravy.hudson.plugins.github.GithubProjectProperty"`
	DurabilityHintJobProperty          durabilityHintJobProperty          `xml:"org.jenkinsci.plugins.workflow.job.properties.DurabilityHintJobProperty"`
	PipelineTriggersJobProperty        pipelineTriggersJobProperty        `xml:"org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty"`
}

type branchSpec struct {
	Name string `xml:"name"`
}

type branches struct {
	BranchSpec branchSpec `xml:"hudson.plugins.git.BranchSpec"`
}

type submoduleCfg struct {
	Class string `xml:"class,attr"`
}

type userRemoteConfig struct {
	URL           string `xml:"url"`
	CredentialsId string `xml:"credentialsId"`
}

type userRemoteConfigs struct {
	UserRemoteConfig userRemoteConfig `xml:"hudson.plugins.git.UserRemoteConfig"`
}

type scm struct {
	Class                             string            `xml:"class,attr"`
	Plugin                            string            `xml:"plugin,attr"`
	ConfigVersion                     int               `xml:"configVersion"`
	UserRemoteConfigs                 userRemoteConfigs `xml:"userRemoteConfigs"`
	Branches                          branches          `xml:"branches"`
	DoGenerateSubmoduleConfigurations bool              `xml:"doGenerateSubmoduleConfigurations"`
	SubmoduleCfg                      submoduleCfg      `xml:"submoduleCfg"`
}

type definition struct {
	Class       string `xml:"class,attr"`
	Plugin      string `xml:"plugin,attr"`
	Scm         scm    `xml:"scm"`
	ScriptPath  string `xml:"scriptPath"`
	Lightweight bool   `xml:"lightweight"`
}

type FlowDefinition struct {
	XMLName          xml.Name   `xml:"flow-definition"`
	Plugin           string     `xml:"plugin,attr"`
	Description      string     `xml:"description"`
	KeepDependencies bool       `xml:"keepDependencies"`
	Properties       properties `xml:"properties"`
	Definition       definition `xml:"definition"`
	Disabled         bool       `xml:"disabled"`
}
