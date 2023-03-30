package list

import (
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

type Items struct {
	Items       []Item
	Selected    map[int]struct{}
	Limit       int
	numSelected int
	Cursor      int
}

type Item struct {
	fuzzy.Match
	Style style.ListItem
	Label string
	*Prefix
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

func NewItems(c []map[string]string) Items {
	items := Items{
		Items:    ChoiceMapToMatch(c),
		Selected: make(map[int]struct{}),
	}
	return items
}

func NewChoiceMap(c []map[string]string) Items {
	items := Items{
		Items:    ChoiceMapToMatch(c),
		Selected: make(map[int]struct{}),
	}
	return items
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

func (m Items) Chosen() []int {
	var chosen []int
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			chosen = append(chosen, k)
		}
	} else if len(m.Items) > m.Cursor && m.Cursor >= 0 {
		chosen = append(chosen, m.Cursor)
	}
	return chosen
}

func (match Item) RenderText() string {
	text := lipgloss.StyleRunes(
		match.Str,
		match.MatchedIndexes,
		match.Style.Match,
		match.Style.Text,
	)
	return text
}

func (m Items) RenderItems(cursor int, items []Item) string {
	var s strings.Builder
	for i, match := range items {
		pre := "x"

		if match.Label != "" {
			pre = match.Label
		}

		switch {
		case i == cursor:
			pre = match.Style.Cursor.Render(pre)
		default:
			if _, ok := m.Selected[match.Index]; ok {
				pre = match.Style.Selected.Render(pre)
			} else if match.Label == "" {
				pre = strings.Repeat(" ", lipgloss.Width(pre))
			} else {
				pre = match.Style.Label.Render(pre)
			}
		}

		s.WriteString("[")
		s.WriteString(pre)
		s.WriteString("]")

		s.WriteString(match.RenderText())
		s.WriteRune('\n')
	}
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
