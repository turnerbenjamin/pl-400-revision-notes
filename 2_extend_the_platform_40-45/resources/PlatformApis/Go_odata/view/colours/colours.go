// Package colours provides ANSI color code functionality for terminal text
// formatting.
// It offers predefined color constants and utilities for applying colours to
// strings.
package colours

import "fmt"

// Colour represents an ANSI color code string that can be used to format
// terminal text. When printed before text content, it changes the display
// properties until reset.
type Colour string

const (
	Purple         Colour = "\033[38;5;127m" // Purple foreground color
	Red            Colour = "\033[38;5;196m" // Red foreground color
	Blue           Colour = "\033[38;5;81m"  // Blue foreground color
	Orange         Colour = "\033[38;5;208m" // Orange foreground color
	Grey           Colour = "\033[38;5;238m" // Grey foreground color
	Green          Colour = "\033[38;5;120m" // Green foreground color
	Reset          Colour = "\033[0m"        // Resets all color formatting
	BlueBackground Colour = "\033[48;5;166m" // Blue background color
)

// ApplyColour wraps the provided string with the specified color code and a
// reset code. This ensures the color only applies to the intended text and
// doesn't affect following output.
//
// Parameters:
//   - s: The string to be coloured
//   - colour: The Colour to apply
//
// Returns:
//   - A string with color codes applied
func ApplyColour(s string, colour Colour) string {
	return fmt.Sprintf("%s%s%s", colour, s, Reset)
}
