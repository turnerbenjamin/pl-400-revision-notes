// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"errors"
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

// minimumOptionsLength defines the minimum number of options required for a
// functional menu. A menu with fewer options wouldn't provide meaningful choice.
const minimumOptionsLength = 2

// ErrOptionCountBelowMin is returned when attempting to create a menu
// with fewer than the minimum required number of options (minimumOptionsLength).
// This ensures that menus always provide users with meaningful choices.
var ErrOptionCountBelowMin = errors.New(
	fmt.Sprintf("menu must contain at least %d options", minimumOptionsLength),
)

// menu represents a selectable list of options in a terminal UI.
// It allows navigation between options using arrow keys and selection with
// Enter. The currently selected item is visually highlighted with an arrow
// indicator.
type menu struct {
	options  []string
	selected int
}

// NewMenuComponent creates a new menu component with the provided options.
// The menu requires at least two options to be functional. It returns an
// implementation of the InteractiveComponent interface defined in screen.go.
func NewMenuComponent(options []string) (InteractiveComponent, error) {
	optionsCopy := make([]string, len(options))
	copy(optionsCopy, options)

	if len(optionsCopy) < minimumOptionsLength {
		return nil, ErrOptionCountBelowMin
	}

	return &menu{
		options:  optionsCopy,
		selected: 0,
	}, nil
}

// render displays all menu options on the terminal. The currently selected
// option is marked with an arrow indicator.
func (m *menu) render() {
	for i, option := range m.options {
		indicator := menuIndicator(i == m.selected)
		fmt.Printf("%s %s\n", indicator, option)
	}
}

// handleKeyboardInput processes keyboard input events for menu navigation. It
// routes different keys to their appropriate handlers and returns an
// updateResponse struct (defined in screen.go) with processing instructions.
func (m *menu) handleKeyboardInput(character rune, key keyboard.Key) (*updateResponse, error) {
	switch key {
	case keyboard.KeyArrowUp:
		return m.handleArrowUpPressed(), nil
	case keyboard.KeyArrowDown:
		return m.handleArrowDownPressed(), nil
	case keyboard.KeyEnter:
		response, err := m.handleEnterPressed()
		return response, err
	default:
		return newUpdateResponse().setContinue(true), nil
	}
}

// handleArrowUpPressed moves the selection cursor up one menu item if not
// already at the top item. Returns an updateResponse indicating that
// input processing should continue.
func (m *menu) handleArrowUpPressed() *updateResponse {
	if m.selected > 0 {
		m.selected--
	}
	return newUpdateResponse().setContinue(true)
}

// handleArrowDownPressed moves the selection cursor down one menu item
// if not already at the bottom item. Returns an updateResponse indicating that
// input processing should continue.
func (m *menu) handleArrowDownPressed() *updateResponse {
	if m.selected < len(m.options)-1 {
		m.selected++
	}
	return newUpdateResponse().setContinue(true)
}

// handleEnterPressed confirms the current menu selection. Returns an
// updateResponse with doContinue set to false and userInput containing the
// selected option's text.
func (m *menu) handleEnterPressed() (*updateResponse, error) {
	if m.selected < 0 || m.selected >= len(m.options) {
		msg := fmt.Sprintf("Selection out of range (%d)", m.selected)
		return nil, errors.New(msg)
	}

	userInput := m.options[m.selected]
	return newUpdateResponse().setContinue(false).setUserInput(userInput), nil
}

// menuIndicator returns the visual indicator for a menu item based
// on its selection state. Selected items are marked with an orange arrow,
// while unselected items have blank space.
func menuIndicator(isSelected bool) string {
	if isSelected {
		return colours.ApplyColour("->", colours.Orange)
	}
	return "  "
}
