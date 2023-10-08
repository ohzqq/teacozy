package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	bullet   = "•"
	ellipsis = "…"
)

// Styles contains style definitions for this list component. By default, these
// values are generated by DefaultStyles.
type Styles struct {
	list.Styles
}

// DefaultStyles returns a set of default style definitions for this list
// component.
func DefaultStyles() (s Styles) {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#EEEEEE", Dark: "#444444"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#BCBCBC", Dark: "#626262"}

	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 0, 2)

	s.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AFFFFF")).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#8E8E8E", Dark: "#747373"})

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#FFFFAF"})

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF87FF", Dark: "#FF87FF"})

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#BCBCBC", Dark: "#626262"}).
		Padding(0, 0, 0, 2)

	s.StatusEmpty = lipgloss.NewStyle().Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})

	s.StatusBarFilterCount = lipgloss.NewStyle().Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#909090", Dark: "#626262"})

	s.ArabicPagination = lipgloss.NewStyle().Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().Padding(0, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#5FAFFF", Dark: "#5FAFFF"}).
		SetString(bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + bullet + " ")

	return s
}

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	list.DefaultItemStyles
	Prefix  lipgloss.Style
	Toggled lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 1)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#AFFF87", Dark: "#AFFFAF"}).
		Padding(0, 0, 0, 1)

	s.Prefix = s.SelectedTitle.Copy().Padding(0)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#AF87FF", Dark: "#AFAFFF"})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 1)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	s.Toggled = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#646464", Dark: "#646464"})

	return s
}
