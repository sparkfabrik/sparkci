package cmd

import (
	"fmt"

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
	SilenceErrors: true,
	SilenceUsage:  true,
}

func SetVersionInfo(version, commit, date, builtBy string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s - built by %s)", version, date, commit, builtBy)
}

func Execute() {
	utils.InitLogging()

	if err := rootCmd.Execute(); err != nil {
		utils.Fatal(fmt.Sprintf("error: %s", err.Error()))
	}
}

func init() {
	rootCmd.AddCommand(gitlab.GitlabCommand)
	rootCmd.AddCommand(gwif.GwifCommand)
}
