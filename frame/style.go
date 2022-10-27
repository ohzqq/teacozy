package frame

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Border     lipgloss.Border
	Sides      []bool
	Padding    []int
	Margin     []int
	Foreground string
	Background string
}

func DefaultStyle() Styles {
	return Styles{
		Border:     lipgloss.HiddenBorder(),
		Sides:      []bool{true},
		Padding:    []int{0},
		Margin:     []int{0},
		Foreground: "#FFBF00",
		Background: "#262626",
	}
}

func (s Styles) Render() lipgloss.Style {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color(s.Background)).
		Foreground(lipgloss.Color(s.Foreground)).
		Border(s.Border, s.Sides...).
		Padding(s.Padding...).
		Margin(s.Margin...)
	return style
}
