package gwif

import (
	"os"

	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

func NewPrintVarsCommand() *cobra.Command {
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
	return printVarsCmd
}
