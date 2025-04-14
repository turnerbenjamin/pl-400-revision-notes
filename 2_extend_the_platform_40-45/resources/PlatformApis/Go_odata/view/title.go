// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"fmt"
	"strings"

	"github.com/turnerbenjamin/go_odata/view/colours"
)

// title represents a styled text header component.
// It displays formatted text that has been converted to uppercase and
// coloured according to the specified colour parameter.
type title struct {
	content string // The pre-formatted content ready for display
}

// render displays the title component on the terminal.
// It prints the pre-formatted title content followed by newlines for spacing.
func (t *title) render() {
	fmt.Printf("%s\n\n", t.content)
}

// NewTitleComponent creates a new title component with the specified text and
// colour. The text is automatically converted to uppercase and formatted with
// the provided colour. It returns an implementation of the Component interface
// that can be added to a screen.
func NewTitleComponent(str string, colour colours.Colour) Component {
	uppercasedTitle := strings.ToUpper(str)
	formattedTitle := colours.ApplyColour(uppercasedTitle, colour)

	return &title{
		content: formattedTitle,
	}
}
