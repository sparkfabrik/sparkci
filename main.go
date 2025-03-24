package main

import (
	"github.com/sparkfabrik/sparkci/cmd"
)

var (
	// Version information - will be set by goreleaser
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date, builtBy)
	cmd.Execute()
}
