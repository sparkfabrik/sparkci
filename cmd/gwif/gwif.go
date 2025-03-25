package gwif

import (
	"github.com/spf13/cobra"
)

var GwifCommand = &cobra.Command{
	Use:   "gwif",
	Short: "Google Cloud Workload Identity Federation utilities",
	Long:  `Commands for working with Google Cloud Workload Identity Federation in GitLab CI.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	Args: cobra.NoArgs,
}

func init() {
	GwifCommand.AddCommand(gcloudExec)
	GwifCommand.AddCommand(printGitlabJwtCmd)
	GwifCommand.AddCommand(getSaTokenCmd)
	GwifCommand.AddCommand(statusCmd)
	GwifCommand.AddCommand(printVarsCmd)
}
