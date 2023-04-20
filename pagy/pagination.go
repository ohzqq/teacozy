package pagy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/keys"
)

type Model struct {
	paginator.Model

	cursor int
	total  int
	start  int
	end    int
	KeyMap keys.KeyMap
}

func New(per, total int) *Model {
	m := &Model{
		KeyMap: DefaultKeyMap(),
		total:  total,
		Model:  paginator.New(),
	}
	m.PerPage = per
	m.SetTotalPages(total)

	return m
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		for _, k := range m.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.SliceBounds()

	return m, tea.Batch(cmds...)
}

func (m Model) Cursor() int {
	return m.cursor
}

func (m Model) Len() int {
	return m.total
}

func (m Model) Start() int {
	return m.start
}

func (m Model) End() int {
	return m.end
}

func (m *Model) SetCursor(n int) *Model {
	m.cursor = n
	return m
}

func (m *Model) SetKeyMap(km keys.KeyMap) *Model {
	m.KeyMap = km
	return m
}

func (m *Model) SetTotal(n int) *Model {
	m.total = n
	m.SetTotalPages(n)
	return m
}

func (m *Model) SetPerPage(n int) *Model {
	m.PerPage = n
	return m
}

func (m Model) Current() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.PerPage {
			return i
		}
	}
	return 0
}

func (m *Model) SliceBounds() (int, int) {
	m.start, m.end = m.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m Model) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m Model) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *Model) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *Model) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *Model) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Model) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return m.Model.View()
}

func DefaultKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.PgUp(),
		keys.PgDown(),
		keys.Up(),
		keys.Down(),
		keys.HalfPgUp(),
		keys.HalfPgDown(),
		keys.Home(),
		keys.End(),
		keys.Quit(),
	}
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
