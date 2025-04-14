package utilities

import (
	"os"

	"golang.org/x/term"
)

func GetConsoleWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return w
}
