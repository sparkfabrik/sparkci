package commands

import (
	"github.com/spf13/cobra"
)

// RootCmdOptions contains options for the root command
type RootCmdOptions struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

// NewRootCmd creates the root command for sparkci
func NewRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "sparkci",
		Short:   "SparkCI - A CLI tool for GitLab CI operations",
		Version: version,
		Long: `SparkCI is a CLI tool designed to enhance GitLab CI workflows,
providing various utilities that can be run both in CI/CD pipelines and locally.`,
	}

	// Add all commands
	rootCmd.AddCommand(NewGwifCmd())
	rootCmd.AddCommand(NewGitlabCmd())

	return rootCmd
}
