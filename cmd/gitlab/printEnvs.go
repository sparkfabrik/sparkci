package gitlab

import (
	"github.com/sparkfabrik/sparkci/pkg/gitlab"
	"github.com/spf13/cobra"
)

var format string

var printEnvs = &cobra.Command{
	Use:           "print-envs",
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

func init() {
	printEnvs.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json, yaml)")
}
