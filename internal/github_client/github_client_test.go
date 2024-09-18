package github_client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v65/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func generateTestPrivateKey(t *testing.T) (string, *rsa.PrivateKey) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA private key: %v", err)
	}

	// Convert private key to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Write private key to a temp file
	tempKeyFile, err := os.CreateTemp("", "testkey")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tempKeyFile.Close()

	if _, err := tempKeyFile.Write(privateKeyPEM); err != nil {
		t.Fatalf("Failed to write to temp key file: %v", err)
	}

	return tempKeyFile.Name(), privateKey
}

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

	testAuth := &AppConfig{
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

	testAuth := &TokenConfig{
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

func TestInitConfigAppWithoutInstallationID(t *testing.T) {
	os.Setenv("GITHUB_AUTH_TYPE", "APP")
	os.Setenv("GITHUB_APP_ID", "1")
	os.Setenv("GITHUB_ORG_NAME", "org")
	os.Setenv("GITHUB_PRIVATE_KEY_PATH", "/home")

	testAuth := &AppConfig{
		AppID:          1,
		OrgName:        "org",
		PrivateKeyPath: "/home",
	}

	appInitConfig := InitConfig()

	assert.Equal(t, appInitConfig, testAuth, "should be equal")
}

func TestAppConfig_InitClient(t *testing.T) {
	testCases := []struct {
		name              string
		orgName           string
		repoName          string
		providedInstallID int64 // InstallationID provided directly in AppConfig
		expectedInstallID int64 // Expected InstallationID after InitClient
		expectedPattern   string
		method            string
	}{
		{
			name:              "WithInstallationID",
			orgName:           "",
			repoName:          "",
			providedInstallID: 654321,
			expectedInstallID: 654321,
			expectedPattern:   "", // No API call expected
			method:            "",
		},
		{
			name:              "WithOrgName",
			orgName:           "testorg",
			repoName:          "",
			providedInstallID: 0, // To be retrieved via API
			expectedInstallID: 654321,
			expectedPattern:   "/orgs/{org}/installation",
			method:            "GET",
		},
		{
			name:              "WithOrgAndRepoName",
			orgName:           "testorg",
			repoName:          "testrepo",
			providedInstallID: 0, // To be retrieved via API
			expectedInstallID: 654321,
			expectedPattern:   "/repos/{owner}/{repo}/installation",
			method:            "GET",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privateKeyPath, _ := generateTestPrivateKey(t)
			defer os.Remove(privateKeyPath)

			appID := int64(123456)
			var httpClient *http.Client

			if tc.expectedPattern != "" {
				// Create a mock HTTP client to simulate API call
				mockClient := mock.NewMockedHTTPClient(
					mock.WithRequestMatchHandler(
						mock.EndpointPattern{Pattern: tc.expectedPattern, Method: tc.method},
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							// Return mock installation data
							installation := &github.Installation{
								ID: github.Int64(tc.expectedInstallID),
							}
							data, _ := json.Marshal(installation)
							w.WriteHeader(http.StatusOK)
							w.Write(data)
						}),
					),
				)
				httpClient = mockClient
			} else {
				httpClient = nil // No HTTP client needed; no API call expected
			}

			// Initialize the AppConfig
			c := &AppConfig{
				AppID:          appID,
				InstallationID: tc.providedInstallID,
				OrgName:        tc.orgName,
				RepoName:       tc.repoName,
				PrivateKeyPath: privateKeyPath,
			}

			client := initAppClient(c, httpClient)
			assert.NotNil(t, client, "Expected client not to be nil")

			assert.Equal(t, tc.expectedInstallID, c.InstallationID, "Expected InstallationID to be set correctly")
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	privateKeyPath, privateKey := generateTestPrivateKey(t)
	defer os.Remove(privateKeyPath)

	appID := int64(123456)
	token := generateJWT(appID, privateKeyPath)

	assert.NotEmpty(t, token, "expected token not to be empty")

	// Verify the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	assert.True(t, parsedToken.Valid, "the token should be valid")

	// Check claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		issuer := claims["iss"]
		assert.Equal(t, fmt.Sprintf("%d", appID), issuer, "expected issuer to be equal app id")

		exp := int64(claims["exp"].(float64))
		now := time.Now().Unix()
		assert.LessOrEqual(t, now, exp, "expected token to not be expired")
	} else {
		t.Error("Failed to parse claims")
	}
}
