package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/util"
)

type Frame struct {
	MinWidth  int
	MinHeight int
	Style     lipgloss.Style
	width     int
	height    int
}

func DefaultFrameStyle() Frame {
	s := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true).
		MarginRight(0)
	return Frame{
		Style:     s,
		MinWidth:  util.TermWidth(),
		MinHeight: util.TermHeight(),
	}
}

func (s *Frame) SetSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *Frame) SetWidth(w int) {
	s.width = w
}

func (s Frame) Width() int {
	return util.CalculateWidth(s.MinWidth, s.width)
}

func (s *Frame) SetHeight(h int) {
	s.height = h
}

func (s Frame) Height() int {
	return util.CalculateHeight(s.MinHeight, s.height)
}
