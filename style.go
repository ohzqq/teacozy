package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type TUIStyle struct {
	Color  Color
	List   list.Styles
	Item   ItemStyle
	Widget WidgetStyle
	Frame  Frame
}

type Frame struct {
	MaxWidth  int
	MaxHeight int
	width     int
	height    int
	Frame     lipgloss.Style
}

type Color struct {
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

func DefaultColors() Color {
	return Color{
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
}

type WidgetStyle struct {
	MaxWidth  int
	MaxHeight int
	width     int
	height    int
}

func (s WidgetStyle) Width() int {
	if s.MaxWidth == 0 {
		s.MaxWidth = TermWidth()
	}

	w := s.MaxWidth

	if s.width != 0 && s.width < s.MaxWidth {
		w = s.width
	}

	return w
}

func (s WidgetStyle) Height() int {
	if s.MaxHeight == 0 {
		s.MaxHeight = TermHeight()
	}

	h := s.MaxHeight

	if s.height != 0 && s.height < s.MaxHeight {
		h = s.height
	}

	return h
}

const (
	Bullet   = "•"
	Ellipsis = "…"
)

type ItemStyle struct {
	NormalItem   lipgloss.Style
	CurrentItem  lipgloss.Style
	SelectedItem lipgloss.Style
	SubItem      lipgloss.Style
}

func ItemStyles() (s ItemStyle) {
	s.NormalItem = lipgloss.NewStyle().Foreground(DefaultColors().DefaultFg)
	s.CurrentItem = lipgloss.NewStyle().Foreground(DefaultColors().Green).Reverse(true)
	s.SelectedItem = lipgloss.NewStyle().Foreground(DefaultColors().Grey)
	s.SubItem = lipgloss.NewStyle().Foreground(DefaultColors().Purple)
	return s
}

func FrameStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true).
		MarginRight(0)
	return s
}

func ListStyles() (s list.Styles) {
	verySubduedColor := DefaultColors().Grey
	subduedColor := DefaultColors().White

	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 0, 0)

	s.Title = lipgloss.NewStyle().
		Background(DefaultColors().Purple).
		Foreground(DefaultColors().Black).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(DefaultColors().Cyan)

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(DefaultColors().Pink)

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(DefaultColors().Yellow)

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(DefaultColors().Blue).
		Padding(0, 0, 1, 2)

	s.StatusEmpty = lipgloss.NewStyle().Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(DefaultColors().Purple)

	s.StatusBarFilterCount = lipgloss.NewStyle().Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(DefaultColors().Grey)

	s.ArabicPagination = lipgloss.NewStyle().Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(DefaultColors().Pink).
		SetString(Bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(Bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + Bullet + " ")

	return s
}
