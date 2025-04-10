package view

import (
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"
)

type Component interface {
	render()
	isInteractive() bool
	handleKeyboardInput(rune, keyboard.Key) *updateResponse
}

type updateResponse struct {
	doContinue bool
	userInput  string
	target     string
}

func (ur *updateResponse) GetUserInput() string {
	return ur.userInput
}

func (ur *updateResponse) GetTarget() string {
	return ur.target
}

type Screen interface {
	Mount()
	Dismount()
	Refresh()
	handleKeyboardInput(rune, keyboard.Key) *updateResponse
}

type screen struct {
	components             *[]Component
	interactiveComponent   Component
	lastDoFullRefreshValue bool
}

func MakeScreen(cs []Component) Screen {
	var interactiveComponent Component = nil
	for _, c := range cs {
		if c.isInteractive() {
			if interactiveComponent != nil {
				log.Fatal("screen cannot contain more than one interactive component (menu | input)")
			}
			interactiveComponent = c
		}
	}

	return &screen{
		components:           &cs,
		interactiveComponent: interactiveComponent,
	}
}

func (s *screen) Mount() {
	fmt.Print("\033[?25l\033[H\033[2J\033[3J")
	s.render()
}

func (s *screen) Dismount() {
	fmt.Print("\033[?25l\033[H\033[2J\033[3J")
}

func (s *screen) Refresh() {
	if s.shouldDoFullRefresh() {
		fmt.Print("\033[2J\033[3J")
	}
	fmt.Print("\033[H")
	s.render()
}

func (s *screen) render() {
	for _, c := range *s.components {
		c.render()
	}
}

func (s *screen) handleKeyboardInput(char rune, key keyboard.Key) *updateResponse {
	if s.interactiveComponent != nil {
		ur := s.interactiveComponent.handleKeyboardInput(char, key)
		return ur
	}
	return &updateResponse{}
}

func (s *screen) shouldDoFullRefresh() bool {

	threshold := 22
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("error getting terminal size:")
	}
	cv := s.lastDoFullRefreshValue
	s.lastDoFullRefreshValue = height <= threshold
	return cv || s.lastDoFullRefreshValue

}
