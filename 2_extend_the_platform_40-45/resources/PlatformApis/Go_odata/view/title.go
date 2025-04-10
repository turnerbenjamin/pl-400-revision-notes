package view

import (
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type title struct {
	content string
	colour  colours.Color
}

func (t *title) render() {

	fmt.Printf("%s%s%s\n\n", t.colour, strings.ToUpper(t.content), colours.RESET)
}

func (t *title) isInteractive() bool {
	return false
}
func (t *title) handleKeyboardInput(c rune, k keyboard.Key) *updateResponse {
	return nil
}

func BuildTitleComponent(str string, colour colours.Color) Component {
	return &title{
		content: str,
		colour:  colour,
	}
}
