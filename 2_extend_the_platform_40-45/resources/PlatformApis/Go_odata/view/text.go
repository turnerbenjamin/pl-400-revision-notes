package view

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type text struct {
	content string
}

func (t *text) render() {
	fmt.Printf("%s\n\n", t.content)
}

func (t *text) isInteractive() bool {
	return false
}
func (t *text) handleKeyboardInput(c rune, k keyboard.Key) *updateResponse {
	return nil
}

func BuildTextComponent(str string) Component {
	return &text{
		content: str,
	}
}
