package cmd

import (
	"github.com/sparkfabrik/sparkci/cmd/gitlab"
	"github.com/sparkfabrik/sparkci/cmd/gwif"
	"github.com/sparkfabrik/sparkci/pkg/utils"
	"github.com/spf13/cobra"
)

// RootCmdOptions contains options for the root command
type RootCmdOptions struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

var rootCmd = &cobra.Command{
	Use:   "sparkci",
	Short: "SparkCI - A CLI tool for GitLab CI operations",
	Long: `SparkCI is a CLI tool designed to enhance GitLab CI workflows,
providing various utilities that can be run both in CI/CD pipelines and locally.`,
}

func SetVersionInfo(version, commit, date, builtBy string) {
	rootCmd.Version = version
	rootCmd.Short += " (commit: " + commit + ", built: " + date + ", builder: " + builtBy + ")"
}

func Execute() {
	utils.InitLogging()

	if err := rootCmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}

func init() {
	rootCmd.AddCommand(gitlab.GitlabCommand)
	rootCmd.AddCommand(gwif.GwifCommand)
}
