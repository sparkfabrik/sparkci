package commands

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

// NewGwifCmd creates a new gwif command
func NewGwifCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gwif",
		Short: "Google Cloud Workload Identity Federation utilities",
		Long:  `Commands for working with Google Cloud Workload Identity Federation in GitLab CI.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(newGwifExecCmd())
	cmd.AddCommand(newPrintJwtToken())
	cmd.AddCommand(newGetGoogleWifToken())

	return cmd
}

func newPrintJwtToken() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "print-jwt-token",
		Short: "Print jwt token",
		Long:  `Print the Gitlab OIDC JWT token.`,
		Run: func(cmd *cobra.Command, args []string) {
			oidc, err := gwif.NewGitlabOidc()
			if err != nil {
				utils.Error(err.Error())
				return
			}
			jwt := oidc.Payload

			switch format {
			case "json":
				res, err := jwt.JsonPrettyPrint()
				if err != nil {
					utils.Error(err.Error())
					return
				}
				fmt.Println(res)
			case "text":
				res := jwt.AsString()
				fmt.Println(res)
			default:
				utils.Error("Unsupported format. Supported formats are: json, text")
			}
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, text)")
	return cmd
}

func newGetGoogleWifToken() *cobra.Command {
	cmd := &cobra.Command{
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
	return cmd
}

// newGwifExecCmd creates an exec subcommand
func newGwifExecCmd() *cobra.Command {
	cmd := &cobra.Command{
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

	return cmd
}
