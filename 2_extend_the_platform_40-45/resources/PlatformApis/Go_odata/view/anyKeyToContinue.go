package view

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type anyKeyToContinue struct {
}

func (t *anyKeyToContinue) render() {
	fmt.Print("\n\nPress any key to continue")
}

func (t *anyKeyToContinue) isInteractive() bool {
	return true
}
func (t *anyKeyToContinue) handleKeyboardInput(c rune, k keyboard.Key) *updateResponse {
	return &updateResponse{
		doContinue: false,
		userInput:  "",
	}
}

func BuildAnyKeyToContinueComponent() Component {
	return &anyKeyToContinue{}
}
