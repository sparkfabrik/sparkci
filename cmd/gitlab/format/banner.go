package format

import (
	"errors"

	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

func NewBannerCommand() *cobra.Command {
	var bannerText string

	cmd := &cobra.Command{
		Use:   "banner",
		Short: "Print a banner in GitLab CI output",
		Long: `Print a highlighted banner text in GitLab CI output.

Example:
  sparkci format banner --text "Deployment Started"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if bannerText == "" {
				return errors.New("banner text is required")
			}

			utils.GitLabPrintBanner(bannerText)
			return nil
		},
		SilenceErrors: true,
	}

	cmd.Flags().StringVarP(&bannerText, "text", "t", "", "Text to display in the banner (required)")
	cmd.MarkFlagRequired("text")

	return cmd
}
