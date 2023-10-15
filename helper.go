package teacozy

import (
	"os"

	"golang.org/x/term"
)

func TermSize() (int, int) {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return w, h
}

func TermHeight() int {
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return h
}

func TermWidth() int {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return w
}
