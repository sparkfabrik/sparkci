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
	gcloudExecCmd := NewGcloudExecCommand()
	gcloudAuthCmd := NewGcloudAuthCommand(nil)
	printGitlabJwtCmd := NewPrintGitlabJwtCommand()
	printSaTokenCmd := NewPrintSaTokenCommand()
	statusCmd := NewStatusCommand()
	printVarsCmd := NewPrintVarsCommand()
	GwifCommand.AddCommand(statusCmd)
	GwifCommand.AddCommand(printVarsCmd)
	GwifCommand.AddCommand(printSaTokenCmd)
	GwifCommand.AddCommand(gcloudExecCmd)
	GwifCommand.AddCommand(gcloudAuthCmd)
	GwifCommand.AddCommand(printGitlabJwtCmd)
}
