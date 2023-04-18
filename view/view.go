package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	Cursor int
	KeyMap keys.KeyMap

	Viewport viewport.Model
	start    int
	end      int
	props    Props
}

type Props struct {
	Editable   bool
	Filterable bool
	Selectable bool
	Selected   map[int]struct{}
	Matches    []item.Item
}

func New(props Props) *Component {
	m := Component{
		Cursor: 0,
		props:  props,
	}
	m.SetKeyMap(DefaultKeyMap())
	m.Viewport = viewport.New(0, 0)
	m.UpdateItems()

	return &m
}

func (m *Component) Props() Props {
	return m.props
}

func (m *Component) Init() tea.Cmd {
	return nil
}

func (m *Component) Update(msg tea.Msg) (*Component, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.PageUpMsg:
		m.MoveUp(m.Height())
	case keys.PageDownMsg:
		m.MoveDown(m.Height())
	case keys.HalfPageUpMsg:
		m.MoveUp(m.Height() / 2)
	case keys.HalfPageDownMsg:
		m.MoveDown(m.Height() / 2)
	case keys.LineDownMsg:
		m.MoveDown(1)
	case keys.LineUpMsg:
		m.MoveUp(1)
	case keys.TopMsg:
		m.GotoTop()
	case keys.BottomMsg:
		m.GotoBottom()
	case tea.KeyMsg:
		for _, k := range m.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Component) View() string {
	return m.Viewport.View()
}

// UpdateItems updates the list content based on the previously defined
// columns and rows.
func (m *Component) UpdateItems() {
	items := make([]string, 0, len(m.Props().Matches))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Height(), 0, m.Cursor)
	} else {
		m.start = 0
		m.SetCursor(0)
	}
	m.end = clamp(m.Cursor+m.Height(), m.Cursor, len(m.Props().Matches))

	if m.Cursor > m.end {
		m.SetCursor(clamp(m.Cursor+m.Height(), m.Cursor, len(m.Props().Matches)-1))
	}

	for i := m.start; i < m.end; i++ {
		items = append(items, m.renderItem(i))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}

func (m *Component) renderItem(rowID int) string {
	item := m.Props().Matches[rowID]

	var s strings.Builder

	if m.Props().Selectable {
		switch {
		case rowID == m.Cursor:
			item.Current = true
		case m.IsSelected(rowID):
			item.Selected = true
		}
	}

	s.WriteString(item.Render(m.Width(), m.Height()))

	return s.String()
}

func (m Component) IsSelected(idx int) bool {
	_, ok := m.Props().Selected[m.Props().Matches[idx].Index]
	return ok
}

// CurrentItem returns the selected row.
// You can cast it to your own implementation.
func (m Component) CurrentItem() int {
	return m.Props().Matches[m.Cursor].Index
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Component) MoveUp(n int) {
	m.SetCursor(clamp(m.Cursor-n, 0, len(m.Props().Matches)-1))
	m.UpdateItems()
	switch {
	case m.start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
	case m.start < m.Height():
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Height())
	}
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Component) MoveDown(n int) {
	m.SetCursor(clamp(m.Cursor+n, 0, len(m.Props().Matches)-1))
	m.UpdateItems()
	switch {
	case m.end == len(m.Props().Matches):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Height()))
	case m.Cursor > (m.end-m.start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Height()-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *Component) GotoTop() {
	m.MoveUp(m.Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Component) GotoBottom() {
	m.MoveDown(len(m.Props().Matches))
}

// SetWidth sets the width of the viewport of the list.
func (m *Component) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateItems()
}

// SetHeight sets the height of the viewport of the list.
func (m *Component) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateItems()
}

// SetKeyMap sets the keymap for the list.
func (m *Component) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = km
}

// Height returns the viewport height of the list.
func (m Component) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the list.
func (m Component) Width() int {
	return m.Viewport.Width
}

// Cursor returns the index of the selected item.
func (m Component) GetCursor() int {
	return m.Cursor
}

// SetCursor sets the cursor position in the list.
func (m *Component) SetCursor(n int) {
	m.Cursor = clamp(n, 0, len(m.Props().Matches)-1)
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

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
