package table

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
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

// Model defines a state for the table widget.
type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Matches     []props.Item
	Cursor      int
	focus       bool
	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List

	Input textinput.Model

	Viewport viewport.Model
	start    int
	end      int
}

type Props struct {
	*props.Items
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option func(*Model)

func (c Model) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (c Model) Name() string {
	return "table"
}

// New creates a new model for the table widget.
func NewTable(opts ...Option) *Model {
	m := Model{
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

func (m *Model) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.Input = textinput.New()
	m.Input.Prompt = m.Prompt
	m.Input.PromptStyle = m.Style.Prompt
	m.Input.Placeholder = m.Placeholder
	m.Input.Width = props.Width
	m.Input.Focus()

	m.Viewport = viewport.New(props.Width, props.Height)
	m.Matches = props.Visible()

	return nil
}

func (m Model) KeyMap() keys.KeyMap {
	var km = keys.KeyMap{
		keys.Quit(),
		keys.ToggleItem(),
		keys.Up().WithKeys("up"),
		keys.Down().WithKeys("down"),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.NewBinding("/").
			WithHelp("filter list").
			Cmd(StartFiltering()),
		keys.NewBinding("esc").
			WithHelp("stop filtering").
			Cmd(StopFiltering()),
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

func (m *Model) ReturnSelections() tea.Cmd {
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.QuitMsg:
		return tea.Quit
	case StopFilteringMsg:
		m.Input.Reset()
		m.Input.Blur()
		cur := m.SelectedRow().Index
		m.Matches = m.Props().Visible()
		m.SetCursor(cur)
		m.Props().SetFooter(strconv.Itoa(m.Props().Cursor))

	case StartFilteringMsg:
		m.Input.Focus()
	case message.ToggleItemMsg:
		if len(m.Matches) > 0 {
			m.Props().SetCurrent(m.Matches[m.Props().Cursor].Index)
			if m.Props().NumSelected == 0 && m.quitting {
				cmds = append(cmds, m.ReturnSelections())
			}
			m.Props().ToggleSelection()
			if m.Props().Limit == 1 {
				return m.ReturnSelections()
			}
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
		if m.Input.Focused() {
			for _, k := range m.KeyMap() {
				if key.Matches(msg, k.Binding) {
					cmds = append(cmds, k.TeaCmd)
				}
			}
			m.Input, cmd = m.Input.Update(msg)
			if v := m.Input.Value(); v != "" {
				m.Matches = m.Props().Visible(v)
			} else {
				m.Matches = m.Props().Visible()
			}
			cmds = append(cmds, cmd)
		} else {
			for _, k := range m.UnfilteredKeyMap() {
				if key.Matches(msg, k.Binding) {
					cmds = append(cmds, k.TeaCmd)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m Model) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height - 1
	m.Viewport.Width = m.Props().Width
	m.UpdateRows()

	view := m.Viewport.View()
	if m.Input.Focused() {
		view = m.Input.View() + "\n" + view
	}

	return view
}

// UpdateRows updates the list content based on the previously defined
// columns and rows.
func (m *Model) UpdateRows() {
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

func (m *Model) renderRow(rowID int) string {
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
func (m Model) SelectedRow() props.Item {
	row := m.Props().GetItem(m.Matches[m.Props().Cursor].Index)
	return row
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
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
func (m *Model) MoveDown(n int) {
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
func (m *Model) GotoTop() {
	m.MoveUp(m.Props().Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model) GotoBottom() {
	m.MoveDown(len(m.Matches))
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

func (m *Model) ToggleAllItems() tea.Cmd {
	return func() tea.Msg {
		var items []int
		for _, item := range m.Props().Items.Items {
			items = append(items, item.Index)
		}
		m.Props().ToggleSelection(items...)
		return nil
	}
}

func (m *Model) quit() tea.Cmd {
	m.quitting = true
	return message.Quit()
}

func (m Model) UnfilteredKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Up().WithKeys("k", "up"),
		keys.Down().WithKeys("j", "down"),
		keys.Next().WithKeys("right", "l"),
		keys.Prev().WithKeys("left", "h"),
		keys.NewBinding("/").
			WithHelp("filter list").
			Cmd(StartFiltering()),
		keys.NewBinding("G").
			WithHelp("list bottom").
			Cmd(message.Bottom()),
		keys.NewBinding("g").
			WithHelp("list top").
			Cmd(message.Top()),
		keys.NewBinding("v").
			WithHelp("toggle all items").
			Cmd(m.ToggleAllItems()),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.ToggleItem().WithKeys("tab", " "),
		keys.ShowHelp(),
		keys.Quit().
			WithKeys("ctrl+c", "q", "esc").
			Cmd(m.quit()),
	}
	return km
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}

// WithRows sets the table rows (data).
func WithRows(rows *props.Items) Option {
	return func(m *Model) {
		m.Matches = rows.Visible()
	}
}

// WithHeight sets the height of the table.
func WithHeight(h int) Option {
	return func(m *Model) {
		m.Viewport.Height = h
	}
}

// WithWidth sets the width of the table.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.Viewport.Width = w
	}
}

// WithFocused sets the focus state of the table.
func WithFocused(f bool) Option {
	return func(m *Model) {
		m.focus = f
	}
}

// WithStyles sets the table styles.
func WithStyles(s style.List) Option {
	return func(m *Model) {
		m.Style = s
	}
}

// Focused returns the focus state of the table.
func (m Model) Focused() bool {
	return m.focus
}

// Focus focuses the table, allowing the user to move around the rows and
// interact.
func (m *Model) Focus() {
	m.focus = true
	m.UpdateRows()
}

// Blur blurs the table, preventing selection or movement.
func (m *Model) Blur() {
	m.focus = false
	m.UpdateRows()
}

// Rows returns the current rows.
func (m Model) Rows() []props.Item {
	return m.Matches
}

// SetRows sets a new rows state.
func (m *Model) SetRows(r []props.Item) {
	m.Matches = r
	m.UpdateRows()
}

// SetWidth sets the width of the viewport of the table.
func (m *Model) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateRows()
}

// SetHeight sets the height of the viewport of the table.
func (m *Model) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateRows()
}

// Height returns the viewport height of the table.
func (m Model) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the table.
func (m Model) Width() int {
	return m.Viewport.Width
}

// Cursor returns the index of the selected row.
func (m Model) GetCursor() int {
	return m.Props().Cursor
}

// SetCursor sets the cursor position in the table.
func (m *Model) SetCursor(n int) {
	m.Props().SetCursor(clamp(n, 0, len(m.Matches)-1))
	m.UpdateRows()
}

type StopFilteringMsg struct{}

func StopFiltering() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}

type StartFilteringMsg struct{}

func StartFiltering() tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}
