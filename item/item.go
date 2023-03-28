package item

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
	Matches     []Item
	Selected    map[int]struct{}
	Limit       int
	numSelected int
	Cursor      int
}

type Item struct {
	fuzzy.Match
	Style     style.ListItem
	Label     string
	isCur     bool
	IsCurrent bool
	selected  bool
	*Prefix
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

func New(c []string) Items {
	items := Items{
		Items:    ChoicesToMatch(c),
		Selected: make(map[int]struct{}),
	}
	items.Matches = items.Items
	return items
}

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

func (m Items) Chosen() []int {
	var chosen []int
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			chosen = append(chosen, k)
		}
	} else if len(m.Matches) > m.Cursor && m.Cursor >= 0 {
		chosen = append(chosen, m.Cursor)
	}
	return chosen
}

func (m *Items) ToggleSelection() {
	idx := m.Matches[m.Cursor].Index
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.Limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
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

func (i *Item) IsCur() {
	i.isCur = true
	i.IsCurrent = true
}

func (i *Item) Cur(cur bool) {
	i.isCur = cur
	i.IsCurrent = true
}

func (i *Item) NotCur() {
	i.isCur = false
	i.IsCurrent = false
}

func (i *Item) Toggle() {
	i.selected = !i.selected
}

func (i *Item) Select() {
	i.selected = true
}

func (i *Item) Deselect() {
	i.selected = false
}

func (i Item) Selected() bool {
	return i.selected
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

func RenderItems(cursor int, items []Item) string {
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
			if match.selected {
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

		//text := lipgloss.StyleRunes(
		//  match.Str,
		//  match.MatchedIndexes,
		//  match.Style.Match,
		//  match.Style.Text,
		//)

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
