package gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/spf13/cobra"
)

func NewGcloudExecCommand() *cobra.Command {
	var gcloudExec = &cobra.Command{
		Use:   "gcloud-exec -- [gcloud commands and arguments]",
		Short: "Execute a gcloud command with WIF authentication",
		Long: `Execute a gcloud command with Google Cloud Workload Identity Federation authentication.
	The -- separator MUST be used to separate sparkci command from gcloud arguments.

	Example:
	  sparkci gwif gcloud-exec -- secrets versions access latest --project="my-project"`,
		// Don't validate args so we can handle them ourselves
		DisableFlagParsing: true,
		SilenceErrors:      true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			gcloudArgs, err := validateArgs(args)
			if err != nil {
				return cmd.Help()
			}

			output, err := gwif.GcloudExec(gcloudArgs)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		},
	}
	return gcloudExec

}

func validateArgs(args []string) ([]string, error) {
	var gcloudArgs []string
	separatorIndex := -1

	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	if separatorIndex == -1 {
		return nil, fmt.Errorf("missing -- separator. The -- separator MUST be used to separate sparkci command from gcloud arguments")
	} else {
		gcloudArgs = args[separatorIndex+1:]
	}

	if len(gcloudArgs) == 0 {
		return nil, fmt.Errorf("no gcloud command provided to execute after the -- separator")
	}

	return gcloudArgs, nil
}
