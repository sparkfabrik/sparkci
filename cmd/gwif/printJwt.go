package cmd_gwif

import (
	"fmt"

	"github.com/sparkfabrik/sparkci/pkg/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

var format string

var printJwtCmd = &cobra.Command{
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

func init() {
	printJwtCmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, text)")
}
