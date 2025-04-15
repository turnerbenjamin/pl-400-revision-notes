// Package ansi provides ANSI escape code constants for terminal control
// operations.
// These codes can be used to control cursor behavior, clear the screen,
// and manage terminal output formatting.
package ansi

const (

	// Cursor control sequences
	CursorHide = "\033[?25l" // Hide terminal cursor
	CursorShow = "\033[?25h" // Show terminal cursor
	CursorHome = "\033[H"    //Move cursor to top-left (1,1)

	// Screen clearing sequences
	ClearScreen     = "\033[2J" // Clear entire screen
	ClearScrollback = "\033[3J" // Clear scrollback buffer
	ClearToEnd      = "\033[J"  // Clear from cursor position to end of screen

	// Combined operations
	ClearAll  = "\033[H\033[2J\033[3J" // Clear screen and scrollback buffer
	ResetView = "\033[H\033[J"         // Clear screen, not scrollback buffer
)
