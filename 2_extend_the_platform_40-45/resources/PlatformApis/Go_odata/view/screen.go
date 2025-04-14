// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"errors"
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/constants/ansi"
	"github.com/turnerbenjamin/go_odata/utilities"
)

// minimumScreenHeightForPartialRefresh defines the minimum terminal height
// required to perform partial screen refreshes. If the terminal is smaller
// than this value, a full refresh will be performed instead.
const (
	minimumScreenHeightForPartialRefresh = 22
)

// ErrNoInteractiveComponent is returned when attempting to create a Screen
// without providing any components that implement the InteractiveComponent
// interface.
// A valid Screen must have exactly one interactive component to handle user
// input.
var ErrNoInteractiveComponent = errors.New(
	"screen creation failed: no interactive component found")

// ErrMultipleInteractiveComponents is returned when attempting to create a
// Screen with more than one component implementing the InteractiveComponent
// interface.
// A valid Screen must have exactly one interactive component to maintain clear
// input handling responsibility.
var ErrMultipleInteractiveComponents = errors.New(
	"screen creation failed: multiple interactive components found")

// Component defines the interface for any UI element that can be rendered
// to the terminal.
type Component interface {
	// render outputs the component's visual representation to the terminal.
	render()
}

// InteractiveComponent extends Component to add user input handling
// capabilities.
// A screen must contain exactly one interactive component.
type InteractiveComponent interface {
	Component
	// handleKeyboardInput processes keyboard input and returns update
	// information and any errors that occurred during processing.
	handleKeyboardInput(rune, keyboard.Key) (*updateResponse, error)
}

// Screen represents a complete terminal UI view composed of multiple
// components.
// It manages component layout, rendering, and input handling.
type Screen interface {
	// Mount prepares the screen for display and renders all components.
	Mount()

	// Dismount cleans up the terminal state when the screen is no longer
	// needed.
	Dismount()

	// Refresh updates the screen content, either partially or fully depending
	// on context.
	Refresh()

	// handleKeyboardInput delegates keyboard input to the interactive
	// component.
	handleKeyboardInput(rune, keyboard.Key) (*updateResponse, error)
}

// screen implements the Screen interface.
type screen struct {
	components           []Component
	interactiveComponent InteractiveComponent
	needsFullRefresh     bool
}

// MakeScreen creates a new Screen from the provided components.
// It requires exactly one component to implement the InteractiveComponent
// interface.
// Returns an error if no interactive component or multiple interactive
// components are provided.
func MakeScreen(cs []Component) (Screen, error) {
	s := screen{
		components: make([]Component, 0, len(cs)),
	}

	for _, c := range cs {
		err := s.registerComponent(c)
		if err != nil {
			return nil, err
		}
	}

	if s.interactiveComponent == nil {
		return nil, ErrNoInteractiveComponent
	}

	return &s, nil
}

// registerComponent adds a component to the screen and identifies interactive
// components.
// If the component implements InteractiveComponent, it's set as the screen's
// interactive component. An error is returned if there is more than one
// interactive component
func (s *screen) registerComponent(c Component) error {
	s.components = append(s.components, c)

	if ic, ok := c.(InteractiveComponent); ok {
		if s.interactiveComponent != nil {
			return ErrMultipleInteractiveComponents
		}
		s.interactiveComponent = ic
	}
	return nil
}

// Mount initializes the screen for display by hiding the cursor,
// clearing the terminal, and rendering all components.
func (s *screen) Mount() {
	fmt.Print(ansi.CursorHide + ansi.ClearAll)
	s.render()
}

// Dismount restores the terminal to its normal state by showing the cursor
// and clearing the screen.
func (s *screen) Dismount() {
	fmt.Print(ansi.CursorShow + ansi.ClearAll)
}

// Refresh updates the screen content, choosing between full or partial refresh
// based on terminal size and internal state.
func (s *screen) Refresh() {
	if s.shouldDoFullRefresh() {
		fmt.Print(ansi.ClearAll)
	} else {
		fmt.Print(ansi.CursorHome)
	}
	s.render()
}

// render displays all components of the screen in sequence.
func (s *screen) render() {
	for _, c := range s.components {
		c.render()
	}
}

// handleKeyboardInput processes keyboard input by delegating to the interactive
// component and tracks whether a full screen refresh is needed.
func (s *screen) handleKeyboardInput(char rune, key keyboard.Key) (*updateResponse, error) {
	sr, err := s.interactiveComponent.handleKeyboardInput(char, key)
	if sr != nil && sr.needsFullRefresh {
		s.needsFullRefresh = true
	}
	return sr, err
}

// shouldDoFullRefresh determines whether a complete screen redraw is needed
// based on terminal size and internal state. It also updates the internal
// refresh flag based on current terminal dimensions.
func (s *screen) shouldDoFullRefresh() bool {
	height := utilities.GetConsoleHeight(minimumScreenHeightForPartialRefresh - 1)

	currentValue := s.needsFullRefresh
	s.needsFullRefresh = height < minimumScreenHeightForPartialRefresh
	return currentValue || s.needsFullRefresh
}
