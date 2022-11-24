package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

type ItemStyle struct {
	NormalItem   lipgloss.Style
	CurrentItem  lipgloss.Style
	SelectedItem lipgloss.Style
	SubItem      lipgloss.Style
}

func ItemStyles() ItemStyle {
	var s ItemStyle
	s.NormalItem = lipgloss.NewStyle().Foreground(color.Foreground)
	s.CurrentItem = lipgloss.NewStyle().Foreground(color.Green).Reverse(true)
	s.SelectedItem = lipgloss.NewStyle().Foreground(color.Grey)
	s.SubItem = lipgloss.NewStyle().Foreground(color.Purple)
	return s
}

func NewDefaultItemStyles() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles()

	s.NormalTitle = lipgloss.NewStyle().
		Foreground(color.Foreground).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(color.Grey)

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(color.Purple).
		Foreground(color.Green).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(color.Pink)

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(color.White).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(color.Grey)

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}
