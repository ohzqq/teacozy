package item

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type Item struct {
	fuzzy.Match

	Label    string
	Current  bool
	Selected bool
	Style    Style
	exec     *exec.Cmd
}

func New() Item {
	return Item{
		Style: DefaultStyle(),
	}
}

func NewItem(idx int, t string) Item {
	item := New()
	item.Match = fuzzy.Match{
		Str:   t,
		Index: idx,
	}
	return item
}

func (i *Item) Exec(cmd *exec.Cmd) {
	i.exec = cmd
}

func (i Item) Render(w, h int) string {
	var s strings.Builder
	pre := Cursor

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
		i.Style.Unselected,
	)
	s.WriteString(lipgloss.NewStyle().Render(text))

	return s.String()
}

func ChoicesToItems(options Choices) []Item {
	matches := make([]Item, len(options))
	for i, option := range options {
		for label, val := range option {
			item := NewItem(i, val)
			item.Label = label
			matches[i] = item
		}
	}
	return matches
}
