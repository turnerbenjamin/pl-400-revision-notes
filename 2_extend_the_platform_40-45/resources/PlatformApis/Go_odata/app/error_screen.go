// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"strings"

	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// newErrorScreen creates a screen that displays an error message to the user.
// The screen includes an error title in red, the provided message text with
// enhanced readability (periods followed by line breaks), and a prompt for
// the user to press any key to continue.
//
// Parameters:
//   - msg: The error message to display
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if screen creation fails
func newErrorScreen(msg string) (view.Screen, error) {
	fm := strings.ReplaceAll(msg, ". ", ".\n")

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("ERROR", colours.Red),
		view.NewTextComponent(fm),
		view.NewAnyKeyToContinueComponent(),
	})
}
