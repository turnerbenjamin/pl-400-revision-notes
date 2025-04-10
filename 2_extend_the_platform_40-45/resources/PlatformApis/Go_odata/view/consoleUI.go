package view

import (
	"github.com/turnerbenjamin/go_odata/view/consoleInput"
)

type ScreenOutputs interface {
	GetUserInput() string
	GetTarget() string
}

type UI interface {
	NavigateTo(Screen) ScreenOutputs
	Exit()
}

type consoleUi struct {
	currentScreen Screen
	inputReader   consoleInput.InputReader
}

func NewConsoleUI() UI {
	ir := consoleInput.CreateInputReader()
	ir.Open()

	ui := consoleUi{
		inputReader: ir,
	}
	return &ui
}

func (ui *consoleUi) NavigateTo(s Screen) ScreenOutputs {
	if ui.currentScreen != nil {
		ui.currentScreen.Dismount()
	}

	ui.currentScreen = s
	ui.currentScreen.Mount()
	return ui.AwaitOutput()
}

func (ui *consoleUi) Exit() {
	if ui.currentScreen != nil {
		ui.currentScreen.Dismount()
	}
	ui.inputReader.Close()
}

func (ui *consoleUi) AwaitOutput() ScreenOutputs {
	for {
		char, key := ui.inputReader.AwaitInput()
		updateResponse := ui.currentScreen.handleKeyboardInput(char, key)

		if !updateResponse.doContinue {
			return updateResponse
		}
		ui.currentScreen.Refresh()
	}
}
