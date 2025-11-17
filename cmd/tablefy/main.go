package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"tablefy/internal/app"
)

func main() {
	autoExpand := pflag.BoolP("auto-expand", "a", false, "Auto-expand focused column if it contains truncated cells")
	pflag.Parse()

	if err := app.Run(app.Config{AutoExpand: *autoExpand}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
