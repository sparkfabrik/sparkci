package gitlab

import (
	formatCmd "github.com/sparkfabrik/sparkci/cmd/gitlab/format"
	"github.com/spf13/cobra"
)

var GitlabCommand = &cobra.Command{
	Use:   "gitlab",
	Short: "GitLab CI utilities",
	Long:  `Commands for working with GitLab CI environments and operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	Args: cobra.NoArgs,
}

func init() {
	cmd := formatCmd.NewFormatCommand()
	GitlabCommand.AddCommand(printEnvs)
	GitlabCommand.AddCommand(cmd)
}
