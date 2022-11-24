package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

const (
	Bullet   = "•"
	Ellipsis = "…"
)

type TUIStyle struct {
	Color  Color
	List   list.Styles
	Item   ItemStyle
	Widget Frame
	Frame  Frame
}

func DefaultTuiStyle() TUIStyle {
	tui := TUIStyle{
		Color:  DefaultColors(),
		Item:   ItemStyles(),
		List:   ListStyles(),
		Widget: DefaultWidgetStyle(),
		Frame:  DefaultFrameStyle(),
	}

	return tui
}

type Frame struct {
	MinWidth  int
	MinHeight int
	Style     lipgloss.Style
	width     int
	height    int
}

func DefaultFrameStyle() Frame {
	return Frame{
		Style:     FrameStyle(),
		MinWidth:  TermWidth(),
		MinHeight: TermHeight(),
	}
}

func (s *Frame) SetSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *Frame) SetWidth(w int) {
	s.width = w
}

func (s Frame) Width() int {
	return CalculateWidth(s.MinWidth, s.width)
}

func (s *Frame) SetHeight(h int) {
	s.height = h
}

func (s Frame) Height() int {
	return CalculateHeight(s.MinHeight, s.height)
}

func CalculateWidth(min, width int) int {
	max := TermWidth()
	w := min

	if width != 0 {
		switch {
		case width < min:
			w = width
		case width > max:
			w = max
		}
	}

	return w
}

func CalculateHeight(min, height int) int {
	max := TermHeight()
	h := min

	if height != 0 {
		switch {
		case height < min:
			h = height
		case height > max:
			h = max
		}
	}

	return h
}

type Color struct {
	Foreground lipgloss.Color
	Background lipgloss.Color
	Black      lipgloss.Color
	Blue       lipgloss.Color
	Cyan       lipgloss.Color
	Green      lipgloss.Color
	Grey       lipgloss.Color
	Pink       lipgloss.Color
	Purple     lipgloss.Color
	Red        lipgloss.Color
	White      lipgloss.Color
	Yellow     lipgloss.Color
}

func DefaultColors() Color {
	return Color{
		Foreground: color.Foreground,
		Background: color.Background,
		Black:      color.Black,
		Blue:       color.Blue,
		Cyan:       color.Cyan,
		Green:      color.Green,
		Grey:       color.Grey,
		Pink:       color.Pink,
		Purple:     color.Purple,
		Red:        color.Red,
		White:      color.White,
		Yellow:     color.Yellow,
	}
}

type WidgetStyle struct {
	MinWidth  int
	MinHeight int
	width     int
	height    int
}

func DefaultWidgetStyle() Frame {
	return Frame{
		MinWidth:  TermWidth(),
		MinHeight: TermHeight() / 3,
	}
}

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

func FrameStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true).
		MarginRight(0)
	return s
}

func ListStyles() list.Styles {
	verySubduedColor := color.Grey
	subduedColor := color.White

	var s list.Styles

	s.TitleBar = lipgloss.NewStyle().
		Padding(0, 0, 0, 0)

	s.Title = lipgloss.NewStyle().
		Background(color.Purple).
		Foreground(color.Black).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(color.Cyan)

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(color.Pink)

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(color.Yellow)

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().
		Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(DefaultColors().Blue).
		Padding(0, 0, 1, 2)

	s.StatusEmpty = lipgloss.NewStyle().
		Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(color.Purple)

	s.StatusBarFilterCount = lipgloss.NewStyle().
		Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(color.Grey)

	s.ArabicPagination = lipgloss.NewStyle().
		Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().
		PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().
		Padding(1, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(color.Pink).
		SetString(Bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(Bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + Bullet + " ")

	return s
}

var fieldStyle = FieldStyle{
	Key:   lipgloss.NewStyle().Foreground(color.Blue),
	Value: lipgloss.NewStyle().Foreground(color.Foreground),
}

type FieldStyle struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}
