package main

import (
	"flag"
	"fmt"
	"os"

	"tablefy/internal/app"
)

func main() {
	autoExpand := flag.Bool("auto-expand", false, "Auto-expand focused column if it contains truncated cells")
	flag.Bool("a", false, "Short flag for auto-expand")
	flag.Parse()

	if err := app.Run(app.Config{AutoExpand: *autoExpand}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
