package gwif

import (
	formatCmd "github.com/sparkfabrik/sparkci/cmd/format"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Step 1: Start the section
			sectionCmd := formatCmd.NewSectionCommand()
			sectionCmd.SetArgs([]string{"--title", "wif", "--description", "Workload Identity Federation"})
			if err := sectionCmd.Execute(); err != nil {
				return err
			}

			// Step 2: Display the banner
			bannerCmd := formatCmd.NewBannerCommand()
			bannerCmd.SetArgs([]string{"--text", "GCP WIF CONFIGURATION"})
			if err := bannerCmd.Execute(); err != nil {
				return err
			}

			// Step 3: Execute the `print-vars` command
			printVarsCmd := NewPrintVarsCommand()
			if err := printVarsCmd.Execute(); err != nil {
				return err
			}

			// Step 4: Execute the `status` command
			statusCmd := NewStatusCommand()
			statusCmd.SetArgs([]string{"--silent=false"})
			if err := statusCmd.Execute(); err != nil {
				return err
			}

			// Step 5: Execute the `gcloud-auth` command
			gcloudAuthCmd := NewGcloudAuthCommand(nil)
			if err := gcloudAuthCmd.Execute(); err != nil {
				return err
			}

			// Step 6: Display the end banner
			bannerCmd.SetArgs([]string{"--text", "END GCP WIF CONFIGURATION"})
			if err := bannerCmd.Execute(); err != nil {
				return err
			}

			// Step 7: End the section
			sectionCmd.SetArgs([]string{"--title", "wif", "--end"})
			if err := sectionCmd.Execute(); err != nil {
				return err
			}

			return nil
		},
	}
}
