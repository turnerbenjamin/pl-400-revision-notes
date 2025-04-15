// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// newSuccessScreen creates a screen that displays a success message to the
// user.
// The screen includes a success title, the provided message text, and a prompt
// for the user to press any key to continue.
//
// Parameters:
//   - msg: The success message to display
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if screen creation fails
func newSuccessScreen(msg string) (view.Screen, error) {
	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("SUCCESS", colours.Green),
		view.NewTextComponent(msg),
		view.NewAnyKeyToContinueComponent(),
	})
}
