package color

import "github.com/charmbracelet/lipgloss"

//go:generate gomplate -f _gen/color.tmpl -o gen_color.go -c .=_gen/color.toml

type Color struct {
	Fg     lipgloss.Color
	Bg     lipgloss.Color
	Black  lipgloss.Color
	Blue   lipgloss.Color
	Cyan   lipgloss.Color
	Green  lipgloss.Color
	Grey   lipgloss.Color
	Pink   lipgloss.Color
	Purple lipgloss.Color
	Red    lipgloss.Color
	White  lipgloss.Color
	Yellow lipgloss.Color
}
