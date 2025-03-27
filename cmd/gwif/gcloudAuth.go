package gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

var gcloudAuth = &cobra.Command{
	Use:                "gcloud-auth",
	Short:              "Authenticate gcloud cli with WIF",
	SilenceErrors:      true,
	DisableFlagParsing: true,
	SilenceUsage:       true,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellExecutor := utils.NewShellExecutor()
		_, err := gwif.CheckGcloudInstalled(shellExecutor)
		if err != nil {
			return err
		}
		wifConfig, err := gwif.NewWorkloadIdentityConfig()
		if err != nil {
			return err
		}
		_, err = gwif.GcloudAuth(shellExecutor, wifConfig)
		if err != nil {
			return err
		}
		fmt.Println("gcloud auth activated with WIF")
		return nil
	},
}
