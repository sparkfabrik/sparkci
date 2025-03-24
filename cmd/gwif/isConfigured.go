package gwif

import (
	"os"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	print  bool
	silent bool
)

var isConfiguredCmd = &cobra.Command{
	Use:   "is-configured",
	Short: "Check if the Workload Identity Federation is configured",
	Long:  `Check if the Workload Identity Federation is configured in the current environment.`,
	Example: `
# Check if the Workload Identity Federation is configured, it will just return 0 or 1 as exit code.
sparkci gwif is-configured

# Check if the Workload Identity Federation is configured, it will print the error message if not configured.
sparkci gwif is-configured --silent=false
	`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		wifConfig, err := gwif.NewWorkloadIdentityConfig()
		if err != nil {
			if silent {

				os.Exit(1)
			}
			return err
		}

		if !print {
			return nil
		}

		wifMap := wifConfig.SafeToMap()
		utils.PrintMap(wifMap)
		return nil
	},
}

var printVarsCmd = &cobra.Command{
	Use:   "print-vars",
	Short: "Print the Workload Identity Federation standard env variables",
	Long: `
Print the Workload Identity Federation standard env variables: GCP_WIF_PROJECT_ID, GCP_WIF_POOL, GCP_WIF_PROVIDER, GCP_WIF_SERVICE_ACCOUNT_EMAIL.`,
	Example: `
# Print the Workload Identity Federation configuration in formatted output.
sparkci gwif print-vars`,
	Run: func(cmd *cobra.Command, args []string) {
		envs := map[string]string{
			"GCP_WIF_PROJECT_ID":            os.Getenv("GCP_WIF_PROJECT_ID"),
			"GCP_WIF_POOL":                  os.Getenv("GCP_WIF_POOL"),
			"GCP_WIF_PROVIDER":              os.Getenv("GCP_WIF_PROVIDER"),
			"GCP_WIF_SERVICE_ACCOUNT_EMAIL": os.Getenv("GCP_WIF_SERVICE_ACCOUNT_EMAIL"),
		}
		utils.PrintFormattedVars("Configured WIF related variables", envs)
	},
}

func init() {
	isConfiguredCmd.Flags().BoolVarP(&print, "print", "p", false, "Print configuration")
	isConfiguredCmd.Flags().BoolVarP(&silent, "silent", "s", true, "Silent mode (no output)")
}
