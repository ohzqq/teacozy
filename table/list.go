package table

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

// List defines a state for the table widget.
type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Matches     []props.Item
	Cursor      int
	focus       bool
	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List

	Viewport viewport.Model
	start    int
	end      int
}

type option func(*List)

func (c List) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (c List) Name() string {
	return "list"
}

// New creates a new model for the table widget.
func New(opts ...option) *List {
	m := List{
		Cursor: 0,

		focus: true,

		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return &m
}

func (m *List) Init(props Props) tea.Cmd {
	m.UpdateProps(props)

	m.Viewport = viewport.New(props.Width, props.Height)
	m.Matches = props.Visible()

	return nil
}

func (m List) KeyMap() keys.KeyMap {
	var km = keys.KeyMap{
		keys.Quit(),
		keys.ToggleItem(),
		keys.Up().WithKeys("up"),
		keys.Down().WithKeys("down"),
		keys.NewBinding("ctrl+u").
			WithHelp("½ page up").
			Cmd(message.Up(m.Viewport.Height / 2)),
		keys.NewBinding("ctrl+d").
			WithHelp("½ page down").
			Cmd(message.Down(m.Viewport.Height / 2)),
		keys.NewBinding("pgup").
			WithHelp("page up").
			Cmd(message.Up(m.Viewport.Height)),
		keys.NewBinding("pgdown").
			WithHelp("page down").
			Cmd(message.Down(m.Viewport.Height)),
		keys.NewBinding("end").
			WithHelp("list bottom").
			Cmd(message.Bottom()),
		keys.NewBinding("home").
			WithHelp("list top").
			Cmd(message.Top()),
	}
	return km
}

func (m *List) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.QuitMsg:
		return tea.Quit

	case message.ToggleItemMsg:
		if len(m.Matches) > 0 {
			m.Props().ToggleSelection()
			cmds = append(cmds, message.Down())
		}

	case message.DownMsg:
		m.MoveDown(msg.Lines)
	case message.UpMsg:
		m.MoveUp(msg.Lines)
	case message.TopMsg:
		m.GotoTop()
	case message.BottomMsg:
		m.GotoBottom()
	case tea.KeyMsg:
		if m.Focused() {
			for _, k := range m.KeyMap() {
				if key.Matches(msg, k.Binding) {
					cmds = append(cmds, k.TeaCmd)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m List) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height - 1
	m.Viewport.Width = m.Props().Width
	m.UpdateRows()

	view := m.Viewport.View()

	return view
}

func (m List) View() string {
	return m.Viewport.View()
}

// UpdateRows updates the list content based on the previously defined
// columns and rows.
func (m *List) UpdateRows() {
	renderedRows := make([]string, 0, len(m.Matches))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Props().Cursor >= 0 {
		m.start = clamp(m.Props().Cursor-m.Viewport.Height, 0, m.Props().Cursor)
	} else {
		m.start = 0
	}
	m.end = clamp(m.Props().Cursor+m.Viewport.Height, m.Props().Cursor, len(m.Matches))
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *List) renderRow(rowID int) string {
	row := m.Matches[rowID]

	var s strings.Builder
	pre := "x"

	if row.Label != "" {
		pre = row.Label
	}

	switch {
	case rowID == m.Props().Cursor:
		pre = row.Style.Cursor.Render(pre)
	default:
		if _, ok := m.Props().Selected[row.Index]; ok {
			pre = row.Style.Selected.Render(pre)
		} else if row.Label == "" {
			pre = strings.Repeat(" ", lipgloss.Width(pre))
		} else {
			pre = row.Style.Label.Render(pre)
		}
	}

	s.WriteString("[")
	s.WriteString(pre)
	s.WriteString("]")

	s.WriteString(row.Render(m.Viewport.Width, m.Viewport.Height))

	return s.String()
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m List) SelectedRow() props.Item {
	row := m.Props().GetItem(m.Matches[m.Props().Cursor].Index)
	return row
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *List) MoveUp(n int) {
	m.Props().SetCursor(clamp(m.Props().Cursor-n, 0, len(m.Matches)-1))
	m.UpdateRows()
	switch {
	case m.start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Props().Cursor))
	case m.start < m.Viewport.Height:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Props().Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Viewport.Height)
	}
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *List) MoveDown(n int) {
	m.Props().SetCursor(clamp(m.Props().Cursor+n, 0, len(m.Matches)-1))
	m.UpdateRows()
	switch {
	case m.end == len(m.Matches):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Viewport.Height))
	case m.Props().Cursor > (m.end-m.start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Props().Cursor))
	case m.Viewport.YOffset > 1:
	case m.Props().Cursor > m.Viewport.YOffset+m.Viewport.Height-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *List) GotoTop() {
	m.MoveUp(m.Props().Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *List) GotoBottom() {
	m.MoveDown(len(m.Matches))
}

func (m *List) ToggleAllItems() tea.Cmd {
	return func() tea.Msg {
		var items []int
		for _, item := range m.Props().AllItems() {
			items = append(items, item.Index)
		}
		m.Props().ToggleSelection(items...)
		return nil
	}
}

func (m *List) quit() tea.Cmd {
	m.quitting = true
	return message.Quit()
}

// Focused returns the focus state of the table.
func (m List) Focused() bool {
	return m.focus
}

// Focus focuses the table, allowing the user to move around the rows and
// interact.
func (m *List) Focus() {
	m.focus = true
	m.UpdateRows()
}

// Blur blurs the table, preventing selection or movement.
func (m *List) Blur() {
	m.focus = false
	m.UpdateRows()
}

// VisibleItems returns the current rows.
func (m List) VisibleItems() []props.Item {
	return m.Matches
}

// SetItems sets a new rows state.
func (m *List) SetItems(r []props.Item) {
	m.Matches = r
	m.UpdateRows()
}

// SetWidth sets the width of the viewport of the table.
func (m *List) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateRows()
}

// SetHeight sets the height of the viewport of the table.
func (m *List) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateRows()
}

// Height returns the viewport height of the table.
func (m List) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the table.
func (m List) Width() int {
	return m.Viewport.Width
}

// Cursor returns the index of the selected row.
func (m List) GetCursor() int {
	return m.Props().Cursor
}

// SetCursor sets the cursor position in the table.
func (m *List) SetCursor(n int) {
	m.Props().SetCursor(clamp(n, 0, len(m.Matches)-1))
	m.UpdateRows()
}
