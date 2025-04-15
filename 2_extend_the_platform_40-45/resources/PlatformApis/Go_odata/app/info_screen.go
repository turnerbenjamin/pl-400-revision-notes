// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// newInfoScreen creates a screen that displays an informational message to the
// user.
// The screen includes an info title in blue, the provided message text, and a
// prompt for the user to press any key to continue.
//
// Parameters:
//   - msg: The informational message to display
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if screen creation fails
func newInfoScreen(msg string) (view.Screen, error) {
	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("INFO", colours.Blue),
		view.NewTextComponent(msg),
		view.NewAnyKeyToContinueComponent(),
	})
}
