package item

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/sahilm/fuzzy"
	"golang.org/x/exp/maps"
)

const (
	PromptPrefix     = "> "
	CursorPrefix     = "x"
	SelectedPrefix   = "◉ "
	UnselectedPrefix = " "
)

type Choices []Choice
type Choice map[string]string

func (c Choices) String(i int) string {
	return c[i].Value()
}

func (c Choices) Len() int {
	return len(c)
}

func (c Choices) Filter(s string) []Item {
	matches := []Item{}
	m := fuzzy.FindFrom(s, c)
	if len(m) == 0 {
		return ChoiceMapToMatch(c)
	}
	for _, match := range m {
		item := New()
		item.Match = match
		item.Label = maps.Keys(c[match.Index])[0]
		matches = append(matches, item)
	}
	return matches
}

func (c Choices) Set(idx int, val string) {
	c[idx] = c[idx].Set(val)
}

func (c Choice) Key() string {
	return maps.Keys(c)[0]
}

func (c Choice) Value() string {
	return maps.Values(c)[0]
}

func (c Choice) Set(v string) Choice {
	for k, _ := range c {
		c[k] = v
		break
	}
	return c
}

type Item struct {
	fuzzy.Match
	*Prefix

	Label    string
	Style    style.ListItem
	Current  bool
	Selected bool
	exec     *exec.Cmd
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

func New() Item {
	item := Item{
		Style:  DefaultItemStyle(),
		Prefix: DefaultPrefix(),
	}
	return item
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

func (i Item) Render(w, h int) string {
	var s strings.Builder
	pre := "x"

	if i.Label != "" {
		pre = i.Label
	}

	switch {
	case i.Current:
		pre = i.Style.Cursor.Render(pre)
	default:
		if i.Selected {
			pre = i.Style.Selected.Render(pre)
		} else if i.Label == "" {
			pre = strings.Repeat(" ", lipgloss.Width(pre))
		} else {
			pre = i.Style.Label.Render(pre)
		}
	}

	s.WriteString("[")
	s.WriteString(pre)
	s.WriteString("]")

	text := lipgloss.StyleRunes(
		i.Str,
		i.MatchedIndexes,
		i.Style.Match,
		i.Style.Text,
	)
	s.WriteString(lipgloss.NewStyle().Render(text))

	return s.String()
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

func ChoiceMapToMatch(options Choices) []Item {
	matches := make([]Item, len(options))
	for i, option := range options {
		for label, val := range option {
			item := NewItem(val, i)
			item.Label = label
			matches[i] = item
		}
	}
	return matches
}

func ChoicesToMatch(options []string) []Item {
	matches := make([]Item, len(options))
	for i, option := range options {
		matches[i] = NewItem(option, i)
	}
	return matches
}

func ChoiceMap(c []map[string]string) Choices {
	choices := make(Choices, len(c))
	for i, ch := range c {
		choices[i] = ch
	}
	return choices
}

func ChoiceSliceToMap[E any](c []E) Choices {
	choices := make([]Choice, len(c))
	for i, val := range c {
		choices[i] = Choice{"": fmt.Sprint(val)}
	}
	return choices
}
