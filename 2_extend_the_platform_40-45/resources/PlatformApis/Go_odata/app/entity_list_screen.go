// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	tableMenuOption "github.com/turnerbenjamin/go_odata/constants/tablemenuoption"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// listControl implements the view.ListControl interface and represents
// a selectable option in the entity list screen.
type listControl struct {
	label string
	value string
	key   rune
}

// Label returns the human-readable text of the control to be displayed in the
// UI.
func (lc listControl) Label() string {
	return lc.label
}

// Value returns the internal value associated with this control.
func (lc listControl) Value() string {
	return lc.value
}

// Key returns the keyboard key associated with this control.
func (lc listControl) Key() rune {
	return lc.key
}

// entityListControls defines the standard set of controls available for all
// entity list screens in the application.
var entityListControls = []view.ListControl{
	listControl{
		label: "Set/Clear search term",
		value: string(tableMenuOption.Search),
		key:   's',
	},
	listControl{
		label: "Create",
		value: string(tableMenuOption.Create),
		key:   'c',
	},
	listControl{
		label: "Update",
		value: string(tableMenuOption.Update),
		key:   'u',
	},
	listControl{
		label: "Delete",
		value: string(tableMenuOption.Delete),
		key:   'd',
	},
	listControl{
		label: "Back to main menu",
		value: string(tableMenuOption.Back),
		key:   'b',
	},
}

// listScreenOptions contains configuration for creating an entity list screen.
// The generic type T represents the entity type to be displayed.
type listScreenOptions[T view.Entity] struct {
	entityList view.EntityList[T]
	columns    []view.ListColumn[T]
}

// newEntityListScreen creates a new screen displaying a list of entities.
// It takes an entityLabel string that will be used as the screen's title and
// listScreenOptions that configure the list display.
// Returns a view.Screen implementation and any error encountered during
// creation.
func newEntityListScreen[T view.Entity](entityLabel string, listScreenOptions listScreenOptions[T]) (view.Screen, error) {

	listOptions := view.ListComponentOptions[T]{
		Controls:   entityListControls,
		EntityList: listScreenOptions.entityList,
		Columns:    listScreenOptions.columns,
	}

	listComponent, err := view.BuildListComponent(listOptions)
	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent(entityLabel, colours.Purple),
		listComponent,
	})
}
