// Package tablemenuoption defines the available menu options for the list view
// of an entity
package tablemenuoption

// TableMenuOption represents a selectable option in the table management
// interface. It's implemented as a string type for type safety when working
// with menu selections.
type TableMenuOption string

// Menu option constants define the available actions that can be performed on
// tables.
const (
	Search TableMenuOption = "Search" // Filter by keyword
	Create TableMenuOption = "Create" // Create new entity
	Update TableMenuOption = "Update" // Update selected entity
	Delete TableMenuOption = "Delete" // Delete selected entity
	Back   TableMenuOption = "Back"   // Return to previous menu
)
