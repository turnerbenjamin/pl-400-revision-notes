package consoleInput

import (
	"sync"

	"github.com/eiannone/keyboard"
)

type InputReader interface {
	AwaitInput() (rune, keyboard.Key)
	Open()
	Close()
}

type inputReader struct {
	inputWaitGroup sync.WaitGroup
}

func CreateInputReader() InputReader {
	return &inputReader{}
}

func (r *inputReader) Open() {
	keyboard.Open()
}

func (r *inputReader) Close() {
	keyboard.Open()
}

func (r *inputReader) AwaitInput() (rune, keyboard.Key) {
	r.inputWaitGroup.Wait()
	r.inputWaitGroup.Add(1)

	char, key, _ := keyboard.GetKey()

	r.inputWaitGroup.Done()
	return char, key

}
