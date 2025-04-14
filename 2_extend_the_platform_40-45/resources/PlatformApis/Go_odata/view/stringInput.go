// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/constants/ansi"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// stringInput represents a text input field in a terminal UI.
// It handles user text entry, validation for required fields,
// and displays error messages when validation fails.
type stringInput struct {
	propertyName string
	value        string
	isRequired   bool
	errorMessage string
	requiredFlag string
}

// NewStringInputComponent creates a new text input field with the given
// property name. If isRequired is true, the field will be marked as required
// and validated before submission.
func NewStringInputComponent(propertyName, value string, isRequired bool) InteractiveComponent {
	si := &stringInput{
		propertyName: propertyName,
		value:        value,
		isRequired:   isRequired,
	}
	if isRequired {
		si.requiredFlag = "(" + colours.ApplyColour("*", colours.Red) + ")"
	}
	return si
}

// render displays the string input component with its current value.
// It shows the property name, current input value, and any error messages.
func (si *stringInput) render() {
	fmt.Print(ansi.CursorShow)
	fmt.Printf("\n%s%s: %s", si.propertyName, si.requiredFlag, si.value)
	if si.errorMessage != "" {
		fmt.Printf("\n\n%s%s%s", colours.Red, si.errorMessage, colours.Reset)
	}
}

// handleKeyboardInput routes keypresses to the appropriate handlers.
// This component handles all validation errors via UI and never returns
// actual errors through the error return value. The error return is
// maintained for compatibility with the InteractiveComponent interface,
// which is used by components that may encounter technical failures.
func (si *stringInput) handleKeyboardInput(c rune, k keyboard.Key) (*updateResponse, error) {
	si.clearErrorMessage()
	switch k {
	case keyboard.KeyEnter:
		return si.handleEnterPressed(), nil
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		return si.handleBackspacePressed(), nil
	case keyboard.KeySpace:
		return si.handleCharEntered(' '), nil
	default:
		return si.handleCharEntered(c), nil
	}
}

// handleEnterPressed processes when the Enter key is pressed.
// It returns an updateResponse with the captured input value, or displays
// an error message if validation fails for required fields.
func (si *stringInput) handleEnterPressed() *updateResponse {
	if si.isRequired && si.value == "" {
		si.showPropertyRequiredError()
		return newUpdateResponse().setContinue(true)
	}
	return newUpdateResponse().setUserInput(si.value)
}

// handleBackspacePressed removes the last character from the input value
// and returns an updateResponse to continue input capture.
func (si *stringInput) handleBackspacePressed() *updateResponse {
	if len(si.value) > 0 {
		si.value = si.value[:len(si.value)-1]
	}
	return newUpdateResponse().setContinue(true).setFullRefresh()
}

// handleCharEntered appends the given character to the input value
// and returns an updateResponse to continue input capture.
func (si *stringInput) handleCharEntered(c rune) *updateResponse {
	si.value += string(c)
	return newUpdateResponse().setContinue(true)
}

// showPropertyRequiredError updates the error message property to inform the
// user that they must enter a value to continue
func (si *stringInput) showPropertyRequiredError() {
	errorText := fmt.Sprintf("%s is required", si.propertyName)
	si.errorMessage = colours.ApplyColour(errorText, colours.Red)
}

// clearErrorMessage resets the error state of the input field.
func (si *stringInput) clearErrorMessage() {
	si.errorMessage = ""
}
