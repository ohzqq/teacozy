package style

import "github.com/charmbracelet/lipgloss"

type color struct {
	DefaultFg lipgloss.Color
	DefaultBg lipgloss.Color
	Black     lipgloss.Color
	Blue      lipgloss.Color
	Cyan      lipgloss.Color
	Green     lipgloss.Color
	Grey      lipgloss.Color
	Pink      lipgloss.Color
	Purple    lipgloss.Color
	Red       lipgloss.Color
	White     lipgloss.Color
	Yellow    lipgloss.Color
}

var Color = color{
	DefaultFg: lipgloss.Color("#FFBF00"),
	DefaultBg: lipgloss.Color("#262626"),
	Black:     lipgloss.Color("#262626"),
	Blue:      lipgloss.Color("#5FAFFF"),
	Cyan:      lipgloss.Color("#AFFFFF"),
	Green:     lipgloss.Color("#AFFFAF"),
	Grey:      lipgloss.Color("#626262"),
	Pink:      lipgloss.Color("#FFAFFF"),
	Purple:    lipgloss.Color("#AFAFFF"),
	Red:       lipgloss.Color("#FF875F"),
	White:     lipgloss.Color("#EEEEEE"),
	Yellow:    lipgloss.Color("#FFFFAF"),
}

func DefaultColors() color {
	return Color
}
