package main

import (
	"flag"
	"fmt"
	"os"

	"tablefy/internal/app"
)

func main() {
	autoExpand := flag.Bool("auto-expand", false, "Auto-expand focused column if it contains truncated cells")
	shortAutoExpand := flag.Bool("a", false, "Short flag for auto-expand")
	flag.Parse()

	// Support both flags
	enableAutoExpand := *autoExpand || *shortAutoExpand

	if err := app.Run(app.Config{AutoExpand: enableAutoExpand}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
