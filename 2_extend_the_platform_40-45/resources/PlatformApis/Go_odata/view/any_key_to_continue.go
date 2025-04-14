// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

// anyKeyToContinue represents a simple interactive terminal component that
// waits for any user keypress before continuing program flow. It implements
// the InteractiveComponent interface.
type anyKeyToContinue struct {
}

// NewAnyKeyToContinueComponent creates and returns a new anyKeyToContinue
// component. This component displays a message prompting the user to press any
// key to continue and handles that input event.
// Returns an InteractiveComponent that can be used in a screen.
func NewAnyKeyToContinueComponent() InteractiveComponent {
	return &anyKeyToContinue{}
}

// render displays the "Press any key to continue" message in the terminal.
// This implements part of the InteractiveComponent interface.
func (t *anyKeyToContinue) render() {
	fmt.Print("\n\nPress any key to continue")
}

// handleKeyboardInput processes any keyboard input and signals completion.
// Any keypress will cause this component to complete its interaction.
// It returns an updateResponse that signals the component is finished and
// returns a nil error.
//
// Note that this method returns an error to comply with the
// InteractiveComponent interface defined in screen.go, though this
// implementation will never return an error.
//
// Parameters:
//   - char: The character code (rune) of the pressed key
//   - key: The keyboard key that was pressed
//
// Returns:
//   - *updateResponse: A response indicating the component is done
//   - error: Always nil in this implementation
func (t *anyKeyToContinue) handleKeyboardInput(char rune, key keyboard.Key) (*updateResponse, error) {
	return newUpdateResponse().setContinue(false), nil
}
