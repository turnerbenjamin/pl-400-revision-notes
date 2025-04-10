package view

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type stringInput struct {
	propertyName string
	value        string
	isRequired   bool
	errorMessage string
	requiredFlag string
}

func BuildStringInputComponent(propertyName string, isRequired bool) Component {
	si := &stringInput{
		propertyName: propertyName,
		value:        "",
		isRequired:   isRequired,
	}
	if isRequired {
		si.requiredFlag = fmt.Sprintf("(%s*%s)", colours.RED, colours.RESET)
	}
	return si
}

func (m *stringInput) render() {
	fmt.Printf("\n%s%s: %s", m.propertyName, m.requiredFlag, m.value)
	if m.errorMessage != "" {
		fmt.Printf("\n\n%s%s%s", colours.RED, m.errorMessage, colours.RESET)
	}
}

func (m *stringInput) isInteractive() bool {
	return true
}

func (m *stringInput) handleKeyboardInput(c rune, k keyboard.Key) *updateResponse {
	ur := updateResponse{
		doContinue: true,
	}
	m.errorMessage = ""

	switch k {
	case keyboard.KeyEnter:
		if m.isRequired && m.value == "" {
			m.errorMessage = fmt.Sprintf("%s is required", m.propertyName)
		} else {
			ur.doContinue = false
			ur.userInput = m.value
		}
	case keyboard.KeyBackspace:
		fallthrough
	case keyboard.KeyBackspace2:
		if len(m.value) > 0 {
			m.value = m.value[:len(m.value)-1]
		}
	case keyboard.KeySpace:
		m.value += " "
	default:
		m.value += string(c)
	}
	return &ur
}
