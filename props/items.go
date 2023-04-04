package props

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/slices"
)

type Items struct {
	Items       []Item
	Selected    map[int]struct{}
	NumSelected int
	Limit       int
	Height      int
	Width       int
	Snapshot    string
	Title       string
	Cur         int
	Line        int
	footer      string
	TotalLines  int
	SetHeader   func(string)
	SetBody     func(string)
	SetFooter   func(string)
	SetHelp     func(keys.KeyMap)
	args        []string
	cmd         string
}

type Opt func(*Items)

func New(opts ...Opt) *Items {
	items := Items{
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

func (i *Items) ChoiceSlice(c []string) {
	i.Items = ChoiceMapToMatch(MapChoices(c))
}

func (i *Items) Opts(opts ...Opt) {
	for _, opt := range opts {
		opt(i)
	}
}

func (i *Items) NoLimit() *Items {
	i.Limit = len(i.Items)
	return i
}

func (i Items) Update() *Items {
	items := &Items{}
	items.Items = i.Items
	items.Limit = i.Limit
	items.Selected = i.Selected
	items.NumSelected = i.NumSelected
	items.Height = i.Height
	items.Width = i.Width
	items.Cur = i.Cur
	items.SetHeader = i.SetHeader
	items.SetFooter = i.SetFooter
	items.Title = i.Title
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
	i.Cur = idx
}

func (i Items) CurrentItem() *Item {
	return &i.Items[i.Cur]
}

func (cp Items) Visible(matches ...string) []Item {
	if len(matches) != 0 {
		return ExactMatches(matches[0], cp.Items)
	}
	return cp.Items
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
	m.Items = ChoiceMapToMatch(choices)
}

func (m Items) RenderItems(cursor int, items []Item) string {
	var s strings.Builder
	for i, match := range items {
		pre := "x"

		if match.Label != "" {
			pre = match.Label
		}

		switch {
		case i == m.CurrentItem().Index:
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

		s.WriteString(match.Render(m.Width, m.Height))
		s.WriteRune('\n')
	}
	view := s.String()
	//m.SetBody(view)
	return view
}

func (m Items) Exec() error {
	if m.cmd != "" {
		for _, c := range m.Chosen() {
			for _, arg := range c {
				args := slices.Clone(m.args)
				args = append(args, arg)
				cmd := exec.Command(m.cmd, args...)
				var (
					stderr bytes.Buffer
					stdout bytes.Buffer
				)
				cmd.Stderr = &stderr
				cmd.Stdout = &stdout

				//fmt.Println(cmd.String())
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("command exited with error: %s\n", stderr.String())
				}

				if out := stdout.String(); out != "" {
					fmt.Println(out)
				}

				if err := stderr.String(); err != "" {
					return fmt.Errorf("command exited with error: %s\n%s\n", err, stdout.String())
				}
			}
		}
	}
	return nil
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

func Width(w int) Opt {
	return func(i *Items) {
		i.Width = w
	}
}

func Size(w, h int) Opt {
	return func(i *Items) {
		i.Width = w
		i.Height = h
	}
}

func Exec(cmd string, args ...string) Opt {
	return func(i *Items) {
		i.args = args
		i.cmd = cmd
	}
}

func Header(t string) Opt {
	return func(i *Items) {
		i.Title = t
	}
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
				text := lipgloss.NewStyle().Width(i.Width).Render(fmt.Sprint(val))
				item := NewItem(text, idx)
				item.Label = fmt.Sprint(label)
				i.Items[idx] = item
			}
		}
	}
}
