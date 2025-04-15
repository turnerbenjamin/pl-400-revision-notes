// Package confirmoption provides constants for confirmation options used in
// interactive prompts within the application.
package confirmoption

// ConfirmOption represents a selectable option for confirmation prompts.
// It's implemented as a string type for type safety when working with user
// confirmation choices.
type ConfirmOption string

// Constants representing confirmation choices.
const (
	Yes ConfirmOption = "Yes"
	No  ConfirmOption = "No"
)
