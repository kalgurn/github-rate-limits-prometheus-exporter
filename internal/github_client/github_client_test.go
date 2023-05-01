package github_client

import (
	"math"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetRemainingLimits(t *testing.T) {
	var (
		limit        = 100
		remaining    = 63
		used         = 37
		seconds_left = 1500
	)
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetRateLimit,
			struct {
				Resources *github.RateLimits
			}{
				Resources: &github.RateLimits{
					Core: &github.Rate{
						Limit:     limit,
						Remaining: remaining,
						Reset:     github.Timestamp{Time: time.Now().Add(time.Second * time.Duration(seconds_left))},
					},
					Search: &github.Rate{},
				},
			},
		),
	)
	c := github.NewClient(mockedHTTPClient)
	limits := GetRemainingLimits(c)

	assert.Equal(t, limit, limits.Limit, "The limits should be equal")
	assert.Equal(t, remaining, limits.Remaining, "The remaining limits should be equal")
	assert.Equal(t, used, limits.Used, "The used value should be equal")
	assert.Equal(t, seconds_left, int(math.Ceil(limits.SecondsLeft)), "The seconds left value should be equal")

	assert.NotEqual(t, 99, limits.Limit, "The limit should not be equal")
	assert.NotEqual(t, 99, limits.Remaining, "The remaining limits should not be equal")
	assert.NotEqual(t, 18, limits.Used, "The used value should not be equal")
	assert.NotEqual(t, 18, limits.Used, "The seconds left value should not be equal")
}

func TestInitConfigApp(t *testing.T) {
	os.Setenv("GITHUB_AUTH_TYPE", "APP")
	os.Setenv("GITHUB_APP_ID", "1")
	os.Setenv("GITHUB_INSTALLATION_ID", "1")
	os.Setenv("GITHUB_PRIVATE_KEY_PATH", "/home")

	testAuth := AppConfig{
		AppID:          1,
		InstallationID: 1,
		PrivateKeyPath: "/home",
	}

	appInitConfig := InitConfig()

	assert.Equal(t, appInitConfig, testAuth, "should be equal")

}

func TestInitConfigPAT(t *testing.T) {
	os.Setenv("GITHUB_AUTH_TYPE", "PAT")
	os.Setenv("GITHUB_TOKEN", "token_ahsd")

	testAuth := TokenConfig{
		Token: "token_ahsd",
	}

	patInitConfig := InitConfig()

	assert.Equal(t, patInitConfig, testAuth, "should be equal")

}

func TestInitConfigFailure(t *testing.T) {
	os.Setenv("GITHUB_AUTH_TYPE", "test")

	patInitConfig := InitConfig()

	assert.Equal(t, nil, patInitConfig)

}
