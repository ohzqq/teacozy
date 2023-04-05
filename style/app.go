package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

type App struct {
	Confirm lipgloss.Style
	Footer  lipgloss.Style
	Header  lipgloss.Style
}

func DefaultAppStyle() App {
	return App{
		Confirm: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
		Footer:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
		Header:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
	}
}
