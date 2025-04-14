// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"fmt"
)

// text represents a simple text component for displaying plain text content.
// It renders text followed by spacing for readability within the terminal UI.
type text struct {
	content string
}

// render displays the text component on the terminal.
// It prints the text content followed by additional newlines for spacing.
func (t *text) render() {
	fmt.Printf("%s\n\n", t.content)
}

// NewTextComponent creates a new text component with the specified content.
// It returns an implementation of the Component interface that can be added to
// a screen.
func NewTextComponent(str string) Component {
	return &text{
		content: str,
	}
}
