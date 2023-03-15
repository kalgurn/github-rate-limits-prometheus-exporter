package github_client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
	"github.com/kalgurn/github-rate-limits-prometheus-exporter/internal/utils"
	"golang.org/x/oauth2"
)

func GetRemainingLimits(c *github.Client) RateLimits {
	ctx := context.Background()

	limits, _, err := c.RateLimits(ctx)
	if err != nil {
		utils.RespError(err)
	}

	return RateLimits{
		Limit:       limits.Core.Limit,
		Remaining:   limits.Core.Remaining,
		Used:        limits.Core.Limit - limits.Core.Remaining,
		SecondsLeft: time.Until(limits.Core.Reset.Time).Seconds(),
	}
}

func (c TokenConfig) InitClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)

}

func (c AppConfig) InitClient() *github.Client {
	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport

	// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewKeyFromFile(tr, c.AppID, c.InstallationID, c.PrivateKeyPath)
	utils.RespError(err)

	// Use installation transport with github.com/google/go-github
	return github.NewClient(&http.Client{Transport: itr})
}

func InitConfig() GithubClient {
	// determine type (app or pat)
	var auth GithubClient
	authType := utils.GetOSVar("GITHUB_AUTH_TYPE")
	if authType == "PAT" {
		auth = TokenConfig{
			Token: utils.GetOSVar("GITHUB_TOKEN"),
		}

	} else if authType == "APP" {
		appID, _ := strconv.ParseInt(utils.GetOSVar("GITHUB_APP_ID"), 10, 64)
		installationID, _ := strconv.ParseInt(utils.GetOSVar("GITHUB_INSTALLATION_ID"), 10, 64)

		auth = AppConfig{
			AppID:          appID,
			InstallationID: installationID,
			PrivateKeyPath: utils.GetOSVar("GITHUB_PRIVATE_KEY_PATH"),
		}
	} else {
		err := fmt.Errorf("invalid auth type")
		utils.RespError(err)
		return nil
	}

	return auth

}
