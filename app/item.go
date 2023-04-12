package app

import (
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"github.com/sahilm/fuzzy"
)

const (
	PromptPrefix     = "> "
	CursorPrefix     = "x"
	SelectedPrefix   = "◉ "
	UnselectedPrefix = " "
)

type Item struct {
	reactea.BasicComponent // It implements AfterUpdate() for us, so we don't have to care!
	reactea.BasicPropfulComponent[ItemProps]

	fuzzy.Match
	*Prefix

	Label string
	exec  *exec.Cmd
	Style style.ListItem
}

type ItemProps struct {
	fuzzy.Match

	Current  bool
	Selected bool
	Label    string
	exec     *exec.Cmd
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

func NewItem(t string, idx int) Item {
	item := Item{
		Match: fuzzy.Match{
			Str:   t,
			Index: idx,
		},
		Style:  DefaultItemStyle(),
		Prefix: DefaultPrefix(),
	}
	return item
}

func (i *Item) Init(props ItemProps) tea.Cmd {
	i.UpdateProps(props)
	return nil
}

func DefaultPrefix() *Prefix {
	return &Prefix{
		Cursor:     CursorPrefix,
		Selected:   SelectedPrefix,
		Unselected: UnselectedPrefix,
	}
}

func (i *Item) Exec(cmd *exec.Cmd) {
	i.exec = cmd
}

func (i *Item) SetValue(val string) {
	i.Str = val
}

func (i Item) Map() map[string]string {
	return map[string]string{i.Label: i.Str}
}

func (i Item) String() string {
	return i.Str
}

func (i Item) Render(w, h int) string {
	var s strings.Builder
	pre := "x"

	if i.Props().Label != "" {
		pre = i.Props().Label
	}

	switch {
	case i.Props().Current:
		pre = i.Style.Cursor.Render(pre)
	default:
		if i.Props().Selected {
			pre = i.Style.Selected.Render(pre)
		} else if i.Props().Label == "" {
			pre = strings.Repeat(" ", lipgloss.Width(pre))
		} else {
			pre = i.Style.Label.Render(pre)
		}
	}

	s.WriteString("[")
	s.WriteString(pre)
	s.WriteString("]")

	text := lipgloss.StyleRunes(
		i.Props().Str,
		i.Props().MatchedIndexes,
		i.Style.Match,
		i.Style.Text,
	)
	s.WriteString(lipgloss.NewStyle().Render(text))

	return s.String()
}

func (i Item) LineHeight() int {
	//return 1
	return lipgloss.Height(i.Render(util.TermSize()))
}

func DefaultItemStyle() style.ListItem {
	var s style.ListItem
	s.Cursor = style.Cursor
	s.Selected = style.Selected
	s.Unselected = style.Unselected
	s.Text = style.Foreground
	s.Label = style.Label
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())

	return s
}
