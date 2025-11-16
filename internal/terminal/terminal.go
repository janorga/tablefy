package terminal

import (
	"os"

	"golang.org/x/term"
)

// GetSize gets the terminal width and height
func GetSize() (width, height int, err error) {
	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24, err // Default fallback
	}
	return width, height, nil
}
