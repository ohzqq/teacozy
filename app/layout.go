package app

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Top = iota + 1
	Bottom
	Right
	Left
)

const (
	Single = iota + 1
	Half
	Third
	Quarter
)

type Layout struct {
	split    int
	mainPos  int
	sections int
	width    int
	height   int
}

func NewLayout() *Layout {
	l := &Layout{
		split:    Vertical,
		sections: Single,
	}
	switch l.split {
	case Vertical:
		l.mainPos = Top
	case Horizontal:
		l.mainPos = Right
	}
	return l
}

func (l *Layout) Vertical() *Layout {
	l.split = Vertical
	return l
}

func (l *Layout) Horizontal() *Layout {
	l.split = Horizontal
	return l
}

func (l *Layout) SetSize(w, h int) *Layout {
	l.width = w
	l.height = h
	return l
}

func (l Layout) Join(li, page string) string {
	if l.split == Horizontal {
		if l.mainPos == Right {
			return lipgloss.JoinHorizontal(lipgloss.Center, page, li)
		}
		return lipgloss.JoinHorizontal(lipgloss.Center, li, page)
	}

	if l.mainPos == Bottom {
		return lipgloss.JoinVertical(lipgloss.Left, page, li)
	}
	return lipgloss.JoinVertical(lipgloss.Left, li, page)
}

func (l *Layout) Top() *Layout {
	l.mainPos = Top
	return l
}

func (l *Layout) Bottom() *Layout {
	l.mainPos = Bottom
	return l
}
func (l *Layout) Left() *Layout {
	l.mainPos = Left
	return l
}
func (l *Layout) Right() *Layout {
	l.mainPos = Right
	return l
}

func (l *Layout) main() (int, int) {
	w, h := l.width, l.height
	switch l.split {
	case Vertical:
		if l.sections > 1 {
			h = h / l.sections * (l.sections - 1)
		}
	case Horizontal:
		if l.sections > 1 {
			w = w / l.sections * (l.sections - 1)
		}
	}
	if h == 0 {
		return w, h
	}
	return w, h - 2
}

func (l *Layout) sub() (int, int) {
	w, h := 0, 0
	switch l.split {
	case Vertical:
		if l.sections > 1 {
			h = l.height / l.sections
			w = l.width
		}
	case Horizontal:
		if l.sections > 1 {
			w = l.width / l.sections
			h = l.height
		}
	}
	if h == 0 {
		return w, h
	}
	return w, h - 2
}

func (l *Layout) Position(p int) *Layout {
	l.mainPos = p
	return l
}

func (l *Layout) Split(s int) *Layout {
	l.split = s
	switch s {
	case Vertical:
		if l.mainPos > 2 {
			l.mainPos = Top
		}
	case Horizontal:
		if l.mainPos < 3 {
			l.mainPos = Right
		}
	}
	return l
}

func (l *Layout) Single() *Layout {
	return l.Sections(Single)
}

func (l *Layout) Half() *Layout {
	return l.Sections(Half)
}

func (l *Layout) Third() *Layout {
	return l.Sections(Third)
}

func (l *Layout) Quarter() *Layout {
	return l.Sections(Quarter)
}

func (l *Layout) Sections(s int) *Layout {
	l.sections = s
	return l
}

func (l Layout) Width() int {
	return l.width
}

func (l Layout) Height() int {
	return l.height
}
