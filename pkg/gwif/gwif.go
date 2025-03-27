package gwif

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sts/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type GitlabOidc struct {
	Header    string
	Payload   GitlabOidcJwt
	Signature string
	FromEnv   string
}

type GitlabOidcJwt struct {
	NamespaceID          string         `json:"namespace_id"`
	NamespacePath        string         `json:"namespace_path"`
	ProjectID            string         `json:"project_id"`
	ProjectPath          string         `json:"project_path"`
	UserID               string         `json:"user_id"`
	UserLogin            string         `json:"user_login"`
	UserEmail            string         `json:"user_email"`
	UserIdentities       []UserIdentity `json:"user_identities"`
	PipelineID           string         `json:"pipeline_id"`
	PipelineSource       string         `json:"pipeline_source"`
	JobID                string         `json:"job_id"`
	Ref                  string         `json:"ref"`
	RefType              string         `json:"ref_type"`
	RefPath              string         `json:"ref_path"`
	RefProtected         string         `json:"ref_protected"`
	GroupsDirect         []string       `json:"groups_direct"`
	Environment          string         `json:"environment"`
	EnvironmentProtected string         `json:"environment_protected"`
	DeploymentTier       string         `json:"deployment_tier"`
	EnvironmentAction    string         `json:"environment_action"`
	RunnerID             int            `json:"runner_id"`
	RunnerEnvironment    string         `json:"runner_environment"`
	Sha                  string         `json:"sha"`
	ProjectVisibility    string         `json:"project_visibility"`
	CiConfigRefUri       string         `json:"ci_config_ref_uri"`
	CiConfigSha          string         `json:"ci_config_sha"`
	Jti                  string         `json:"jti"`
	Iss                  string         `json:"iss"`
	Iat                  int64          `json:"iat"`
	Nbf                  int64          `json:"nbf"`
	Exp                  int64          `json:"exp"`
	Sub                  string         `json:"sub"`
	Aud                  string         `json:"aud"`
}

type UserIdentity struct {
	Provider  string `json:"provider"`
	ExternUID string `json:"extern_uid"`
}

type WorkloadIdentityConfig struct {
	ProjectNumber   string
	PoolID          string
	ProviderID      string
	ServiceAccount  string
	GitLabOIDCToken GitlabOidc
}

func (gwif *WorkloadIdentityConfig) SafeToMap() map[string]string {
	safeConfig := map[string]string{
		"project_number":  gwif.ProjectNumber,
		"pool_id":         gwif.PoolID,
		"provider_id":     gwif.ProviderID,
		"service_account": gwif.ServiceAccount,
	}
	return safeConfig
}

func (jwt *GitlabOidcJwt) JsonPrettyPrint() (string, error) {
	jsonBytes, err := json.MarshalIndent(jwt, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (wic *GitlabOidcJwt) AsString() string {
	jsonBytes, err := json.Marshal(wic)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func NewGitlabOidc() (*GitlabOidc, error) {
	tokenString := os.Getenv("GITLAB_OIDC_TOKEN")
	if tokenString == "" {
		return nil, errors.New("GITLAB_OIDC_TOKEN environment variable not set")
	}
	parts, err := parseGitlabOidc(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GitLab OIDC token: %w", err)
	}
	jwt, err := NewGitlabCiCdJWT(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab OIDC token: %w", err)
	}
	return &GitlabOidc{
		Header:    parts[0],
		Payload:   *jwt,
		Signature: parts[2],
		FromEnv:   tokenString,
	}, nil
}

func parseGitlabOidc(tokenString string) ([]string, error) {
	if tokenString == "" {
		return nil, errors.New("invalid JWT format: token is empty")
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format: token must have exactly three parts separated by dots")
	}

	for i, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid JWT format: part %d is empty", i+1)
		}
	}

	_, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Arbitrary minimum length to ensure it's a meaningful payload
	if len(parts[1]) < 8 {
		return nil, errors.New("invalid JWT format: payload is too short")
	}

	return parts, nil
}

func NewWorkloadIdentityConfig() (*WorkloadIdentityConfig, error) {
	if err := checkEnvVars(); err != nil {
		return nil, err
	}

	oidc, err := NewGitlabOidc()
	if err != nil {
		return nil, fmt.Errorf("error parsing GitLab OIDC token: %v", err)
	}

	return &WorkloadIdentityConfig{
		ProjectNumber:   os.Getenv("GCP_WIF_PROJECT_ID"),
		PoolID:          os.Getenv("GCP_WIF_POOL"),
		ProviderID:      os.Getenv("GCP_WIF_PROVIDER"),
		ServiceAccount:  os.Getenv("GCP_WIF_SERVICE_ACCOUNT_EMAIL"),
		GitLabOIDCToken: *oidc,
	}, nil
}

func NewGitlabCiCdJWT(tokenString string) (*GitlabOidcJwt, error) {
	payload, err := base64.RawURLEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var jwt GitlabOidcJwt
	if err := json.Unmarshal(payload, &jwt); err != nil {
		return nil, fmt.Errorf("failed to parse JWT payload: %w", err)
	}

	if jwt.UserIdentities == nil {
		jwt.UserIdentities = []UserIdentity{}
	}
	if jwt.GroupsDirect == nil {
		jwt.GroupsDirect = []string{}
	}

	return &jwt, nil
}

func checkEnvVars() error {
	var envVars = []string{
		"GITLAB_OIDC_TOKEN",
		"GCP_WIF_PROJECT_ID",
		"GCP_WIF_POOL",
		"GCP_WIF_PROVIDER",
		"GCP_WIF_SERVICE_ACCOUNT_EMAIL",
	}
	for _, env := range envVars {
		if os.Getenv(env) == "" {
			return fmt.Errorf("missing environment variable: %s", env)
		}
	}
	return nil
}

func (config *WorkloadIdentityConfig) getAudience() string {
	return fmt.Sprintf("//iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/providers/%s",
		config.ProjectNumber, config.PoolID, config.ProviderID)
}

// GetFederatedToken exchanges a GitLab OIDC token for a GCP federated token
func GetFederatedToken(ctx context.Context, config *WorkloadIdentityConfig) (string, error) {
	// 1. Initialize the STS (Security Token Service) client
	stsService, err := sts.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", fmt.Errorf("failed to create STS service: %w", err)
	}

	// 2. Create the resource name for the provider pool
	audience := config.getAudience()

	// 3. Exchange the GitLab OIDC token for a GCP federated token
	exchangeReq := &sts.GoogleIdentityStsV1ExchangeTokenRequest{
		Audience:           audience,
		GrantType:          "urn:ietf:params:oauth:grant-type:token-exchange",
		RequestedTokenType: "urn:ietf:params:oauth:token-type:access_token",
		SubjectTokenType:   "urn:ietf:params:oauth:token-type:jwt",
		SubjectToken:       config.GitLabOIDCToken.FromEnv,
		Scope:              "https://www.googleapis.com/auth/cloud-platform",
	}

	exchangeToken, err := stsService.V1.Token(exchangeReq).Do()
	if err != nil {
		return "", fmt.Errorf("token exchange error: %w", err)
	}

	return exchangeToken.AccessToken, nil
}

// GetServiceAccountToken impersonates a service account using a federated token
func GetServiceAccountToken(ctx context.Context, federatedToken, serviceAccountEmail string) (*oauth2.Token, error) {
	// 1. Create IAM Credentials client with federated token
	iamCredsOpts := []option.ClientOption{
		option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: federatedToken,
		})),
	}

	iamCredsClient, err := credentials.NewIamCredentialsClient(ctx, iamCredsOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM credentials client: %w", err)
	}
	defer iamCredsClient.Close()

	// 2. Generate an access token for the service account
	serviceAccountName := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail)

	// 3. Using the Cloud Client Libraries for service account impersonation
	accessTokenReq := &credentialspb.GenerateAccessTokenRequest{
		Name:  serviceAccountName,
		Scope: []string{"https://www.googleapis.com/auth/cloud-platform"},
		// Default lifetime is 1 hour if not specified
		Lifetime: &durationpb.Duration{Seconds: 3600},
	}

	accessTokenResp, err := iamCredsClient.GenerateAccessToken(ctx, accessTokenReq)
	if err != nil {
		return nil, fmt.Errorf("service account impersonation error: %w", err)
	}

	// 4. Parse the expiry time
	expiry := accessTokenResp.ExpireTime.AsTime()
	if expiry.IsZero() {
		// Fallback: use the current time plus one hour if parsing fails
		expiry = time.Now().Add(time.Hour)
	}

	// 5. Create an OAuth2 token to return
	token := &oauth2.Token{
		AccessToken: accessTokenResp.AccessToken,
		TokenType:   "Bearer",
		Expiry:      expiry,
	}

	return token, nil
}

// GetGCPToken gets a GCP access token using a GitLab OIDC token
func GetGCPToken(config *WorkloadIdentityConfig) (*oauth2.Token, error) {
	if config == nil {
		return nil, errors.New("WorkloadIdentityConfig is nil")
	}

	ctx := context.Background()

	// 1. Get federated token
	federatedToken, err := GetFederatedToken(ctx, config)
	if err != nil {
		return nil, err
	}

	// 2. Use the federated token to get a service account token
	return GetServiceAccountToken(ctx, federatedToken, config.ServiceAccount)
}

func GcloudExec(args []string) (output string, err error) {
	cmd := exec.Command("gcloud", args...)
	wifConfig, err := NewWorkloadIdentityConfig()

	// Set the token if we have a WIF config.
	if err == nil && wifConfig != nil {
		gwifToken, err := GetGCPToken(wifConfig)
		if err == nil && gwifToken != nil && gwifToken.AccessToken != "" {
			utils.Debug("Using WIF token for gcloud command")
			cmd.Env = append(os.Environ(), "CLOUDSDK_AUTH_ACCESS_TOKEN="+gwifToken.AccessToken)
		}
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}
	return stdout.String(), nil
}

func GcloudAuth(shellExecutor utils.Executor, wifConfig *WorkloadIdentityConfig) (string, error) {
	audience := wifConfig.getAudience()

	// remove //iam.googleapis.com/ from audience.
	audience = strings.Replace(audience, "//iam.googleapis.com/", "", 1)
	oidcToken := wifConfig.GitLabOIDCToken.FromEnv
	if oidcToken == "" {
		return "", fmt.Errorf("GITLAB_OIDC_TOKEN is not set or empty")
	}

	// generate an empty temporary file for the OIDC token
	tmpFile, err := utils.WriteTempFile("oidc_token_*.jwt", oidcToken)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpFile.Close()

	// Create a second temporary file for the credential config
	credFile, err := utils.WriteTempFile("gcloud_cred_*.json", "")
	if err != nil {
		return "", fmt.Errorf("failed to create credential file: %w", err)
	}
	credFile.Close()

	// Create cred config command.
	_, err = shellExecutor.Run("gcloud", "iam", "workload-identity-pools", "create-cred-config", audience,
		"--service-account", wifConfig.ServiceAccount,
		"--output-file", credFile.Name(),
		"--credential-source-file", tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to create cred config: %w", err)
	}

	// Now login using the credential file
	out, err := shellExecutor.Run("gcloud", "auth", "login", "--cred-file", credFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}
	return out, nil
}

func CheckGcloudInstalled(shellExecutor utils.Executor) (bool, error) {
	_, err := shellExecutor.Run("gcloud", "--version")
	if err != nil {
		return false, err
	}
	return true, nil
}
