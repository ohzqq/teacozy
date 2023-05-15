package cmpnt

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
)

type Pager struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[*teacozy.Page]
	Model paginator.Model

	Width  int
	Height int

	cursor int
	total  int
	start  int
	end    int

	keyMap keys.KeyMap
	Style  Style
}

func New() *Pager {
	c := &Pager{
		Width:  util.TermWidth(),
		Height: util.TermHeight() - 2,
		Model:  paginator.New(),
	}
	c.SetKeyMap(keys.VimKeyMap())
	c.AddKey(keys.Help())
	return c
}

func (c *Pager) Initializer(props *teacozy.Page) teacozy.PageComponent {
	c.Init(props)
	return c
}

func (c *Pager) Mount() reactea.SomeComponent {
	return c
}

func (c *Pager) Init(props *teacozy.Page) tea.Cmd {
	c.UpdateProps(props)
	c.SetPerPage(c.Height)
	c.SetTotal(c.Props().Items().Len())
	c.SliceBounds()
	return nil
}

func (c *Pager) Update(msg tea.Msg) tea.Cmd {
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

func (c *Pager) Render(w, h int) string {
	var s strings.Builder
	//h = h - 2
	h = c.Height - 2

	// get matched items
	//items := c.ExactMatches(c.Search)
	items := teacozy.SourceToMatches(c.Props().Items())

	c.SetPerPage(h)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	c.SetTotal(len(items))

	props := NewItems(items[c.Start():c.End()])
	props.ReadOnly = false
	props.Highlighted = c.Highlighted
	props.IsSelected = c.Props().IsSelected
	props.GetLabel = c.Props().Items().Label
	view := props.Render()
	s.WriteString(view)

	return s.String()
}

func (c *Pager) AddKey(k *keys.Binding) *Pager {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (c *Pager) SetWidth(n int) *Pager {
	c.Width = n
	return c
}

func (c *Pager) SetHeight(n int) *Pager {
	c.Height = n
	return c
}

func (c *Pager) SetSize(w, h int) *Pager {
	c.Width = w
	c.Height = h
	return c
}

func (m Pager) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m *Pager) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m Pager) Len() int {
	return m.total
}

func (m Pager) Start() int {
	return m.start
}

func (m Pager) End() int {
	return m.end
}

func (m *Pager) SetCursor(n int) *Pager {
	m.cursor = n
	return m
}

func (m *Pager) SetKeyMap(km keys.KeyMap) {
	m.keyMap = km
}

func (m *Pager) SetTotal(n int) *Pager {
	m.total = n
	m.Model.SetTotalPages(n)
	m.SliceBounds()
	return m
}

func (m *Pager) SetPerPage(n int) *Pager {
	m.Model.PerPage = n
	m.Model.SetTotalPages(m.total)
	m.SliceBounds()
	return m
}

func (m Pager) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.Model.PerPage {
			m.Props().SetCurrent(m.cursor)
			return i
		}
	}
	return 0
}

func (m *Pager) SliceBounds() (int, int) {
	m.start, m.end = m.Model.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m Pager) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m Pager) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *Pager) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *Pager) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *Pager) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Pager) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Pager) DisableKeys() {
	m.keyMap = keys.NewKeyMap(keys.Quit())
}

func (m Pager) KeyMap() keys.KeyMap {
	return m.keyMap
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
