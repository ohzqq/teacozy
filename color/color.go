package color

import "github.com/charmbracelet/lipgloss"

var (
	Bg       lipgloss.Color = "#262626"
	Black    lipgloss.Color = "#262626"
	Blue     lipgloss.Color = "#5FAFFF"
	Cyan     lipgloss.Color = "#AFFFFF"
	Fg       lipgloss.Color = "#FFBF00"
	Green    lipgloss.Color = "#AFFFAF"
	Grey     lipgloss.Color = "#626262"
	DarkGrey lipgloss.Color = "#444444"
	Pink     lipgloss.Color = "#FFAFFF"
	Purple   lipgloss.Color = "#AF87FF"
	Red      lipgloss.Color = "#FF5F5F"
	White    lipgloss.Color = "#FFFFFF"
	Yellow   lipgloss.Color = "#FFFF87"
)

func SetBg(c string) {
	Bg = lipgloss.Color(c)
}

func SetBlack(c string) {
	Black = lipgloss.Color(c)
}

func SetBlue(c string) {
	Blue = lipgloss.Color(c)
}

func SetCyan(c string) {
	Cyan = lipgloss.Color(c)
}

func SetFg(c string) {
	Fg = lipgloss.Color(c)
}

func SetGreen(c string) {
	Green = lipgloss.Color(c)
}

func SetGrey(c string) {
	Grey = lipgloss.Color(c)
}

func SetPink(c string) {
	Pink = lipgloss.Color(c)
}

func SetPurple(c string) {
	Purple = lipgloss.Color(c)
}

func SetRed(c string) {
	Red = lipgloss.Color(c)
}

func SetWhite(c string) {
	White = lipgloss.Color(c)
}

func SetYellow(c string) {
	Yellow = lipgloss.Color(c)
}
