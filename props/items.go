package props

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/util"
)

type Items struct {
	NumSelected int
	Choices     []map[string]string
	Items       []Item
	Selected    map[int]struct{}
	Limit       int
	Height      int
	Width       int
	Snapshot    string
	Current     *Item
}

type Opt func(*Items)

func NewItems(c []map[string]string, opts ...Opt) *Items {
	items := Items{
		Choices:  c,
		Items:    ChoiceMapToMatch(c),
		Selected: make(map[int]struct{}),
	}
	items.Opts(opts...)

	w, h := util.TermSize()
	if items.Height == 0 {
		items.Height = h - 4
	}
	if items.Width == 0 {
		items.Width = w
	}

	items.SetCurrent(0)

	return &items
}

func (i *Items) Opts(opts ...Opt) {
	for _, opt := range opts {
		opt(i)
	}
}

func (i Items) Update() *Items {
	items := NewItems(i.Choices)
	items.Limit = i.Limit
	items.Selected = i.Selected
	items.NumSelected = i.NumSelected
	items.Height = i.Height
	items.Width = i.Width
	items.Current = i.Current
	return items
}

func (m Items) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			c := map[string]string{
				m.Items[k].Label: m.Items[k].Str,
			}
			chosen = append(chosen, c)
		}
	}
	return chosen
}

func (m Items) Map() []map[string]string {
	var items []map[string]string
	for _, item := range m.Items {
		items = append(items, item.Map())
	}
	return items
}

func (i *Items) SetCurrent(idx int) {
	i.Current = &i.Items[idx]
}

func (cp Items) Visible(matches ...string) []Item {
	if len(matches) != 0 {
		return ExactMatches(matches[0], cp.Items)
	}
	return cp.Items
}

func ItemSlice(i []string) *Items {
	items := NewItems(MapChoices(i))
	return items
}

func (m *Items) ToggleSelection(idx int) {
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.NumSelected--
	} else if m.NumSelected < m.Limit {
		m.Selected[idx] = struct{}{}
		m.NumSelected++
	}
}

func (m *Items) ChoiceMap(choices []map[string]string) {
	m.Choices = choices
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

func Limit(l int) Opt {
	return func(i *Items) {
		i.Limit = l
	}
}

func Height(h int) Opt {
	return func(i *Items) {
		i.Height = h
	}
}
