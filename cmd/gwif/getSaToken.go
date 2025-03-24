package cmd_gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

var getSaTokenCmd = &cobra.Command{
	Use:   "get-gwif-sa-token",
	Short: "Get Google Workload Identify token",
	Long:  `Get Google Workload Identity Federation token for the service account.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := gwif.NewWorkloadIdentityConfig()
		if err != nil {
			utils.Error(err.Error())
			return
		}
		token, err := gwif.GetGCPToken(config)
		if err != nil {
			utils.Error(err.Error())
			return
		}
		fmt.Println(token.AccessToken)
	},
}
