package app

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
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

func ChoiceMapToMatch(options []map[string]string) []Item {
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

func ExactMatches(search string, choices []Item) []Item {
	matches := []Item{}
	for _, choice := range choices {
		search = strings.ToLower(search)
		matchedString := strings.ToLower(choice.Str)

		index := strings.Index(matchedString, search)
		if index >= 0 {
			for s := range search {
				choice.MatchedIndexes = append(choice.MatchedIndexes, index+s)
			}
			matches = append(matches, choice)
		}
	}

	return matches
}

func MapChoices[E any](c []E) []map[string]string {
	choices := make([]map[string]string, len(c))
	for i, val := range c {
		choices[i] = map[string]string{"": fmt.Sprint(val)}
	}
	return choices
}
