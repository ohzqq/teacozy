package teacozy

import (
	"github.com/charmbracelet/bubbles/paginator"
)

type State struct {
	Model      paginator.Model
	Items      Items
	Selected   map[int]struct{}
	InputValue string
	ReadOnly   bool

	cursor  int
	total   int
	start   int
	end     int
	current int
}

func New(items Items) *State {
	m := &State{
		Items: items,
		Model: paginator.New(),
	}
	m.SetPerPage(10)
	m.SetTotal(items.Len())
	return m
}

func (m State) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m State) Len() int {
	return m.total
}

func (m State) Current() int {
	return m.current
}

func (m State) Start() int {
	return m.start
}

func (m State) End() int {
	return m.end
}

func (m *State) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m *State) SetCursor(n int) *State {
	m.cursor = clamp(n, 0, m.end-1)
	return m
}

func (m *State) SetCurrent(n int) {
	m.current = n
}

func (m *State) SetTotal(n int) {
	m.total = n
	m.Model.SetTotalPages(n)
	m.updateSliceBounds()
}

func (m *State) SetPerPage(n int) {
	m.Model.PerPage = n
	m.Model.SetTotalPages(m.total)
	m.updateSliceBounds()
}

func (m *State) SetInputValue(val string) {
	m.InputValue = val
}

func (m State) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.Model.PerPage {
			return i
		}
	}
	return 0
}

func (m *State) updateSliceBounds() {
	m.start, m.end = m.Model.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
}

func (m *State) SliceBounds() (int, int) {
	m.updateSliceBounds()
	return m.start, m.end
}

func (m State) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m State) OnFirstItem() bool {
	return m.cursor == 0
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
