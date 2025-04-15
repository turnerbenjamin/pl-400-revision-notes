// Package mainmenuoption defines the available top-level menu options for the
// application.
package mainmenuoption

// MainMenuOption represents a selectable option in the main application menu.
// It's implemented as a string type for type safety when working with menu
// selections.
type MainMenuOption string

// Menu option constants define the available choices in the main menu.
const (
	Accounts MainMenuOption = "Accounts" // Account entity list
	Contacts MainMenuOption = "Contacts" // Contact entity list
	Exit     MainMenuOption = "Exit"     // Quit application
	Invalid  MainMenuOption = "Invalid"  // Invalid selection
)
