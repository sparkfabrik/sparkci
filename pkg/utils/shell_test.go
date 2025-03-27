package utils

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShellExecutor(t *testing.T) {
	executor := NewShellExecutor()
	assert.NotNil(t, executor, "NewShellExecutor should return a non-nil Executor")
}

func TestShellExecutor_Command(t *testing.T) {
	executor := NewShellExecutor()

	// get the command path for echo.
	helloPath, _ := executor.Run("which", "echo")
	helloPath = strings.TrimSpace(helloPath)

	cmd := executor.Command("echo", "hello")
	assert.NotNil(t, cmd, "Command should return a non-nil *exec.Cmd")
	assert.Equal(t, helloPath, cmd.Path, "Command should set the correct command name")
	assert.Equal(t, []string{"hello"}, cmd.Args[1:], "Command should set the correct arguments")
}

func TestShellExecutor_RunCommand_Success(t *testing.T) {
	executor := NewShellExecutor()
	cmd := executor.Command("echo", "hello")
	output, err := executor.RunCommand(cmd)

	assert.NoError(t, err, "RunCommand should not return an error for a valid command")
	assert.Equal(t, "hello\n", output, "RunCommand should return the correct output")
}

func TestShellExecutor_RunCommand_Error(t *testing.T) {
	executor := NewShellExecutor()
	cmd := executor.Command("invalid_command")
	_, err := executor.RunCommand(cmd)

	assert.Error(t, err, "RunCommand should return an error for an invalid command")
}

func TestShellExecutor_Run_Success(t *testing.T) {
	executor := NewShellExecutor()
	output, err := executor.Run("echo", "hello")

	assert.NoError(t, err, "Run should not return an error for a valid command")
	assert.Equal(t, "hello\n", output, "Run should return the correct output")
}

func TestShellExecutor_Run_Error(t *testing.T) {
	executor := NewShellExecutor()
	_, err := executor.Run("invalid_command")

	assert.Error(t, err, "Run should return an error for an invalid command")
}

func TestShellExecutor_RunCommand_Stderr(t *testing.T) {
	executor := NewShellExecutor()
	cmd := executor.Command("ls", "nonexistent_file")
	_, err := executor.RunCommand(cmd)

	var execErr *exec.ExitError
	assert.True(t, errors.As(err, &execErr), "RunCommand should return an error containing *exec.ExitError for a failing command")
}
