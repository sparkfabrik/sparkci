package gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

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
		// Look for the -- separator
		var gcloudArgs []string
		separatorIndex := -1

		for i, arg := range args {
			if arg == "--" {
				separatorIndex = i
				break
			}
		}

		if separatorIndex == -1 {
			return cmd.Help()
		} else {
			// Skip the -- separator itself
			gcloudArgs = args[separatorIndex+1:]
		}

		if len(gcloudArgs) == 0 {
			utils.Error("No gcloud command provided to execute after the -- separator.")
			return nil
		}

		output, err := gwif.GcloudExec(gcloudArgs)
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil
	},
}
