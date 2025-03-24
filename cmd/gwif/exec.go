package cmd_gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec -- [gcloud commands and arguments]",
	Short: "Execute a gcloud command with WIF authentication",
	Long: `Execute a gcloud command with Google Cloud Workload Identity Federation authentication.
The -- separator MUST be used to separate sparkci command from gcloud arguments.

Example:
  sparkci gwif exec -- secrets versions access latest --project="my-project"`,
	// Don't validate args so we can handle them ourselves
	DisableFlagParsing: true,
	SilenceErrors:      true,
	SilenceUsage:       true,
	SuggestFor:         []string{"exec", "gwif"},
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
			utils.Error("The -- separator is required to separate sparkci command from gcloud arguments.")
			utils.Error("Example: sparkci gwif exec -- secrets versions access latest --project=\"my-project\"")
			return nil
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
