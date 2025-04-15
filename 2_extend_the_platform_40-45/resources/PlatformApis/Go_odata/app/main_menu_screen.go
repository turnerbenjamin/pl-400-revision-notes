// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	mainMenuOption "github.com/turnerbenjamin/go_odata/constants/mainmenuoption"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// newMainMenuScreen creates the main menu screen for the application.
// It constructs a menu with options for different tables (Accounts, Contacts)
// and an Exit option.
//
// The screen includes:
// - A title "Table Selection" in purple color
// - Instructional text "Choose a table"
// - A menu with table and exit options
//
// Returns:
//   - A Screen object ready to be rendered
//   - An error if menu component creation fails
func newMainMenuScreen() (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(mainMenuOption.Accounts),
		string(mainMenuOption.Contacts),
		string(mainMenuOption.Exit),
	})

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("Table Selection", colours.Purple),
		view.NewTextComponent("Choose a table"),
		menu})
}
