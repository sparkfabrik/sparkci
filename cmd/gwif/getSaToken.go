package gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/spf13/cobra"
)

var getSaTokenCmd = &cobra.Command{
	Use:           "print-sa-token",
	Short:         "Print Google Workload Identify token",
	Long:          `Print Google Workload Identity Federation token for the service account.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := gwif.NewWorkloadIdentityConfig()
		if err != nil {
			return err
		}
		token, err := gwif.GetGCPToken(config)
		if err != nil {
			return err
		}
		fmt.Println(token.AccessToken)
		return nil
	},
}
