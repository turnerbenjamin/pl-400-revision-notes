package view

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type menu struct {
	options  *[]string
	selected int
}

func BuildMenuComponent(options []string) Component {

	if len(options) < 2 {
		log.Fatal("menu must contain at least 2 options")
	}

	return &menu{
		options:  &options,
		selected: 0,
	}
}

func (m *menu) render() {
	for i, o := range *m.options {
		b := getBullet(i == m.selected)
		fmt.Printf("%s %s\n", b, o)
	}
}

func (m *menu) isInteractive() bool {
	return true
}

func (m *menu) handleKeyboardInput(c rune, k keyboard.Key) *updateResponse {

	ur := updateResponse{
		doContinue: true,
	}

	if k == keyboard.KeyArrowUp && m.selected > 0 {
		m.selected--
	}
	if k == keyboard.KeyArrowDown && m.selected < len(*m.options)-1 {
		m.selected++
	}
	if k == keyboard.KeyEnter {
		ur.doContinue = false
		ur.userInput = (*&*m.options)[m.selected]
	}
	return &ur
}

func getBullet(isSelected bool) string {
	if isSelected {
		return fmt.Sprintf("%s->%s", colours.ORANGE, colours.RESET)
	}
	return "  "
}
