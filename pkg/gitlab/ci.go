package gitlab

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/sparkfabrik/sparkci/pkg/utils"
)

// CIEnvironment represents the GitLab CI environment variables
type CIEnvironment struct {
	ProjectPath       string `json:"project_path" human:"Project Path" env:"CI_PROJECT_PATH"`
	ProjectURL        string `json:"project_url" human:"Project URL" env:"CI_PROJECT_URL"`
	CommitSHA         string `json:"commit_sha" human:"Commit SHA" env:"CI_COMMIT_SHA"`
	CommitBranch      string `json:"commit_branch" human:"Commit Branch" env:"CI_COMMIT_BRANCH"`
	CommitTag         string `json:"commit_tag" human:"Commit Tag" env:"CI_COMMIT_TAG"`
	JobName           string `json:"job_name" human:"Job Name" env:"CI_JOB_NAME"`
	PipelineID        string `json:"pipeline_id" human:"Pipeline ID" env:"CI_PIPELINE_ID"`
	JobID             string `json:"job_id" human:"Job ID" env:"CI_JOB_ID"`
	RunnerDescription string `json:"runner_description" human:"Runner Description" env:"CI_RUNNER_DESCRIPTION"`
	ProjectVisibility string `json:"project_visibility" human:"Project Visibility" env:"CI_PROJECT_VISIBILITY"`
	GitlabUserName    string `json:"gitlab_user_name" human:"GitLab User Name" env:"GITLAB_USER_NAME"`
	GitlabUserEmail   string `json:"gitlab_user_email" human:"GitLab User Email" env:"GITLAB_USER_EMAIL"`
	GitlabUserID      string `json:"gitlab_user_id" human:"GitLab User ID" env:"GITLAB_USER_ID"`
}

// GetCIEnvironment returns the current GitLab CI environment
func GetCIEnvironment() (*CIEnvironment, error) {
	if os.Getenv("GITLAB_CI") == "" {
		return nil, fmt.Errorf("not running in GitLab CI")
	}

	// Create it empty.
	env := &CIEnvironment{}

	// Use reflection to populate the struct fields from environment variables
	t := reflect.TypeOf(*env)
	v := reflect.ValueOf(env)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envVar := field.Tag.Get("env")
		if envVar == "" {
			return nil, fmt.Errorf("missing env tag for field %s", field.Name)
		}
		envValue := os.Getenv(envVar)
		if envValue != "" {
			fieldValue := v.Elem().Field(i)
			// env variables are always strings.
			fieldValue.SetString(envValue)
		}
	}
	return env, nil
}

// ToMap converts the stuct to a map using human-readable keys.
func (env *CIEnvironment) ToMap() map[string]string {
	result := make(map[string]string)

	t := reflect.TypeOf(*env)
	v := reflect.ValueOf(*env)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		humanName := field.Tag.Get("human")
		if humanName == "" {
			humanName = field.Name // Fall back to field name if no human tag
		}

		// Handle different types of values
		var valStr string
		if field.Type.Kind() == reflect.Bool {
			valStr = fmt.Sprintf("%v", v.Field(i).Bool())
		} else {
			valStr = v.Field(i).String()
		}

		result[humanName] = valStr
	}

	return result
}

// GetJobURL returns the URL of the current job
func (env *CIEnvironment) GetJobURL() string {
	projectURL := env.ProjectURL
	jobID := env.JobID

	if projectURL == "" || jobID == "" {
		return ""
	}

	return fmt.Sprintf("%s/-/jobs/%s", projectURL, jobID)
}

// GetPipelineURL returns the URL of the current pipeline
func (env *CIEnvironment) GetPipelineURL() string {
	projectURL := env.ProjectURL
	pipelineID := env.PipelineID

	if projectURL == "" || pipelineID == "" {
		return ""
	}

	return fmt.Sprintf("%s/-/pipelines/%s", projectURL, pipelineID)
}

// PrintEnvironment outputs the current GitLab environment in the specified format
func PrintEnvironment(format string) error {
	env, err := GetCIEnvironment()
	if env == nil {
		return err
	}
	envMap := env.ToMap()

	switch format {
	case "json":
		data, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml", "yml":
		for key, value := range envMap {
			fmt.Printf("%s: %s\n", key, value)
		}
	default:
		if jobURL := env.GetJobURL(); jobURL != "" {
			envMap["Job URL"] = jobURL
		}
		if pipelineURL := env.GetPipelineURL(); pipelineURL != "" {
			envMap["Pipeline URL"] = pipelineURL
		}

		// Add header with job info
		fmt.Printf("🔄 GitLab CI Job: %s (Project: %s)\n\n", env.JobName, env.ProjectPath)

		// Print using simple key-value formatting with tabs
		utils.PrintMap(envMap)
	}

	return nil
}
