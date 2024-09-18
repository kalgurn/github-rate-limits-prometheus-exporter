package github_client

import (
	"github.com/google/go-github/v65/github"
)

type AppConfig struct {
	AppID          int64
	InstallationID int64
	OrgName        string
	RepoName       string
	PrivateKeyPath string
}

type TokenConfig struct {
	Token string
}

type RateLimits struct {
	Limit       int
	Remaining   int
	Used        int
	SecondsLeft float64
}

type GithubClient interface {
	InitClient() *github.Client
}
