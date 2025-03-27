package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Executor interface defines methods to create and run commands
type Executor interface {
	Command(name string, args ...string) *exec.Cmd
	RunCommand(cmd *exec.Cmd) (string, error)
	Run(name string, args ...string) (string, error)
}

// ShellExecutor implements the Executor interface
type ShellExecutor struct {
	*exec.Cmd
}

// NewShellExecutor creates a new ShellExecutor
func NewShellExecutor() Executor {
	return &ShellExecutor{}
}

// Command creates a new exec.Cmd
func (e *ShellExecutor) Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// RunCommand executes a pre-created command and returns its output as a string
func (e *ShellExecutor) RunCommand(cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", stderr.String(), err)
	}

	return stdout.String(), nil
}

// Run creates and executes a command in one step
func (e *ShellExecutor) Run(name string, args ...string) (string, error) {
	cmd := e.Command(name, args...)
	return e.RunCommand(cmd)
}
