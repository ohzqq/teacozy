package pagy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/keys"
)

type Paginator struct {
	paginator.Model

	cursor int
	total  int
	start  int
	end    int
	Index  int
	KeyMap keys.KeyMap
}

func New(per, total int) *Paginator {
	m := &Paginator{
		KeyMap: keys.DefaultKeyMap(),
		total:  total,
		Model:  paginator.New(),
	}
	m.PerPage = per
	m.SetTotalPages(total)
	m.SliceBounds()

	return m
}

func (m *Paginator) Update(msg tea.Msg) (*Paginator, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case keys.PageUpMsg:
		m.cursor = clamp(m.cursor-m.PerPage, 0, m.total-1)
		m.PrevPage()
	case keys.PageDownMsg:
		m.cursor = clamp(m.cursor+m.PerPage, 0, m.total-1)
		m.NextPage()
	case keys.HalfPageUpMsg:
		m.HalfUp()
		if m.cursor < m.start {
			m.PrevPage()
		}
	case keys.HalfPageDownMsg:
		m.HalfDown()
		if m.cursor >= m.end {
			m.NextPage()
		}
	case keys.LineDownMsg:
		m.NextItem()
		if m.cursor >= m.end {
			m.NextPage()
		}
	case keys.LineUpMsg:
		m.PrevItem()
		if m.cursor < m.start {
			m.PrevPage()
		}
	case keys.TopMsg:
		m.cursor = 0
		m.Page = 0
	case keys.BottomMsg:
		m.cursor = m.total - 1
		m.Page = m.TotalPages - 1
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		for _, k := range m.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.SliceBounds()

	return m, tea.Batch(cmds...)
}

func (m Paginator) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m *Paginator) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m Paginator) Len() int {
	return m.total
}

func (m Paginator) Current() int {
	return m.Index
}

func (m Paginator) Start() int {
	return m.start
}

func (m Paginator) End() int {
	return m.end
}

func (m *Paginator) SetCursor(n int) *Paginator {
	m.cursor = n
	return m
}

func (m *Paginator) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = km
}

func (m *Paginator) SetCurrent(n int) *Paginator {
	m.Index = n
	return m
}

func (m *Paginator) SetTotal(n int) *Paginator {
	m.total = n
	m.SetTotalPages(n)
	m.SliceBounds()
	return m
}

func (m *Paginator) SetPerPage(n int) *Paginator {
	m.PerPage = n
	m.SetTotalPages(m.total)
	m.SliceBounds()
	return m
}

func (m Paginator) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.PerPage {
			return i
		}
	}
	return 0
}

func (m *Paginator) SliceBounds() (int, int) {
	m.start, m.end = m.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m Paginator) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m Paginator) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *Paginator) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *Paginator) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *Paginator) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Paginator) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Paginator) DisableKeys() {
	m.KeyMap = keys.NewKeyMap(keys.Quit())
}

func (m *Paginator) Init() tea.Cmd {
	return nil
}

func (m *Paginator) View() string {
	return m.Model.View()
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
