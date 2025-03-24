package gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/spf13/cobra"
)

var format string

var printGitlabJwtCmd = &cobra.Command{
	Use:           "print-gitlab-jwt",
	Short:         "Print the Gitlab OIDC JWT token",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		oidc, err := gwif.NewGitlabOidc()
		if err != nil {
			return err
		}
		jwt := oidc.Payload

		switch format {
		case "json":
			res, err := jwt.JsonPrettyPrint()
			if err != nil {
				return err
			}
			fmt.Println(res)
		case "text":
			res := jwt.AsString()
			fmt.Println(res)
		default:
			return fmt.Errorf("unsupported format: %s. Supported formats are: json, text", format)
		}
		return nil
	},
}

func init() {
	printGitlabJwtCmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, text)")
}
