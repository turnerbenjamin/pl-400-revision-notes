// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"errors"

	"github.com/turnerbenjamin/go_odata/view/consoleinputreader"
)

// ErrNilScreen is returned when attempting to navigate to a nil screen.
// This error prevents the UI from entering an invalid state with no active
// screen.
var ErrNilScreen = errors.New("cannot navigate to nil screen")

// UI defines the interface for user interface controllers.
// It provides methods for screen navigation and application termination.
type UI interface {
	// NavigateTo changes the current screen to the provided Screen instance,
	// mounts it, and returns the screen's output after user interaction.
	NavigateTo(Screen) (ScreenOutput, error)

	// Exit performs cleanup operations including screen dismounting
	// and input reader closure.
	Exit()
}

// consoleUI implements the UI interface for terminal-based applications.
// It manages the currently active screen and handles keyboard input.
type consoleUI struct {
	// currentScreen holds the active screen being displayed to the user
	currentScreen Screen
	// inputReader provides terminal input capabilities
	inputReader consoleinputreader.InputReader
}

// NewConsoleUI creates and initializes a new console UI controller.
// It sets up the input reader and returns an implementation of the UI
// interface.
// Returns an error if the input reader cannot be initialized.
func NewConsoleUI() (UI, error) {
	ir := consoleinputreader.NewInputReader()
	err := ir.Open()
	if err != nil {
		return nil, err
	}

	ui := consoleUI{
		inputReader: ir,
	}
	return &ui, nil
}

// NavigateTo switches to a new screen within the UI.
// If there's a current screen, it will be dismounted first.
// The new screen is mounted and its output is processed through AwaitOutput.
//
// Parameters:
//   - s: The screen to navigate to. Must not be nil.
//
// Returns:
//   - ScreenOutput: The output from the screen after user interaction
//   - error: ErrNilScreen if s is nil, or any error that occurs during mounting
//     or user interaction with the screen
func (c *consoleUI) NavigateTo(s Screen) (ScreenOutput, error) {
	if s == nil {
		return nil, ErrNilScreen
	}

	if c.currentScreen != nil {
		c.currentScreen.Dismount()
	}

	c.currentScreen = s
	c.currentScreen.Mount()
	return c.AwaitOutput()
}

// Exit performs cleanup operations when exiting the application.
// It dismounts the current screen if one exists and closes the input reader.
//
// Parameters:
//   - None
//
// Returns:
//   - None
func (c *consoleUI) Exit() {
	if c.currentScreen != nil {
		c.currentScreen.Dismount()
	}
	c.inputReader.Close()
}

// AwaitOutput enters a processing loop for user input on the current screen.
// It continually reads input, passes it to the current screen for handling, and
// returns when the screen signals completion or an error occurs.
// The screen is refreshed after each input that doesn't result in navigation.
//
// Parameters:
//   - None
//
// Returns:
//   - ScreenOutput: The output from the screen after user interaction
//   - error: Any error that occurs during input handling
func (c *consoleUI) AwaitOutput() (ScreenOutput, error) {
	for {
		char, key, err := c.inputReader.AwaitInput()
		if err != nil {
			return nil, err
		}

		updateResponse, err := c.currentScreen.handleKeyboardInput(char, key)
		if err != nil {
			return nil, err
		}

		if !updateResponse.doContinue {
			return updateResponse, nil
		}
		c.currentScreen.Refresh()
	}
}
