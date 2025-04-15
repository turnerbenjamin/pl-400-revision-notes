// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// newStringInputScreen creates a screen that prompts the user to input a string
// value.
// The screen includes a title, instructions text, and a text input field.
//
// Parameters:
//   - title: The title text to display at the top of the screen
//   - text: Instructions or explanation text to display
//   - propertyName: Name of the property being edited (shown as field label)
//   - value: The initial value to display in the input field
//   - isRequired: Whether the field must contain a value before proceeding
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if screen creation fails
func newStringInputScreen(title, text, propertyName, value string, isRequired bool) (view.Screen, error) {
	return view.MakeScreen([]view.Component{
		view.NewTitleComponent(title, colours.Purple),
		view.NewTextComponent(text),
		view.NewStringInputComponent(propertyName, value, isRequired),
	})
}
