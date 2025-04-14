// Package utilities provides helper functions for common operations across the
// application.
package utilities

import (
	"os"

	"golang.org/x/term"
)

// GetConsoleWidth returns the width (number of columns) of the current terminal.
// If the terminal dimensions cannot be determined, it returns the provided
// defaultValue.
//
// Parameters:
//   - defaultValue: The fallback value to return if terminal width cannot be
//     determined
//
// Returns:
//   - The width of the terminal in columns, or defaultValue if an error occurs
func GetConsoleWidth(defaultValue int) int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return defaultValue
	}
	return w
}

// GetConsoleHeight returns the height (number of rows) of the current terminal.
// If the terminal dimensions cannot be determined, it returns the provided
// defaultValue.
//
// Parameters:
//   - defaultValue: The fallback value to return if terminal height cannot be
//     determined
//
// Returns:
//   - The height of the terminal in rows, or defaultValue if an error occurs
func GetConsoleHeight(defaultValue int) int {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return defaultValue
	}
	return h
}
