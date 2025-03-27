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

// NewStatusCommand creates a new command to check the status of Workload Identity Federation
func NewStatusCommand() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check Workload Identity Federation status",
		Long:  `Verify if Workload Identity Federation is properly configured in the current environment.`,
		Example: `
# Check if the Workload Identity Federation is configured, it will just return 0 or 1 as exit code.
sparkci gwif status

# Check if the Workload Identity Federation is configured, it will print the error message if not configured.
sparkci gwif status --silent=false
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
	statusCmd.Flags().BoolVarP(&print, "print", "p", false, "Print configuration")
	statusCmd.Flags().BoolVarP(&silent, "silent", "s", true, "Silent mode (no output)")
	return statusCmd
}
