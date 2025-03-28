package format

import (
	"github.com/spf13/cobra"
)

// FormatCommand is the main command for formatting GitLab CI output
var FormatCommand = newFormatCommand()

func newFormatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format",
		Short: "Format GitLab CI output with sections and banners",
		Long: `Format GitLab CI output with different formatting options:

- Create collapsible sections in GitLab CI output
- Print banners to highlight important information
- Format output to improve readability in GitLab CI logs`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(newSectionCommand())
	cmd.AddCommand(newBannerCommand())

	return cmd
}
