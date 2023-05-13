package cmpnt

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
)

type Filter struct {
	reactea.BasicComponent
	Model paginator.Model

	ConfirmChoices bool
	readOnly       bool

	Width  int
	Height int

	NumSelected int
	Limit       int
	CurrentItem int
	NoLimit     bool
	ReadOnly    bool

	cursor int
	total  int
	start  int
	end    int

	Choices teacozy.Items
	keyMap  keys.KeyMap
	Style   Style

	help keys.KeyMap

	teacozy.State
}

type FilterProps struct {
	ReadOnly   bool
	Selected   map[int]struct{}
	SetCurrent func(int)
}

func NewFilter(choices teacozy.Items) *Filter {
	c := &Filter{
		Width:   util.TermWidth(),
		Height:  util.TermHeight() - 2,
		Limit:   10,
		Model:   paginator.New(),
		Style:   DefaultStyle(),
		Choices: choices,
	}

	if c.NoLimit {
		c.Limit = c.Choices.Len()
	}

	if !c.readOnly {
		c.AddKey(keys.Toggle().AddKeys(" "))
	}

	c.SetPerPage(c.Height)
	c.SetTotal(c.Choices.Len())
	c.SliceBounds()

	c.SetKeyMap(keys.VimKeyMap())

	c.AddKey(keys.Help())

	return c
}

func (c *Filter) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.PageUpMsg:
		c.cursor = clamp(c.cursor-c.Model.PerPage, 0, c.total-1)
		c.Model.PrevPage()
	case keys.PageDownMsg:
		c.cursor = clamp(c.cursor+c.Model.PerPage, 0, c.total-1)
		c.Model.NextPage()
	case keys.HalfPageUpMsg:
		c.HalfUp()
		if c.cursor < c.start {
			c.Model.PrevPage()
		}
	case keys.HalfPageDownMsg:
		c.HalfDown()
		if c.cursor >= c.end {
			c.Model.NextPage()
		}
	case keys.LineDownMsg:
		c.NextItem()
		if c.cursor >= c.end {
			c.Model.NextPage()
		}
	case keys.LineUpMsg:
		c.PrevItem()
		if c.cursor < c.start {
			c.Model.PrevPage()
		}
	case keys.TopMsg:
		c.cursor = 0
		c.Model.Page = 0
	case keys.BottomMsg:
		c.cursor = c.total - 1
		c.Model.Page = c.Model.TotalPages - 1

	case keys.ShowHelpMsg:
		cmds = append(cmds, keys.ChangeRoute("help"))
		//help := NewProps(c.help)
		//help.SetName("help")
		//return ChangeRoute(&help)

	case keys.UpdateItemMsg:
		return msg.Cmd(c.Current())

	case keys.ToggleItemsMsg, keys.ToggleItemMsg:
		c.ToggleItems(c.Current())
		cmds = append(cmds, keys.LineDown)

	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}

		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	c.SliceBounds()

	return tea.Batch(cmds...)
}

func (c *Filter) Render(w, h int) string {
	var s strings.Builder
	//h = h - 2
	h = c.Height - 2

	// get matched items
	items := c.ExactMatches(c.Search)

	c.SetPerPage(h)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	c.SetTotal(len(items))

	for i, m := range items[c.Start():c.End()] {
		var cur bool
		if i == c.Highlighted() {
			c.SetCurrent(m.Index)
			cur = true
		}

		var sel bool
		if _, ok := c.Selected[m.Index]; ok {
			sel = true
		}

		label := c.Items.Label(m.Index)
		pre := c.State.PrefixText(label, sel, cur)
		style := c.State.PrefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if !c.ReadOnly || label != "" {
			s.WriteString(style.Render(pre))
		}

		// render the rest of the line
		text := lipgloss.StyleRunes(
			m.Str,
			m.MatchedIndexes,
			c.Style.Match,
			c.Style.Normal.Style,
		)

		s.WriteString(lipgloss.NewStyle().Render(text))
		s.WriteString("\n")
	}

	return s.String()
}

func (c Filter) ToggleItem() {
	c.ToggleItems(c.Current())
}

func (c *Filter) AddKey(k *keys.Binding) *Filter {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (c *Filter) ToggleItems(items ...int) {
	for _, idx := range items {
		c.CurrentItem = idx
		if _, ok := c.Selected[idx]; ok {
			delete(c.Selected, idx)
			c.NumSelected--
		} else if c.NumSelected < c.Limit {
			c.Selected[idx] = struct{}{}
			c.NumSelected++
		}
	}
}

func (m Filter) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			l := m.Choices.Label(k)
			v := m.Choices.String(k)
			chosen = append(chosen, map[string]string{l: v})
		}
	}
	return chosen
}

func (m *Filter) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *Filter) Current() int {
	return m.CurrentItem
}

func (c *Filter) SetWidth(n int) *Filter {
	c.Width = n
	return c
}

func (c *Filter) SetHelp(km keys.KeyMap) {
	c.help = km
}

func (c *Filter) SetHeight(n int) *Filter {
	c.Height = n
	return c
}

func (c *Filter) SetSize(w, h int) *Filter {
	c.Width = w
	c.Height = h
	return c
}

func (m Filter) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m *Filter) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m Filter) Len() int {
	return m.total
}

func (m Filter) Start() int {
	return m.start
}

func (m Filter) End() int {
	return m.end
}

func (m *Filter) SetCursor(n int) *Filter {
	m.cursor = n
	return m
}

func (m *Filter) SetKeyMap(km keys.KeyMap) {
	m.keyMap = km
}

func (m *Filter) SetTotal(n int) *Filter {
	m.total = n
	m.Model.SetTotalPages(n)
	m.SliceBounds()
	return m
}

func (m *Filter) SetPerPage(n int) *Filter {
	m.Model.PerPage = n
	m.Model.SetTotalPages(m.total)
	m.SliceBounds()
	return m
}

func (m Filter) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.Model.PerPage {
			return i
		}
	}
	return 0
}

func (m *Filter) SliceBounds() (int, int) {
	m.start, m.end = m.Model.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m Filter) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m Filter) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *Filter) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *Filter) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *Filter) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Filter) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Filter) DisableKeys() {
	m.keyMap = keys.NewKeyMap(keys.Quit())
}

func (m Filter) KeyMap() keys.KeyMap {
	return m.keyMap
}
