package props

import (
	"os/exec"

	"github.com/charmbracelet/lipgloss"
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
	fuzzy.Match
	Style    style.ListItem
	selected bool
	Label    string
	Width    int
	Depth    int
	exec     *exec.Cmd
	*Prefix
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
	text := lipgloss.StyleRunes(
		i.Str,
		i.MatchedIndexes,
		i.Style.Match,
		i.Style.Text,
	)
	s := lipgloss.NewStyle().Render(text)
	return s
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
