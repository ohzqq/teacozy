package props

import (
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
	*Prefix
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

func MapChoices(c []string) []map[string]string {
	choices := make([]map[string]string, len(c))
	for i, val := range c {
		choices[i] = map[string]string{"": val}
	}
	return choices
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

func (i *Item) SetValue(val string) {
	i.Str = val
}

func (i Item) Map() map[string]string {
	return map[string]string{i.Label: i.Str}
}

func (match Item) RenderText() string {
	text := lipgloss.StyleRunes(
		match.Str,
		match.MatchedIndexes,
		match.Style.Match,
		match.Style.Text,
	)
	w := util.TermWidth()
	s := lipgloss.NewStyle().Width(w).Render(text)
	return s
}

func (i Item) LineHeight() int {
	return lipgloss.Height(i.RenderText())
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
