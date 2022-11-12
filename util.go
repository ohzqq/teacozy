package teacozy

import (
	"os"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func ClearScreen() {
	termenv.ClearScreen()
}

func TermSize() (int, int) {
	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	return w, h
}

func TermWidth() int {
	w, _, _ := term.GetSize(int(os.Stdin.Fd()))
	return w
}

func TermHeight() int {
	_, h, _ := term.GetSize(int(os.Stdin.Fd()))
	return h
}
