package github_client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v65/github"
	"github.com/kalgurn/github-rate-limits-prometheus-exporter/internal/utils"
	"golang.org/x/oauth2"
)

func GetRemainingLimits(c *github.Client) RateLimits {
	ctx := context.Background()

	limits, _, err := c.RateLimit.Get(ctx)
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

func (c *TokenConfig) InitClient() *github.Client {
	return initTokenClient(c, http.DefaultClient)
}

func (c *AppConfig) InitClient() *github.Client {
	return initAppClient(c, http.DefaultClient)
}

func InitConfig() GithubClient {
	// determine type (app or pat)
	var auth GithubClient
	authType := utils.GetOSVar("GITHUB_AUTH_TYPE")
	if authType == "PAT" {
		auth = &TokenConfig{
			Token: utils.GetOSVar("GITHUB_TOKEN"),
		}

	} else if authType == "APP" {
		appID, _ := strconv.ParseInt(utils.GetOSVar("GITHUB_APP_ID"), 10, 64)

		var installationID int64
		envInstallationID := utils.GetOSVar("GITHUB_INSTALLATION_ID")
		if envInstallationID != "" {
			installationID, _ = strconv.ParseInt(envInstallationID, 10, 64)
		}

		auth = &AppConfig{
			AppID:          appID,
			InstallationID: installationID,
			OrgName:        utils.GetOSVar("GITHUB_ORG_NAME"),
			RepoName:       utils.GetOSVar("GITHUB_REPO_NAME"),
			PrivateKeyPath: utils.GetOSVar("GITHUB_PRIVATE_KEY_PATH"),
		}
	} else {
		err := fmt.Errorf("invalid auth type")
		utils.RespError(err)
		return nil
	}

	return auth

}

// Helper function to allow testing client initialization with custom http clients
func initTokenClient(c *TokenConfig, httpClient *http.Client) *github.Client {
	if httpClient == http.DefaultClient {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.Token},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}
	return github.NewClient(httpClient)
}

// Helper function to allow testing client initialization with custom http clients
func initAppClient(c *AppConfig, httpClient *http.Client) *github.Client {
	if c.InstallationID == 0 && c.OrgName != "" {
		// Retrieve the installation ID if not provided
		auth := &TokenConfig{
			Token: generateJWT(c.AppID, c.PrivateKeyPath),
		}
		client := initTokenClient(auth, httpClient)

		var err error
		var installation *github.Installation
		ctx := context.Background()
		if c.RepoName != "" {
			installation, _, err = client.Apps.FindRepositoryInstallation(ctx, c.OrgName, c.RepoName)
		} else {
			installation, _, err = client.Apps.FindOrganizationInstallation(ctx, c.OrgName)
		}
		utils.RespError(err)

		c.InstallationID = installation.GetID()
	}

	if httpClient == http.DefaultClient {
		tr := http.DefaultTransport
		itr, err := ghinstallation.NewKeyFromFile(tr, c.AppID, c.InstallationID, c.PrivateKeyPath)
		utils.RespError(err)
		httpClient = &http.Client{Transport: itr}
	} else {
		// Wrap the existing transport
		tr := httpClient.Transport
		if tr == nil {
			tr = http.DefaultTransport
		}
		itr, err := ghinstallation.NewKeyFromFile(tr, c.AppID, c.InstallationID, c.PrivateKeyPath)
		utils.RespError(err)
		httpClient.Transport = itr
	}

	return github.NewClient(httpClient)
}

// Helper function to generate JWT for GitHub App
func generateJWT(appID int64, privateKeyPath string) string {
	privateKey, err := os.ReadFile(privateKeyPath)
	utils.RespError(err)

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	utils.RespError(err)

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("%d", appID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(parsedKey)
	utils.RespError(err)

	return signedToken
}
