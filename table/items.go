package table

import (
	"fmt"
	"os/exec"
	"strings"

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

type Items struct {
	Items       []Item
	Selected    map[int]struct{}
	NumSelected int
	Limit       int
	Cur         int
	Lines       int
	args        []string
	cmd         string
}

type Item struct {
	fuzzy.Match
	Style    style.ListItem
	selected bool
	Label    string
	Width    int
	Depth    int
	exec     *exec.Cmd
	*Prefix
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

type Opt func(*Items)

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

func (m Items) Selections() []int {
	var chosen []int
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			chosen = append(chosen, k)
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

func (m Items) Slice() []string {
	var items []string
	for _, item := range m.Items {
		items = append(items, item.String())
	}
	return items
}

func (m *Items) ChoiceMap(choices []map[string]string) {
	m.Items = ChoiceMapToMatch(choices)
}

func (cp Items) Visible(matches ...string) []Item {
	if len(matches) != 0 {
		return ExactMatches(matches[0], cp.Items)
	}
	return cp.Items
}

func (i Items) CurrentItem() *Item {
	return &i.Items[i.Cur]
}

func (m *Items) ToggleSelection(items ...int) {
	if len(items) == 0 {
		items = []int{m.CurrentItem().Index}
	}
	for _, idx := range items {
		if _, ok := m.Selected[idx]; ok {
			delete(m.Selected, idx)
			m.NumSelected--
		} else if m.NumSelected < m.Limit {
			m.Selected[idx] = struct{}{}
			m.NumSelected++
		}
	}
}

func (i Item) Map() map[string]string {
	return map[string]string{i.Label: i.Str}
}

func (i Item) String() string {
	return i.Str
}

func (i Item) Render(w, h int) string {
	text := lipgloss.StyleRunes(
		i.Str,
		i.MatchedIndexes,
		i.Style.Match,
		i.Style.Text,
	)
	s := lipgloss.NewStyle().Render(text)
	return s
}

func (i Item) LineHeight() int {
	return lipgloss.Height(i.Render(util.TermSize()))
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

func ChoiceSlice[E any](choices []E) Opt {
	return func(i *Items) {
		i.ChoiceMap(MapChoices(choices))
	}
}

func ChoiceMap[M ~map[K]V, K comparable, V any](choices []M) Opt {
	return func(i *Items) {
		i.Items = make([]Item, len(choices))
		for idx, option := range choices {
			for label, val := range option {
				text := lipgloss.NewStyle().Render(fmt.Sprint(val))
				item := NewItem(text, idx)
				item.Label = fmt.Sprint(label)
				i.Items[idx] = item
			}
		}
	}
}
