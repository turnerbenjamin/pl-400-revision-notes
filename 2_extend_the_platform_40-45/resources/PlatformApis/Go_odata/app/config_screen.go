// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	authMode "github.com/turnerbenjamin/go_odata/constants/authmode"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// GetConfigScreen creates and returns a configuration screen for the
// application.
// The screen allows the user to select an authentication mode (Application or
// User) for connecting to the Dataverse API.
//
// The screen includes:
// - A title "CONFIGURATION" in purple color
// - Instructional text "Select authentication mode"
// - A menu with authentication mode options
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if menu component creation fails
func GetConfigScreen() (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(authMode.Application),
		string(authMode.User),
	},
	)

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("CONFIGURATION", colours.Purple),
		view.NewTextComponent("Select authentication mode"),
		menu,
	})
}
