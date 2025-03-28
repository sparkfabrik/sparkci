package format

import (
	"errors"

	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

func newSectionCommand() *cobra.Command {
	var sectionTitle string
	var sectionDescription string
	var endSection bool

	cmd := &cobra.Command{
		Use:   "section",
		Short: "Create collapsible sections in GitLab CI output",
		Long: `Create or end collapsible sections in GitLab CI output.

Examples:
  # Start a section
  sparkci format section --title "build-logs" --description "Build logs for the project"

  # End a section
  sparkci format section --title "build-logs" --end`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if sectionTitle == "" {
				return errors.New("section title is required")
			}

			if endSection {
				utils.GitLabSectionEnd(sectionTitle)
				return nil
			}

			utils.GitLabSectionStart(sectionTitle, sectionDescription)
			return nil
		},
		SilenceErrors: true,
	}

	cmd.Flags().StringVarP(&sectionTitle, "title", "t", "", "Title for the section (required)")
	cmd.Flags().StringVarP(&sectionDescription, "description", "d", "", "Description for the section (defaults to title if not provided)")
	cmd.Flags().BoolVarP(&endSection, "end", "e", false, "End the section instead of starting it")
	cmd.MarkFlagRequired("title")

	return cmd
}
