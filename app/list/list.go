package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/app/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/style"
)

type Option func(*Component)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	focus    bool
	quitting bool
	Cursor   int
	height   int
	width    int
	Style    style.List
	KeyMap   keys.KeyMap

	Viewport viewport.Model
	start    int
	end      int
}

type Choices []map[string]string

func (i Choices) String(idx int) string {
	var str string
	item := i[idx]
	for _, v := range item {
		str = v
	}
	return str
}

func (i Choices) Len() int {
	return len(i)
}

type Props struct {
	Matches     []item.Item
	Selected    map[int]struct{}
	ToggleItems func(...int)
	SetContent  func(string)
}

func NewList(opts ...Option) *Component {
	m := Component{
		Cursor: 0,
		Style:  style.ListDefaults(),
	}
	m.KeyMap = keys.DefaultListKeyMap()

	for _, opt := range opts {
		opt(&m)
	}

	return &m
}

func (m *Component) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.Viewport = viewport.New(0, 0)
	m.UpdateItems()
	return nil
}

func (m *Component) AfterUpdate() tea.Cmd {
	m.UpdateItems()
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.QuitMsg:
		cmds = append(cmds, tea.Quit)

	case message.ToggleItemMsg:
		cur := m.Props().Matches[m.Cursor].Index
		m.Props().ToggleItems(cur)
		m.MoveDown(1)

	case keys.PageUpMsg:
		m.MoveUp(m.Viewport.Height)
	case keys.PageDownMsg:
		m.MoveDown(m.Viewport.Height)
	case keys.HalfPageUpMsg:
		m.MoveUp(m.Viewport.Height / 2)
	case keys.HalfPageDownMsg:
		m.MoveDown(m.Viewport.Height / 2)
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

	return tea.Batch(cmds...)
}

func (m *Component) Render(w, h int) string {
	m.SetWidth(w)
	m.SetHeight(h)
	m.UpdateItems()
	return m.Viewport.View()
}

// UpdateItems updates the list content based on the previously defined
// columns and rows.
func (m *Component) UpdateItems() {
	renderedRows := make([]string, 0, len(m.Props().Matches))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Viewport.Height, 0, m.Cursor)
	} else {
		m.start = 0
		m.SetCursor(0)
	}
	m.end = clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.Props().Matches))

	if m.Cursor > m.end {
		m.SetCursor(clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.Props().Matches)-1))
	}

	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *Component) renderRow(rowID int) string {
	row := m.Props().Matches[rowID]

	var s strings.Builder

	switch {
	case rowID == m.Cursor:
		row.Current = true
	case m.isSelected(rowID):
		row.Selected = true
	}

	s.WriteString(row.Render(m.Viewport.Width, m.Viewport.Height))

	return s.String()
}

func (m Component) isSelected(idx int) bool {
	_, ok := m.Props().Selected[m.Props().Matches[idx].Index]
	return ok
}

// SelectedRow returns the selected row.
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
	case m.start < m.Viewport.Height:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Viewport.Height)
	}
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Component) MoveDown(n int) {
	m.SetCursor(clamp(m.Cursor+n, 0, len(m.Props().Matches)-1))
	m.UpdateItems()
	switch {
	case m.end == len(m.Props().Matches):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Viewport.Height))
	case m.Cursor > (m.end-m.start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Viewport.Height-1:
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

//func (m *List) ToggleAllItems() tea.Cmd {
//  return func() tea.Msg {
//    var items []int
//    for _, item := range m.Props().AllItems() {
//      items = append(items, item.Index)
//    }
//    m.Props().ToggleSelection(items...)
//    return nil
//  }
//}

func (m *Component) quit() tea.Cmd {
	m.quitting = true
	return message.Quit()
}

// VisibleItems returns the current rows.
func (m Component) VisibleItems() []item.Item {
	return m.Props().Matches
}

// SetItems sets a new rows state.
func (m *Component) SetItems(r []item.Item) {
	//m.Props().Matches = r
	m.UpdateItems()
}

// SetWidth sets the width of the viewport of the table.
func (m *Component) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateItems()
}

// SetHeight sets the height of the viewport of the table.
func (m *Component) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateItems()
}

func (m *Component) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = km
}

// Height returns the viewport height of the table.
func (m Component) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the table.
func (m Component) Width() int {
	return m.Viewport.Width
}

// Cursor returns the index of the selected row.
func (m Component) GetCursor() int {
	return m.Cursor
}

// SetCursor sets the cursor position in the table.
func (m *Component) SetCursor(n int) {
	m.Cursor = clamp(n, 0, len(m.Props().Matches)-1)
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
