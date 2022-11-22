package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/keybind"
)

func itemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		//var title string
		var curItem *Field
		if i, ok := m.SelectedItem().(*Field); ok {
			curItem = i
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keybind.EditField):
				return EditFormItemCmd(curItem)
			}
		}
		return nil
	}
	d.SetSpacing(0)
	//d.Styles = itemStyles()
	//d.ShowDescription = false

	return d
}

func itemStyles() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles()
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AFFFAF")).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.Color("#AFFFAF"))

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.Color("#AFFFAF")).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(lipgloss.Color("#AFFFAF"))

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AFFFAF")).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.Color("#AFFFAF"))

	s.FilterMatch = lipgloss.NewStyle().Underline(true)
	return s
}
