package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"tablefy/internal/app"
)

// Version is set by ldflags during build
// Default value for local builds without version information
var Version = "dev"

// CommitHash is set by ldflags during build
// Default value for local builds without commit information
var CommitHash = "unknown"

func main() {
	version := pflag.BoolP("version", "v", false, "Show version information")
	autoExpand := pflag.BoolP("auto-expand", "a", false, "Auto-expand focused column if it contains truncated cells")
	pflag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("tablefy %s (commit: %s)\n", Version, CommitHash)
		os.Exit(0)
	}

	if err := app.Run(app.Config{AutoExpand: *autoExpand}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
