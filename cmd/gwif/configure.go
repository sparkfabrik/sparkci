package gwif

import (
	"fmt"

	formatCmd "github.com/sparkfabrik/sparkci/cmd/format"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

// newConfigureCommand creates the `configure` subcommand for gwif.
func newConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "Configure Workload Identity Federation (WIF)",
		Long: `This command orchestrates the setup of Workload Identity Federation (WIF)
by running all necessary steps, including printing variables, checking status,
and authenticating with Google Cloud.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Step: Display the banner
			bannerCmd := formatCmd.NewBannerCommand()
			bannerCmd.SetArgs([]string{"--text", "GCP WIF CONFIGURATION"})
			if err := bannerCmd.Execute(); err != nil {
				utils.Error("%v", err)
				return fmt.Errorf("")
			}

			// Ensure the banner is always closed
			defer func() {
				bannerCmd.SetArgs([]string{"--text", "END GCP WIF CONFIGURATION"})
				if err := bannerCmd.Execute(); err != nil {
					utils.Warn("Failed to display end banner: %v", err)
				}
			}()

			// Step: Execute the `print-vars` command
			printVarsCmd := NewPrintVarsCommand()
			if err := printVarsCmd.Execute(); err != nil {
				utils.Error("%v", err)
				return fmt.Errorf("")
			}

			// Step: Execute the `status` command
			statusCmd := NewStatusCommand()
			statusCmd.SetArgs([]string{"--silent=false"})
			if err := statusCmd.Execute(); err != nil {
				utils.Error("%v", err)
				return fmt.Errorf("")
			}

			// Step: Execute the `gcloud-auth` command
			gcloudAuthCmd := NewGcloudAuthCommand(nil)
			if err := gcloudAuthCmd.Execute(); err != nil {
				utils.Error("%v", err)
				return fmt.Errorf("")
			}

			return nil
		},
	}
}
