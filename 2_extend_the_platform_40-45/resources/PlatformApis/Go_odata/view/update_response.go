// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

// ScreenOutput defines the interface for retrieving user input from screen
// components.
// It provides access to the user's entered text and target selections.
type ScreenOutput interface {
	// UserInput returns the text entered by the user.
	UserInput() string
	// Target returns the selected target or destination.
	Target() string
}

// updateResponse implements the ScreenOutput interface and handles screen
// interaction results.
// This type is package-private and used internally to manage UI component
// state.
type updateResponse struct {
	// Controls whether the application should continue running
	doContinue bool
	// Indicates if the screen requires a complete refresh
	needsFullRefresh bool
	// Stores text entered by the user
	userInput string
	// Stores the selected target or destination
	target string
}

// UserInput returns the text entered by the user.
func (ur *updateResponse) UserInput() string {
	return ur.userInput
}

// Target returns the selected target or destination.
func (ur *updateResponse) Target() string {
	return ur.target
}

// newUpdateResponse creates a new updateResponse with default values.
// By default, doContinue and needsFullRefresh are set to false, and
// userInput and target are empty strings.
func newUpdateResponse() *updateResponse {
	return &updateResponse{
		doContinue:       false,
		needsFullRefresh: false,
		userInput:        "",
		target:           "",
	}
}

// setContinue updates the continue state and returns the updated response.
// This enables method chaining for fluent configuration.
func (ur *updateResponse) setContinue(cont bool) *updateResponse {
	ur.doContinue = cont
	return ur
}

// setFullRefresh marks the response as requiring a full screen refresh and
// returns the updated response for method chaining.
func (ur *updateResponse) setFullRefresh() *updateResponse {
	ur.needsFullRefresh = true
	return ur
}

// setUserInput updates the user input text and returns the updated response.
// This enables method chaining for fluent configuration.
func (ur *updateResponse) setUserInput(userInput string) *updateResponse {
	ur.userInput = userInput
	return ur
}

// setTarget updates the target selection and returns the updated response.
// This enables method chaining for fluent configuration.
func (ur *updateResponse) setTarget(target string) *updateResponse {
	ur.target = target
	return ur
}
