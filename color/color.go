package color

import "github.com/charmbracelet/lipgloss"

var (
	Foreground = lipgloss.Color("#FFBF00")
	Background = lipgloss.Color("#262626")
	Black      = lipgloss.Color("#262626")
	Blue       = lipgloss.Color("#5FAFFF")
	Cyan       = lipgloss.Color("#AFFFFF")
	Green      = lipgloss.Color("#AFFFAF")
	Grey       = lipgloss.Color("#626262")
	Pink       = lipgloss.Color("#FFAFFF")
	Purple     = lipgloss.Color("#AFAFFF")
	Red        = lipgloss.Color("#FF875F")
	White      = lipgloss.Color("#EEEEEE")
	Yellow     = lipgloss.Color("#FFFFAF")
)

func SetColors(c map[string]string) {
	for color, val := range c {
		Set(color, val)
	}
}

func Set(color, val string) {
	switch color {
	case "foreground":
		Foreground = lipgloss.Color(val)
	case "background":
		Background = lipgloss.Color(val)
	case "black":
		Black = lipgloss.Color(val)
	case "blue":
		Blue = lipgloss.Color(val)
	case "cyan":
		Cyan = lipgloss.Color(val)
	case "green":
		Green = lipgloss.Color(val)
	case "grey":
		Grey = lipgloss.Color(val)
	case "pink":
		Pink = lipgloss.Color(val)
	case "purple":
		Purple = lipgloss.Color(val)
	case "red":
		Red = lipgloss.Color(val)
	case "white":
		White = lipgloss.Color(val)
	case "yellow":
		Yellow = lipgloss.Color(val)
	}
}

func Colors() []string {
	return []string{
		"foreground",
		"background",
		"black",
		"blue",
		"cyan",
		"green",
		"grey",
		"pink",
		"purple",
		"red",
		"white",
		"yellow",
	}
}
