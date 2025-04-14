// Package consoleinputreader provides a simple interface for reading keyboard input from the console.
// It wraps the github.com/eiannone/keyboard package to allow for dependency injection and easier testing.
package consoleinputreader

import (
	"github.com/eiannone/keyboard"
)

// InputReader defines the interface for reading keyboard input from the
// console.
// Implementations must support opening a keyboard connection, waiting for user
// input, and properly closing the keyboard connection when finished.
type InputReader interface {
	// AwaitInput blocks until a key is pressed and returns the character,
	// key code, and any error.
	AwaitInput() (rune, keyboard.Key, error)

	// Open initializes the keyboard input connection. This must be called
	// before AwaitInput can be used.
	Open() error

	// Close terminates the keyboard input connection and releases associated
	// resources.
	Close() error
}

// inputReader implements the InputReader interface using the keyboard package.
type inputReader struct{}

// NewInputReader creates and returns a new instance of InputReader.
// Usage:
//
//	reader := consoleinputreader.NewInputReader()
//	err := reader.Open()
//	defer reader.Close()
func NewInputReader() InputReader {
	return &inputReader{}
}

// Open initializes the keyboard connection.
func (r *inputReader) Open() error {
	return keyboard.Open()
}

// Close terminates the keyboard connection and releases resources.
func (r *inputReader) Close() error {
	return keyboard.Close()
}

// AwaitInput blocks until a key is pressed and returns the character,
// the key code, and any error that occurred.
func (r *inputReader) AwaitInput() (rune, keyboard.Key, error) {
	return keyboard.GetKey()
}
