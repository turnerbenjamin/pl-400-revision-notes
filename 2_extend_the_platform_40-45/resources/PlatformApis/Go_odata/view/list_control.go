// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

// ListControl represents a custom action that can be performed on list items.
// Implementations must provide methods for getting the associated key,
// display label, and action value.
type ListControl interface {
	// GetKey returns the keyboard key that triggers this control
	GetKey() rune
	// GetLabel returns the display text for this control
	GetLabel() string
	// GetValue returns the action identifier for this control
	GetValue() string
}
