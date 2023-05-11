package item

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
	Style Style
}

type Props struct {
	Items    teacozy.Items
	Selected bool
	Current  bool
	ReadOnly bool
}

type Prefix struct {
	Fmt   string
	Text  string
	Style lipgloss.Style
}

type Style struct {
	Cursor   Prefix
	Label    Prefix
	Normal   Prefix
	Selected Prefix
	Match    lipgloss.Style
}

func New() *Component {
	return &Component{
		Style: DefaultStyle(),
	}
}

func (c Component) Render(int, int) string {
}

func DefaultStyle() Style {
	return Style{
		Match: lipgloss.NewStyle().Foreground(color.Cyan()),
		Cursor: Prefix{
			Fmt:   currentFmt,
			Text:  currentPrefix,
			Style: lipgloss.NewStyle().Foreground(color.Green()),
		},
		Selected: Prefix{
			Fmt:   selectedFmt,
			Text:  selectedPrefix,
			Style: lipgloss.NewStyle().Foreground(color.Grey()),
		},
		Normal: Prefix{
			Fmt:   unselectedFmt,
			Text:  unselectedPrefix,
			Style: lipgloss.NewStyle().Foreground(color.Fg()),
		},
		Label: Prefix{
			Fmt:   labelFmt,
			Style: lipgloss.NewStyle().Foreground(color.Purple()),
		},
	}
}

const (
	selectedPrefix   = "x"
	selectedFmt      = "[%s]"
	unselectedPrefix = " "
	unselectedFmt    = "[%s]"
	currentPrefix    = "x"
	currentFmt       = "[%s]"
	labelFmt         = "[%s]"
)
