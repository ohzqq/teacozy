package item

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
	"github.com/sahilm/fuzzy"
)

type Item struct {
	fuzzy.Match

	Label    string
	Current  bool
	Selected bool
	exec     *exec.Cmd
}

func New() Item {
	item := Item{}
	return item
}

func NewItem(idx int, t string) Item {
	item := Item{
		Match: fuzzy.Match{
			Str:   t,
			Index: idx,
		},
	}
	return item
}

func (i *Item) Exec(cmd *exec.Cmd) {
	i.exec = cmd
}

func (i Item) Render(w, h int) string {
	var s strings.Builder
	pre := "x"

	if i.Label != "" {
		style.Prefix().Cursor().Set(i.Label)
		style.Prefix().Selected().Set(i.Label)
	}

	switch {
	case i.Current:
		pre = style.Prefix().Cursor().Render()
	default:
		if i.Selected {
			pre = style.Prefix().Selected().Render()
		} else if i.Label == "" {
			pre = strings.Repeat(" ", lipgloss.Width(style.Prefix().Cursor().Text))
		} else {
			pre = style.Label.Render(i.Label)
		}
	}

	s.WriteString("[")
	s.WriteString(pre)
	s.WriteString("]")

	text := lipgloss.StyleRunes(
		i.Str,
		i.MatchedIndexes,
		style.Match,
		style.Foreground,
	)
	s.WriteString(lipgloss.NewStyle().Render(text))

	return s.String()
}

func ChoiceMapToItems(options Choices) []Item {
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
