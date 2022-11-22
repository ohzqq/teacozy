package util

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

func CalculateWidth(min, width int) int {
	max := TermWidth()
	w := min

	if width != 0 {
		switch {
		case width < min:
			w = width
		case width > max:
			w = max
		}
	}

	return w
}

func CalculateHeight(min, height int) int {
	max := TermHeight()
	h := min

	if height != 0 {
		switch {
		case height < min:
			h = height
		case height > max:
			h = max
		}
	}

	return h
}
