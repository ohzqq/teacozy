package color

import "github.com/charmbracelet/lipgloss"

func SetColors(c map[string]string) {
	colors = c
}

func Set(color, val string) {
	colors[color] = val
}

var colors = map[string]string{
	"foreground": "#FFBF00",
	"background": "#262626",
	"black":      "#262626",
	"blue":       "#5FAFFF",
	"cyan":       "#AFFFFF",
	"green":      "#AFFFAF",
	"grey":       "#626262",
	"pink":       "#FFAFFF",
	"purple":     "#AFAFFF",
	"red":        "#FF875F",
	"white":      "#EEEEEE",
	"Yellow":     "#FFFFAF",
}

var (
	Foreground = lipgloss.Color(colors["foreground"])
	Background = lipgloss.Color(colors["background"])
	Black      = lipgloss.Color(colors["black"])
	Blue       = lipgloss.Color(colors["blue"])
	Cyan       = lipgloss.Color(colors["cyan"])
	Green      = lipgloss.Color(colors["green"])
	Grey       = lipgloss.Color(colors["grey"])
	Pink       = lipgloss.Color(colors["pink"])
	Purple     = lipgloss.Color(colors["purple"])
	Red        = lipgloss.Color(colors["red"])
	White      = lipgloss.Color(colors["white"])
	Yellow     = lipgloss.Color(colors["yellow"])
)
