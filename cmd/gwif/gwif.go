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
}

func init() {
	GwifCommand.AddCommand(execCmd)
	GwifCommand.AddCommand(printJwtCmd)
	GwifCommand.AddCommand(getSaTokenCmd)
}
