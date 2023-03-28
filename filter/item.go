package filter

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/sahilm/fuzzy"
)

type Item struct {
	fuzzy.Match
	Style    style.ListItem
	Label    string
	isCur    bool
	selected bool
	*Prefix
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

const (
	PromptPrefix     = "> "
	CursorPrefix     = "x"
	SelectedPrefix   = "◉ "
	UnselectedPrefix = " "
)

func NewItem(t string, idx int) Item {
	return Item{
		Match: fuzzy.Match{
			Str:   t,
			Index: idx,
		},
		//Label:  "poot",
		Style:  DefaultItemStyle(),
		Prefix: DefaultPrefix(),
	}
}

func DefaultPrefix() *Prefix {
	return &Prefix{
		Cursor:     CursorPrefix,
		Selected:   SelectedPrefix,
		Unselected: UnselectedPrefix,
	}
}

func (match Item) RenderPrefix() string {
	pre := "x"

	if match.Label != "" {
		pre = match.Label
	}

	if match.isCur {
		pre = match.Style.Cursor.Render(pre)
	} else {
		if match.selected {
			pre = match.Style.Selected.Render(pre)
		} else if match.Label == "" {
			pre = strings.Repeat(" ", lipgloss.Width(pre))
		} else {
			pre = match.Style.Label.Render(pre)
		}
	}
	return "[" + pre + "]"
}

func (i *Item) IsCur() {
	i.isCur = true
}

func (i *Item) NotCur() {
	i.isCur = false
}

func (i *Item) Toggle() {
	i.selected = !i.selected
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
