package list

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/sahilm/fuzzy"
)

type Item struct {
	fuzzy.Match
	Style            style.ListItem
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	isCur            bool
	selected         bool
	prefix           *Prefix
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
	Style      style.ItemPrefix
}

func NewItem(t string, idx int) Item {
	return Item{
		Match: fuzzy.Match{
			Str:   t,
			Index: idx,
		},
		Style:  DefaultItemStyle(),
		prefix: DefaultPrefix(),
	}
}

func (i Item) Prefix() string {
	var s strings.Builder
	if i.isCur {
		s.WriteString(i.prefix.Style.Cursor.Render(i.prefix.Cursor))
	} else {
		if i.selected {
			s.WriteString(i.prefix.Style.Selected.Render(i.prefix.Selected))
		} else if !i.isCur {
			s.WriteString(i.prefix.Style.Unselected.Render(i.prefix.Unselected))
		}
	}
	return s.String()
}

func (i *Item) IsCur() {
	i.isCur = true
}

func DefaultPrefix() *Prefix {
	return &Prefix{
		Cursor:     "> ",
		Selected:   "◉ ",
		Unselected: "○ ",
		Style: style.ItemPrefix{
			Cursor:     style.Cursor,
			Selected:   style.Selected,
			Unselected: style.Unselected,
		},
	}
}

func DefaultItemStyle() style.ListItem {
	var s style.ListItem
	s.Cursor = style.Cursor
	s.SelectedPrefix = style.Selected
	s.UnselectedPrefix = style.Unselected
	s.Text = style.Foreground
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())

	s.Prefix = style.ItemPrefix{
		Cursor:     style.Cursor,
		Selected:   style.Selected,
		Unselected: style.Unselected,
	}
	return s
}

func ChoicesToMatch(options []string) []Item {
	matches := make([]Item, len(options))
	for i, option := range options {
		matches[i] = NewItem(option, i)
	}
	return matches
}

func exactMatches(search string, choices []Item) []Item {
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
