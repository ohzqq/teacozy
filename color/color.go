package color

import "github.com/charmbracelet/lipgloss"

//go:generate gomplate -f _gen/color.tmpl -o gen_color.go -c .=_gen/color.toml

type Color struct {
	Fg     lipgloss.Color `json:"fg" toml:"fg" yaml:"fg"`
	Bg     lipgloss.Color `json:"bg" toml:"bg" yaml:"bg"`
	Black  lipgloss.Color `json:"black" toml:"black" yaml:"black"`
	Blue   lipgloss.Color `json:"blue" toml:"blue" yaml:"blue"`
	Cyan   lipgloss.Color `json:"cyan" toml:"cyan" yaml:"cyan"`
	Green  lipgloss.Color `json:"green" toml:"green" yaml:"green"`
	Grey   lipgloss.Color `json:"grey" toml:"grey" yaml:"grey"`
	Pink   lipgloss.Color `json:"pink" toml:"pink" yaml:"pink"`
	Purple lipgloss.Color `json:"purple" toml:"purple" yaml:"purple"`
	Red    lipgloss.Color `json:"red" toml:"red" yaml:"red"`
	White  lipgloss.Color `json:"white" toml:"white" yaml:"white"`
	Yellow lipgloss.Color `json:"yellow" toml:"yellow" yaml:"yellow"`
}

const (
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
