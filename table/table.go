package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

	rows        []string
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

// New creates a new model for the table widget.
func New(opts ...Option) *Model {
	m := Model{
		Cursor:   0,
		Viewport: viewport.New(0, 20),

		focus: true,

		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}

	for _, opt := range opts {
		opt(&m)
	}

	//m.UpdateViewport()

	return &m
}

func (c Model) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (c Model) Name() string {
	return "filter"
}

func (m Model) Km() keys.KeyMap {
	var km = keys.KeyMap{
		keys.Quit(),
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

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.QuitMsg:
		return tea.Quit
	case message.DownMsg:
		//m.Cursor = clamp(m.Cursor+msg.Lines, 0, len(m.rows)-1)
		//m.Props().SetCurrent(m.Cursor)
		//m.UpdateViewport()

		//switch {
		//case m.end == len(m.rows):
		//  m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-msg.Lines, 1, m.Viewport.Height))
		//case m.Cursor > (m.end-m.start)/2:
		//  m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-msg.Lines, 1, m.Cursor))
		//case m.Viewport.YOffset > 1:
		//case m.Cursor > m.Viewport.YOffset+m.Viewport.Height-1:
		//  m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
		//}

		m.MoveDown(msg.Lines)
	case message.UpMsg:
		m.MoveUp(msg.Lines)
	case message.TopMsg:
		m.GotoTop()
	case message.BottomMsg:
		m.GotoBottom()
	case tea.KeyMsg:
		for _, k := range m.Km() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the component.
func (m Model) View() string {
	return m.Viewport.View()
}

func (m *Model) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height
	m.Viewport.Width = m.Props().Width
	m.UpdateViewport()
	//m.Viewport.SetContent(m.Props().RenderItems(m.Matches))

	view := m.Input.View() + "\n" + m.Viewport.View()
	return view
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *Model) UpdateViewport() {
	renderedRows := make([]props.Item, 0, len(m.Props().Visible()))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Viewport.Height, 0, m.Cursor)
	} else {
		m.start = 0
	}
	m.end = clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.rows))
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.Props().Items.Items[i])
	}

	//m.Props().SetCurrent(m.Cursor)
	m.Viewport.SetContent(m.Props().RenderItems(renderedRows))
	//lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	//)
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m Model) SelectedRow() string {
	if m.Cursor < 0 || m.Cursor >= len(m.rows) {
		return ""
	}

	return m.rows[m.Cursor]
}

// Rows returns the current rows.
func (m Model) Rows() []string {
	return m.rows
}

// SetRows sets a new rows state.
func (m *Model) SetRows(r []string) {
	m.rows = r
	//m.UpdateViewport()
}

// SetWidth sets the width of the viewport of the table.
func (m *Model) SetWidth(w int) {
	m.Viewport.Width = w
	//m.UpdateViewport()
}

// SetHeight sets the height of the viewport of the table.
func (m *Model) SetHeight(h int) {
	m.Viewport.Height = h
	//m.UpdateViewport()
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
	return m.Cursor
}

// SetCursor sets the cursor position in the table.
func (m *Model) SetCursor(n int) {
	m.Cursor = clamp(n, 0, len(m.rows)-1)
	//m.UpdateViewport()
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
	m.Cursor = clamp(m.Cursor-n, 0, len(m.rows)-1)
	switch {
	case m.start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
	case m.start < m.Viewport.Height:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Viewport.Height)
	}
	m.Props().SetCurrent(m.Cursor)
	//m.UpdateViewport()
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Model) MoveDown(n int) {
	m.Cursor = clamp(m.Cursor+n, 0, len(m.rows)-1)
	//m.UpdateViewport()

	switch {
	case m.end == len(m.rows):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Viewport.Height))
	case m.Cursor > (m.end-m.start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Viewport.Height-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}

	m.Props().SetCurrent(m.Cursor)
}

// GotoTop moves the selection to the first row.
func (m *Model) GotoTop() {
	m.MoveUp(m.Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model) GotoBottom() {
	m.MoveDown(len(m.rows))
}

//func (m Model) Init() tea.Cmd { return nil }

func (tm *Model) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)

	tm.Input = textinput.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.PromptStyle = tm.Style.Prompt
	tm.Input.Placeholder = tm.Placeholder
	tm.Input.Width = tm.Props().Width

	tm.Viewport = viewport.New(0, 0)
	tm.Matches = props.Visible()
	tm.rows = props.Slice()
	tm.focus = true

	tm.UpdateViewport()

	tm.Input.Blur()

	return nil
}

func (m *Model) renderRow(rowID int) string {
	row := m.rows[rowID]

	//row := lipgloss.JoinHorizontal(lipgloss.Left, s...)

	if rowID == m.Cursor {
		return m.Style.Match.Render(row)
	}

	return row
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

// WithRows sets the table rows (data).
func WithRows(rows []string) Option {
	return func(m *Model) {
		m.rows = rows
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
	m.UpdateViewport()
}

// Blur blurs the table, preventing selection or movement.
func (m *Model) Blur() {
	m.focus = false
	m.UpdateViewport()
}
