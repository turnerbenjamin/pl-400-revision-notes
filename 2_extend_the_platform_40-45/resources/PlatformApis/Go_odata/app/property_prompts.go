// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"fmt"

	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/view"
)

// propertyPrompt defines a text prompt for collecting entity property values
// from the user. It includes information about the property and functions to
// get and set its value for a specific entity type.
type propertyPrompt[T view.Entity] struct {
	propertyName string          // Name of the property to display to the user
	promptText   string          // Text to display when prompting for input
	isRequired   bool            // Whether the property is required
	getter       func(T) string  // Function to retrieve current property value
	setter       func(T, string) // Function to set the property value
}

// accountPropertyPrompts defines the collection of prompts for Account entity
// properties.
// Each prompt includes display text, validation rules, and getter/setter
// functions.
var accountPropertyPrompts = []propertyPrompt[*model.Account]{
	{
		propertyName: "Name",
		promptText:   "Enter account name",
		isRequired:   true,
		getter: func(a *model.Account) string {
			return a.Name
		},
		setter: func(a *model.Account, value string) {
			a.Name = value
		},
	},
	{
		propertyName: "City",
		promptText:   "Enter account city",
		isRequired:   true,
		getter: func(a *model.Account) string {
			return a.City
		},
		setter: func(a *model.Account, value string) {
			a.City = value
		},
	},
}

// contactPropertyPrompts defines the collection of prompts for Contact entity
// properties.
// Each prompt includes display text, validation rules, and getter/setter
// functions.
var contactPropertyPrompts = []propertyPrompt[*model.Contact]{
	{
		propertyName: "First name",
		promptText:   "Enter contact's first name",
		isRequired:   true,
		getter: func(a *model.Contact) string {
			return a.FirstName
		},
		setter: func(a *model.Contact, value string) {
			a.FirstName = value
		},
	},
	{
		propertyName: "Last name",
		promptText:   "Enter contact's last name",
		isRequired:   false,
		getter: func(a *model.Contact) string {
			return a.LastName
		},
		setter: func(a *model.Contact, value string) {
			a.LastName = value
		},
	},
	{
		propertyName: "Email",
		promptText:   "Enter contact's email address",
		isRequired:   false,
		getter: func(a *model.Contact) string {
			return a.Email
		},
		setter: func(a *model.Contact, value string) {
			a.Email = value
		},
	},
}

// screenOutputFunc represents a function that displays a screen to the user
// and returns the output from that screen. It takes a function that creates
// a Screen object and returns the ScreenOutput containing user input or an error.
// This type is used to abstract the UI interaction when collecting property values.
type screenOutputFunc func(func() (view.Screen, error)) (view.ScreenOutput, error)

// getEntityDetails collects property values from the user for a given entity
// type.
// It displays a series of prompts defined by the prompts parameter, and returns
// a new entity with the collected values.
//
// Parameters:
//   - defaultValues: Initial entity values to display in the prompts
//   - title: Title to display on the input screens
//   - prompts: Collection of property prompts defining the fields to collect
//   - getScreenOutput: Function to display screens and collect user input
//
// Returns the updated entity with user-provided values or an error if input
// fails.
func getEntityDetails[T view.Entity](
	defaultValues T,
	title string,
	prompts []propertyPrompt[T],
	getScreenOutput screenOutputFunc) (T, error) {
	// entityDetails will be modified with new user values
	entityDetails := defaultValues
	// defaultValuesCopy preserved to show original values in prompts
	defaultValuesCopy := defaultValues

	for _, p := range prompts {
		promptOutput, err := getScreenOutput(func() (view.Screen, error) {
			return newStringInputScreen(
				title,
				p.promptText,
				p.propertyName,
				p.getter(defaultValuesCopy),
				p.isRequired)
		})
		if err != nil {
			var zeroValue T
			return zeroValue, fmt.Errorf("getting property %s: %w", p.propertyName, err)
		}
		p.setter(entityDetails, promptOutput.UserInput())
	}
	return entityDetails, nil
}
