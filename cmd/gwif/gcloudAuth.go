package gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

// GcloudAuthOptions contains the dependencies for the auth command
type GcloudAuthOptions struct {
	ShellExecutor          utils.Executor
	WorkloadIdentityConfig func() (*gwif.WorkloadIdentityConfig, error)
	AuthFunc               func(utils.Executor, *gwif.WorkloadIdentityConfig) (string, error)
	GcloudCheck            func(utils.Executor) (bool, error)
}

// NewGcloudAuthCommand creates a new command to authenticate gcloud with WIF
func NewGcloudAuthCommand(opts *GcloudAuthOptions) *cobra.Command {
	// If no options provided, use defaults
	if opts == nil {
		opts = &GcloudAuthOptions{
			ShellExecutor:          utils.NewShellExecutor(),
			WorkloadIdentityConfig: gwif.NewWorkloadIdentityConfig,
			AuthFunc:               gwif.GcloudAuth,
			GcloudCheck:            gwif.CheckGcloudInstalled,
		}
	}

	var gcloudAuth = &cobra.Command{
		Use:                "gcloud-auth",
		Short:              "Authenticate gcloud cli with WIF",
		SilenceErrors:      true,
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := opts.GcloudCheck(opts.ShellExecutor)
			if err != nil {
				return err
			}

			wifConfig, err := opts.WorkloadIdentityConfig()
			if err != nil {
				return err
			}

			_, err = opts.AuthFunc(opts.ShellExecutor, wifConfig)
			if err != nil {
				return err
			}

			fmt.Println("gcloud auth activated with WIF")
			return nil
		},
	}

	return gcloudAuth
}
