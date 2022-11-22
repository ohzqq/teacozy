package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type ItemStyle struct {
	NormalItem   lipgloss.Style
	CurrentItem  lipgloss.Style
	SelectedItem lipgloss.Style
	SubItem      lipgloss.Style
}

func ItemStyles() ItemStyle {
	var s ItemStyle
	s.NormalItem = lipgloss.NewStyle().Foreground(Color.DefaultFg)
	s.CurrentItem = lipgloss.NewStyle().Foreground(Color.Green).Reverse(true)
	s.SelectedItem = lipgloss.NewStyle().Foreground(Color.Grey)
	s.SubItem = lipgloss.NewStyle().Foreground(Color.Purple)
	return s
}

func NewDefaultItemStyles() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles()

	s.NormalTitle = lipgloss.NewStyle().
		Foreground(Color.DefaultFg).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(Color.Grey)

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(Color.Purple).
		Foreground(Color.Green).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(Color.Pink)

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(Color.White).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(Color.Grey)

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}
