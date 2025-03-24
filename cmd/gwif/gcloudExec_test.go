package gwif

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expected    []string
		expectedErr string
	}{
		{
			name:        "No separator",
			args:        []string{"secrets", "versions", "access", "latest"},
			expected:    nil,
			expectedErr: "missing -- separator. The -- separator MUST be used to separate sparkci command from gcloud arguments",
		},
		{
			name:        "No gcloud command after separator",
			args:        []string{"--"},
			expected:    nil,
			expectedErr: "no gcloud command provided to execute after the -- separator",
		},
		{
			name:        "Valid gcloud command",
			args:        []string{"--", "secrets", "versions", "access", "latest"},
			expected:    []string{"secrets", "versions", "access", "latest"},
			expectedErr: "",
		},
		{
			name:        "Valid gcloud command with additional args",
			args:        []string{"--", "secrets", "versions", "access", "latest", "--project=my-project"},
			expected:    []string{"secrets", "versions", "access", "latest", "--project=my-project"},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateArgs(tt.args)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
