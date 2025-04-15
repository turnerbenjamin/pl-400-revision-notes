// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	confirmOption "github.com/turnerbenjamin/go_odata/constants/confirm_option"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// NewConfirmationScreen creates a screen that prompts the user to confirm an
// action.
// The screen includes a confirmation title in purple, the provided message
// explaining what action is being confirmed, and a Yes/No menu for user
// selection.
//
// Parameters:
//   - msg: The message explaining what action requires confirmation
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if menu component creation fails
func NewConfirmationScreen(msg string) (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(confirmOption.Yes),
		string(confirmOption.No),
	})

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("Confirmation Required", colours.Purple),
		view.NewTextComponent(msg),
		menu,
	},
	)
}
