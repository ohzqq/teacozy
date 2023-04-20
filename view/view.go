package view

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
)

type Model struct {
	Cursor int
	KeyMap keys.KeyMap

	Viewport viewport.Model
	Start    int
	End      int
	props    Props
}

type Props struct {
	item.Props
	Editable   bool
	Filterable bool
	Width      int
	Height     int
}

func NewModel(props Props) *Model {
	m := Model{
		Cursor: props.Cursor,
		Start:  props.Start,
		End:    props.End,
		props:  props,
	}
	m.SetKeyMap(DefaultKeyMap())
	m.Viewport = viewport.New(props.Width, props.Height)
	return &m
}

func NewView() *Model {
	m := Model{
		Cursor: 0,
	}
	m.SetKeyMap(DefaultKeyMap())
	m.Viewport = viewport.New(0, 0)
	return &m
}

func NewProps(c item.Choices) Props {
	return Props{
		Props: item.Props{
			Choices:  c,
			Selected: make(map[int]struct{}),
			Start:    0,
			End:      10,
			Cursor:   0,
		},
	}
}

func (m *Model) Props() Props {
	return m.props
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
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

func (m Model) View() string {
	return m.Viewport.View()
}

func (m Model) ItemProps() item.Props {
	p := item.Props{
		Choices:  m.Props().Choices,
		Start:    m.Start,
		End:      m.End,
		Selected: m.Props().Selected,
	}
	p.Cursor = m.Cursor
	return p
}

// UpdateItems updates the list content based on the previously defined
// columns and rows.
func (m *Model) UpdateItems() {

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height

	if m.Cursor >= 0 {
		m.Start = clamp(m.Cursor-m.Height(), 0, m.Cursor)
	} else {
		m.Start = 0
		m.SetCursor(0)
	}

	m.End = clamp(m.Cursor+m.Height(), m.Cursor, len(m.Props().Matches()))
	if m.Cursor > m.End {
		m.SetCursor(clamp(m.Cursor+m.Height(), m.Cursor, len(m.Props().Matches())-1))
	}

	l := item.Renderer(m.ItemProps(), m.Width(), m.Height())
	m.Viewport.SetContent(l)
}

// CurrentItem returns the selected row.
// You can cast it to your own implementation.
func (m Model) CurrentItem() int {
	return m.Props().Matches()[m.Cursor].Index
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
	m.SetCursor(clamp(m.Cursor-n, 0, len(m.Props().Matches())-1))
	m.UpdateItems()
	switch {
	case m.Start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
	case m.Start < m.Height():
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Height())
	}
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Model) MoveDown(n int) {

	m.SetCursor(clamp(m.Cursor+n, 0, len(m.Props().Matches())-1))
	m.UpdateItems()
	switch {
	case m.End == len(m.Props().Matches()):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Height()))
	case m.Cursor > (m.End-m.Start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Height()-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}

}

// GotoTop moves the selection to the first row.
func (m *Model) GotoTop() {
	m.MoveUp(m.Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model) GotoBottom() {
	m.MoveDown(len(m.Props().Matches()))
}

// SetWidth sets the width of the viewport of the list.
func (m *Model) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateItems()
}

// SetHeight sets the height of the viewport of the list.
func (m *Model) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateItems()
}

// SetKeyMap sets the keymap for the list.
func (m *Model) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = km
}

// Height returns the viewport height of the list.
func (m Model) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the list.
func (m Model) Width() int {
	return m.Viewport.Width
}

// Cursor returns the index of the selected item.
func (m Model) GetCursor() int {
	return m.Cursor
}

// SetCursor sets the cursor position in the list.
func (m *Model) SetCursor(n int) {
	m.Cursor = clamp(n, 0, len(m.Props().Matches())-1)
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
