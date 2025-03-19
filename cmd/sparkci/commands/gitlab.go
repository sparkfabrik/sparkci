package commands

import (
	"github.com/sparkfabrik/sparkci/pkg/gitlab"
	"github.com/spf13/cobra"
)

// NewGitlabCmd creates a new gitlab command
func NewGitlabCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "GitLab CI utilities",
		Long:  `Commands for working with GitLab CI environments and operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newGitlabPrintEnvCmd())
	return cmd
}

// newGitlabEnvCmd creates an env subcommand
func newGitlabPrintEnvCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:           "print-env",
		Short:         "Print GitLab CI environment information",
		Long:          `Display information about the current GitLab CI environment.`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := gitlab.PrintEnvironment(format); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json, yaml)")

	return cmd
}
