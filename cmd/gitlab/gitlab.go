package cmd_gitlab

import (
	"github.com/spf13/cobra"
)

var GitlabCommand = &cobra.Command{
	Use:   "gitlab",
	Short: "GitLab CI utilities",
	Long:  `Commands for working with GitLab CI environments and operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GitlabCommand.AddCommand(printEnv)
}
