package cmpnt

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	*Pager

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

type ListProps struct {
	PagerProps
	Selected   func() map[int]struct{}
	ToggleItem func(int)
}

func NewList(p *Pager, choices teacozy.Items) *List {
	c := &List{
		Width:   util.TermWidth(),
		Height:  util.TermHeight() - 2,
		Limit:   10,
		Style:   DefaultStyle(),
		Choices: choices,
	}

	if c.NoLimit {
		c.Limit = c.Choices.Len()
	}

	c.AddKey(keys.Toggle().AddKeys(" "))

	c.SetPerPage(c.Height)
	c.SetTotal(c.Choices.Len())
	c.SliceBounds()

	c.SetKeyMap(keys.VimKeyMap())

	c.AddKey(keys.Help())

	return c
}

func (c *List) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (c List) ToggleItem() {
	c.ToggleItems(c.Current())
}

func (c *List) AddKey(k *keys.Binding) *List {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (c *List) ToggleItems(items ...int) {
	for _, idx := range items {
		c.Props().SetCurrent(idx)
		if _, ok := c.Selected[idx]; ok {
			delete(c.Selected, idx)
			c.NumSelected--
		} else if c.NumSelected < c.Limit {
			c.Selected[idx] = struct{}{}
			c.NumSelected++
		}
	}
}

func (m List) Chosen() []map[string]string {
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

func (m *List) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *List) Current() int {
	return m.CurrentItem
}

func (c *List) SetWidth(n int) *List {
	c.Width = n
	return c
}

func (c *List) SetHelp(km keys.KeyMap) {
	c.help = km
}

func (c *List) SetHeight(n int) *List {
	c.Height = n
	return c
}

func (c *List) SetSize(w, h int) *List {
	c.Width = w
	c.Height = h
	return c
}

func (m List) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m *List) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m List) Len() int {
	return m.total
}

func (m List) Start() int {
	return m.start
}

func (m List) End() int {
	return m.end
}

func (m *List) SetCursor(n int) *List {
	m.cursor = n
	return m
}

func (m *List) SetKeyMap(km keys.KeyMap) {
	m.keyMap = km
}

func (m *List) SetTotal(n int) *List {
	m.total = n
	m.Model.SetTotalPages(n)
	m.SliceBounds()
	return m
}

func (m *List) SetPerPage(n int) *List {
	m.Model.PerPage = n
	m.Model.SetTotalPages(m.total)
	m.SliceBounds()
	return m
}

func (m List) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.Model.PerPage {
			return i
		}
	}
	return 0
}

func (m *List) SliceBounds() (int, int) {
	m.start, m.end = m.Model.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m List) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m List) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *List) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *List) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *List) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *List) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *List) DisableKeys() {
	m.keyMap = keys.NewKeyMap(keys.Quit())
}

func (m List) KeyMap() keys.KeyMap {
	return m.keyMap
}
