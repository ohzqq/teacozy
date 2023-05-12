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
)

type Pager struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[PagerProps]
	Model paginator.Model

	cursor int
	total  int
	start  int
	end    int
	KeyMap keys.KeyMap
	slug   string
}

type PagerProps struct {
	name       string
	Items      Items
	Selected   map[int]struct{}
	Search     string
	ReadOnly   bool
	SetCurrent func(int)
	Style      Style
}

func NewPager() *Pager {
	c := &Pager{
		KeyMap: keys.DefaultKeyMap(),
		Model:  paginator.New(),
		Style:  DefaultStyle(),
	}
	return c
}

func NewPagerProps(items teacozy.Items) PagerProps {
	return PagerProps{
		Items: items,
	}
}

func (c *Pager) Init(props PagerProps) tea.Cmd {
	c.UpdateProps(props)
	c.Model.SetTotalPages(c.Props().Items.Len())
	c.SetPerPage(10)
	c.SliceBounds()
	return nil
}

func (m *Pager) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case keys.PageUpMsg:
		m.cursor = clamp(m.cursor-m.Model.PerPage, 0, m.total-1)
		m.Model.PrevPage()
	case keys.PageDownMsg:
		m.cursor = clamp(m.cursor+m.Model.PerPage, 0, m.total-1)
		m.Model.NextPage()
	case keys.HalfPageUpMsg:
		m.HalfUp()
		if m.cursor < m.start {
			m.Model.PrevPage()
		}
	case keys.HalfPageDownMsg:
		m.HalfDown()
		if m.cursor >= m.end {
			m.Model.NextPage()
		}
	case keys.LineDownMsg:
		m.NextItem()
		if m.cursor >= m.end {
			m.Model.NextPage()
		}
	case keys.LineUpMsg:
		m.PrevItem()
		if m.cursor < m.start {
			m.Model.PrevPage()
		}
	case keys.TopMsg:
		m.cursor = 0
		m.Model.Page = 0
	case keys.BottomMsg:
		m.cursor = m.total - 1
		m.Model.Page = m.Model.TotalPages - 1
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return tea.Quit
		}
		for _, k := range m.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.SliceBounds()

	return tea.Batch(cmds...)
}

func (c *Pager) Render(w, h int) string {
	var s strings.Builder
	h = h - 2

	// get matched items
	items := teacozy.ExactMatches(c.Props().Search, c.Props().Items)

	c.SetPerPage(h)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	c.SetTotal(len(items))

	for i, m := range items[c.Start():c.End()] {
		var cur bool
		if i == c.Highlighted() {
			c.Props().SetCurrent(m.Index)
			cur = true
		}

		var sel bool
		if _, ok := c.Props().Selected[m.Index]; ok {
			sel = true
		}

		label := c.Props().Items.Label(m.Index)
		pre := c.PrefixText(label, sel, cur)
		style := c.PrefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if !c.Props().ReadOnly || label != "" {
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
	m.KeyMap = km
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
	m.KeyMap = keys.NewKeyMap(keys.Quit())
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
