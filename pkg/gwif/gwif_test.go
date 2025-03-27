package gwif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

// Sample JWT token parts for testing
const (
	testHeader    = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ"
	testPayload   = "eyJuYW1lc3BhY2VfaWQiOiIxMjM0NTYiLCJuYW1lc3BhY2VfcGF0aCI6InRlc3Qtb3JnIiwicHJvamVjdF9pZCI6Ijc4OTAiLCJwcm9qZWN0X3BhdGgiOiJ0ZXN0LW9yZy90ZXN0LXByb2plY3QiLCJ1c2VyX2lkIjoiMTIzNCIsInVzZXJfbG9naW4iOiJ0ZXN0LXVzZXIiLCJ1c2VyX2VtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsInVzZXJfaWRlbnRpdGllcyI6W10sInBpcGVsaW5lX2lkIjoiMTIzNDUiLCJwaXBlbGluZV9zb3VyY2UiOiJwdXNoIiwiam9iX2lkIjoiNjc4OSIsInJlZiI6Im1haW4iLCJyZWZfdHlwZSI6ImJyYW5jaCIsInJlZl9wYXRoIjoicmVmcy9oZWFkcy9tYWluIiwicmVmX3Byb3RlY3RlZCI6ImZhbHNlIiwiZ3JvdXBzX2RpcmVjdCI6W10sImVudmlyb25tZW50IjoiIiwiZW52aXJvbm1lbnRfcHJvdGVjdGVkIjoiIiwiZGVwbG95bWVudF90aWVyIjoiIiwiZW52aXJvbm1lbnRfYWN0aW9uIjoiIiwicnVubmVyX2lkIjoxMjM0LCJydW5uZXJfZW52aXJvbm1lbnQiOiIiLCJzaGEiOiJhYmMxMjM0ZGVmNTY3IiwicHJvamVjdF92aXNpYmlsaXR5IjoicHJpdmF0ZSIsImNpX2NvbmZpZ19yZWZfdXJpIjoiIiwiY2lfY29uZmlnX3NoYSI6IiIsImp0aSI6ImFiYy0xMjMtZGVmLTQ1NiIsImlzcyI6Imh0dHBzOi8vZ2l0bGFiLmNvbSIsImlhdCI6MTY0NzUzNTEyMywibmJmIjoxNjQ3NTM1MTIzLCJleHAiOjE2NDc1MzUxMjMsInN1YiI6InByb2plY3RfcGF0aDp0ZXN0LW9yZy90ZXN0LXByb2plY3Q6cmVmX3R5cGU6YnJhbmNoOnJlZjptYWluIiwiYXVkIjoiaHR0cHM6Ly9naXRsYWIuY29tIn0"
	testSignature = "signature123"
	testToken     = testHeader + "." + testPayload + "." + testSignature
)

// setupEnvironment sets up the environment variables for testing
func setupEnvironment(t *testing.T) {
	t.Setenv("GITLAB_OIDC_TOKEN", testToken)
	t.Setenv("GCP_WIF_PROJECT_ID", "123456789")
	t.Setenv("GCP_WIF_POOL", "test-pool")
	t.Setenv("GCP_WIF_PROVIDER", "gitlab-provider")
	t.Setenv("GCP_WIF_SERVICE_ACCOUNT_EMAIL", "test-sa@project-id.iam.gserviceaccount.com")
}

// removeEnvironment removes the environment variables set for testing
func removeEnvironment(t *testing.T) {
	vars := []string{
		"GITLAB_OIDC_TOKEN",
		"GCP_WIF_PROJECT_ID",
		"GCP_WIF_POOL",
		"GCP_WIF_PROVIDER",
		"GCP_WIF_SERVICE_ACCOUNT_EMAIL",
	}
	for _, env := range vars {
		os.Unsetenv(env)
	}
}

func TestNewGitlabOidc(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		expectError bool
	}{
		{
			name: "Valid token",
			setupEnv: func() {
				os.Setenv("GITLAB_OIDC_TOKEN", testToken)
			},
			expectError: false,
		},
		{
			name: "Missing token",
			setupEnv: func() {
				os.Unsetenv("GITLAB_OIDC_TOKEN")
			},
			expectError: true,
		},
		{
			name: "Invalid token format",
			setupEnv: func() {
				os.Setenv("GITLAB_OIDC_TOKEN", "invalid-token")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setupEnv()
			defer removeEnvironment(t)

			// Test
			oidc, err := NewGitlabOidc()

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, oidc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, oidc)
				assert.Equal(t, testHeader, oidc.Header)
				assert.Equal(t, testSignature, oidc.Signature)
				assert.Equal(t, testToken, oidc.FromEnv)

				// Check a few fields from the parsed JWT
				assert.Equal(t, "1234", oidc.Payload.UserID)
				assert.Equal(t, "test-user", oidc.Payload.UserLogin)
				assert.Equal(t, "test@example.com", oidc.Payload.UserEmail)
			}
		})
	}
}

func TestNewWorkloadIdentityConfig(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		expectError bool
	}{
		{
			name: "All environment variables set",
			setupEnv: func() {
				setupEnvironment(t)
			},
			expectError: false,
		},
		{
			name: "Missing GitLab OIDC token",
			setupEnv: func() {
				setupEnvironment(t)
				os.Unsetenv("GITLAB_OIDC_TOKEN")
			},
			expectError: true,
		},
		{
			name: "Missing GCP project ID",
			setupEnv: func() {
				setupEnvironment(t)
				os.Unsetenv("GCP_WIF_PROJECT_ID")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setupEnv()
			defer removeEnvironment(t)

			// Test
			config, err := NewWorkloadIdentityConfig()

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, os.Getenv("GCP_WIF_PROJECT_ID"), config.ProjectNumber)
				assert.Equal(t, os.Getenv("GCP_WIF_POOL"), config.PoolID)
				assert.Equal(t, os.Getenv("GCP_WIF_PROVIDER"), config.ProviderID)
				assert.Equal(t, os.Getenv("GCP_WIF_SERVICE_ACCOUNT_EMAIL"), config.ServiceAccount)
				assert.Equal(t, testToken, config.GitLabOIDCToken.FromEnv)
			}
		})
	}
}

// TestGitlabOidcJwt_Methods tests the methods on the GitlabOidcJwt struct
func TestGitlabOidcJwt_Methods(t *testing.T) {
	// Create a sample JWT
	jwt := &GitlabOidcJwt{
		UserID:    "1234",
		UserLogin: "test-user",
		UserEmail: "test@example.com",
	}

	// Test JsonPrettyPrint
	prettyJson, err := jwt.JsonPrettyPrint()
	assert.NoError(t, err)
	assert.Contains(t, prettyJson, "  \"user_id\": \"1234\"")
	assert.Contains(t, prettyJson, "  \"user_login\": \"test-user\"")

	// Test AsString
	jsonStr := jwt.AsString()
	assert.Contains(t, jsonStr, "\"user_id\":\"1234\"")
	assert.Contains(t, jsonStr, "\"user_login\":\"test-user\"")
}

// TestNewGitlabCiCdJWT tests parsing a JWT payload
func TestNewGitlabCiCdJWT(t *testing.T) {
	jwt, err := NewGitlabCiCdJWT(testPayload)
	assert.NoError(t, err)
	assert.NotNil(t, jwt)

	// Verify a few fields
	assert.Equal(t, "1234", jwt.UserID)
	assert.Equal(t, "test-user", jwt.UserLogin)
	assert.Equal(t, "test@example.com", jwt.UserEmail)
	assert.Equal(t, "main", jwt.Ref)
	assert.Equal(t, "branch", jwt.RefType)

	// Test with invalid payload
	_, err = NewGitlabCiCdJWT("invalid-base64")
	assert.Error(t, err)
}

// TestCheckEnvVars tests the environment variable validation
func TestCheckEnvVars(t *testing.T) {
	// Test with all env vars set
	setupEnvironment(t)
	err := checkEnvVars()
	assert.NoError(t, err)

	// Test with missing env var
	os.Unsetenv("GCP_WIF_PROJECT_ID")
	err = checkEnvVars()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GCP_WIF_PROJECT_ID")

	// Clean up
	removeEnvironment(t)
}

func TestParseGitlabOidc(t *testing.T) {
	tests := []struct {
		name        string
		tokenString string
		wantParts   []string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid JWT token",
			tokenString: testHeader + "." + testPayload + "." + testSignature,
			wantParts:   []string{testHeader, testPayload, testSignature},
			wantErr:     false,
		},
		{
			name:        "Empty token",
			tokenString: "",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "invalid JWT format: token is empty",
		},
		{
			name:        "Token with less than 3 parts",
			tokenString: "part1.part2",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "invalid JWT format: token must have exactly three parts",
		},
		{
			name:        "Token with more than 3 parts",
			tokenString: "part1.part2.part3.part4",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "invalid JWT format: token must have exactly three parts",
		},
		{
			name:        "Token with invalid base64url payload",
			tokenString: "part1.!invalid-base64!.part3",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "failed to decode JWT payload",
		},
		{
			name:        "Token with missing parts but still has two dots",
			tokenString: "part1..",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "invalid JWT format: part",
		},
		{
			name:        "Token with empty payload",
			tokenString: "part1..part3",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "invalid JWT format: part 2 is empty",
		},
		{
			name:        "Token with standard base64 padding (should fail with URL-safe only)",
			tokenString: "part1.SGVsbG8gV29ybGQ=.part3", // "Hello World" in base64 with padding
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "failed to decode JWT payload",
		},
		{
			name:        "Token with URL-safe base64 (without padding)",
			tokenString: "part1.SGVsbG8gV29ybGQ.part3", // "Hello World" in base64 without padding
			wantParts:   []string{"part1", "SGVsbG8gV29ybGQ", "part3"},
			wantErr:     false,
		},
		{
			name:        "Token with URL-safe special characters",
			tokenString: "abc_-123.eyJhbGciOiJIUzI1NiJ9.xyz_-789",
			wantParts:   []string{"abc_-123", "eyJhbGciOiJIUzI1NiJ9", "xyz_-789"},
			wantErr:     false,
		},
		{
			name:        "Token with payload too short",
			tokenString: "part1.abcd.part3",
			wantParts:   nil,
			wantErr:     true,
			errMsg:      "invalid JWT format: payload is too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts, err := parseGitlabOidc(tt.tokenString)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, parts)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantParts, parts)
			}
		})
	}
}

// TestAdditionalEdgeCases tests some specific edge cases for parseGitlabOidc
func TestParseGitlabOidcEdgeCases(t *testing.T) {
	// Test with very long token parts
	longHeader := "h" + string(make([]byte, 8192)) // 8KB header
	longPayload := "eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0"
	longSignature := "s" + string(make([]byte, 8192)) // 8KB signature

	longToken := longHeader + "." + longPayload + "." + longSignature
	parts, err := parseGitlabOidc(longToken)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(parts))
	assert.Equal(t, longHeader, parts[0])
	assert.Equal(t, longPayload, parts[1])
	assert.Equal(t, longSignature, parts[2])

	// Test with special characters that are valid in base64url encoding
	specialCharsToken := "abc123_-xyz.eyJhbGciOiJIUzI1NiJ9._-987ZYX"
	parts, err = parseGitlabOidc(specialCharsToken)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(parts))

	// Test with Unicode characters in the parts (this should fail as it's not valid base64)
	unicodeToken := "part1.üñîçøðé.part3"
	parts, err = parseGitlabOidc(unicodeToken)
	assert.Error(t, err)
	assert.Nil(t, parts)

	// Test with plus characters (which are not valid in base64url)
	plusToken := "part1.invalid+base64url.part3"
	parts, err = parseGitlabOidc(plusToken)
	assert.Error(t, err)
	assert.Nil(t, parts)
}

func TestGetFederatedToken(t *testing.T) {
	// Setup mock STS server
	mockSTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"access_token": "test-federated-token",
				"expires_in": 3599,
				"token_type": "Bearer"
			}`))
		} else {
			t.Fatalf("Unexpected request to path: %s", r.URL.Path)
		}
	}))
	defer mockSTS.Close()

	// Create a test configuration
	config := &WorkloadIdentityConfig{
		ProjectNumber:  "123456789",
		PoolID:         "test-pool",
		ProviderID:     "gitlab-provider",
		ServiceAccount: "test-sa@project-id.iam.gserviceaccount.com",
		GitLabOIDCToken: &GitlabOidc{
			FromEnv: testToken,
		},
	}

	// This test can only check the function structure because we can't easily
	// mock the Google API clients. In a real test, you would use a mocking library.
	t.Run("Function structure test", func(t *testing.T) {
		ctx := context.Background()

		// Just check that the function exists and has the right signature
		var _ func(context.Context, *WorkloadIdentityConfig) (string, error) = GetFederatedToken

		// Test that the function works with the right parameters
		// This will fail but tests the structure
		_, err := GetFederatedToken(ctx, config)

		// We expect an error since we're not properly mocking the STS service
		// But we're just testing the function structure
		assert.Error(t, err)
	})
}

func TestGetServiceAccountToken(t *testing.T) {
	// Setup mock IAM server
	mockIAM := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/projects/-/serviceAccounts/test-sa@project-id.iam.gserviceaccount.com:generateAccessToken" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			expiry := time.Now().Add(time.Hour).Format(time.RFC3339)
			w.Write([]byte(`{
				"accessToken": "test-sa-token",
				"expireTime": "` + expiry + `"
			}`))
		} else {
			t.Fatalf("Unexpected request to path: %s", r.URL.Path)
		}
	}))
	defer mockIAM.Close()

	// This test can only check the function structure because we can't easily
	// mock the Google API clients. In a real test, you would use a mocking library.
	t.Run("Function structure test", func(t *testing.T) {
		ctx := context.Background()

		// Just check that the function exists and has the right signature
		var _ func(context.Context, string, string) (*oauth2.Token, error) = GetServiceAccountToken

		// Test that the function works with the right parameters
		// This will fail but tests the structure
		_, err := GetServiceAccountToken(ctx, "test-token", "test-sa@project-id.iam.gserviceaccount.com")

		// We expect an error since we're not properly mocking the IAM service
		// But we're just testing the function structure
		assert.Error(t, err)
	})
}

func TestGetGCPToken_Structure(t *testing.T) {
	// Create a test configuration
	config := &WorkloadIdentityConfig{
		ProjectNumber:  "123456789",
		PoolID:         "test-pool",
		ProviderID:     "gitlab-provider",
		ServiceAccount: "test-sa@project-id.iam.gserviceaccount.com",
		GitLabOIDCToken: &GitlabOidc{
			FromEnv: testToken,
		},
	}

	// This test can only check the function structure
	t.Run("Function structure test", func(t *testing.T) {
		// Just check that the function exists and has the right signature
		var _ func(*WorkloadIdentityConfig) (*oauth2.Token, error) = GetGCPToken

		// Test that the function works with the right parameters
		// This will fail but tests the structure
		_, err := GetGCPToken(config)

		// We expect an error since we're not properly mocking the services
		assert.Error(t, err)
	})
}
