package main

import (
	"fmt"
	"os"

	"github.com/sparkfabrik/sparkci/cmd/sparkci/commands"
	"github.com/sparkfabrik/sparkci/pkg/utils"
)

var (
	// Version information - will be set by goreleaser
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Initialize logging
	utils.InitLogging()

	// Create root command with version information
	rootCmd := commands.NewRootCmd(version)

	// Add version information
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s, builder: %s)",
		version, commit, date, builtBy)

	if err := rootCmd.Execute(); err != nil {
		utils.Error("Command execution failed: %v", err)
		os.Exit(1)
	}
}
