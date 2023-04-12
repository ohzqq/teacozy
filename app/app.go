package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/match"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Option func(*List)

type List struct {
	reactea.BasicComponent // It implements AfterUpdate() for us, so we don't have to care!
	reactea.BasicPropfulComponent[reactea.NoProps]

	items       *match.Component
	choices     []map[string]string
	focus       bool
	quitting    bool
	Selected    map[int]struct{}
	Cursor      int
	height      int
	width       int
	NumSelected int
	Limit       int
	Style       style.List
	KeyMap      keys.KeyMap
	props       *Props
	matches     []string

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
	Items   []Item
	Matches []Item
}

func New(props []string, opts ...Option) *List {
	m := List{
		focus:  true,
		width:  util.TermWidth(),
		height: util.TermHeight(),
		Style:  style.ListDefaults(),
		props: &Props{
			Items:   ChoicesToMatch(props),
			Matches: ChoicesToMatch(props),
		},
		items:   match.New(),
		choices: MapChoices(props),
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.Viewport = viewport.New(m.width, m.height)

	m.KeyMap = m.DefaultKeyMap()
	m.UpdateItems()

	return &m
}

//func (m *List) UpdateProps(props *props.Items) {
//  m.props = props
//}

func (m *List) Props() *Props {
	return m.props
}

//func (m *List) Init() tea.Cmd {
//  return nil
//}

func (c *List) Init(reactea.NoProps) tea.Cmd {
	//reactea.SetCurrentRoute("list")
	return nil
}

func (m List) DefaultKeyMap() keys.KeyMap {
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
		cmds = append(cmds, tea.Quit)

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
			for _, k := range m.KeyMap {
				if key.Matches(msg, k.Binding) {
					cmds = append(cmds, k.TeaCmd)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m List) View() string {
	return m.Viewport.View()
}

func (m *List) Render(w, h int) string {
	m.SetWidth(m.width)
	m.SetHeight(m.height)
	m.UpdateItems()

	renderedRows := make([]string, 0, len(m.matches))

	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.matches[i])
	}
	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
	return m.View()
}

// UpdateItems updates the list content based on the previously defined
// columns and rows.
func (m *List) UpdateItems() {
	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Viewport.Height, 0, m.Cursor)
	} else {
		m.start = 0
		m.SetCursor(0)
	}
	m.end = clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.matches))

	if m.Cursor > m.end {
		m.SetCursor(clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.matches)-1))
	}

	props := match.Props{
		Choices:  m.choices,
		Selected: m.Selected,
		Cursor:   m.Cursor,
		Matches:  m.SetMatches,
		Search:   "b",
	}
	m.items.Init(props)
	m.items.Render(m.Viewport.Width, m.Viewport.Height)

}

func (m *List) renderRow(rowID int) string {
	row := m.matches[rowID]

	var s strings.Builder

	s.WriteString(row)

	return s.String()
}

func (m *List) SetMatches(matches []string) {
	m.matches = matches
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m List) CurrentItem() Item {
	return m.Props().Matches[m.Cursor]
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *List) MoveUp(n int) {
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
func (m *List) MoveDown(n int) {
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
func (m *List) GotoTop() {
	m.MoveUp(m.Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *List) GotoBottom() {
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
	m.Props().Matches = m.Props().Items
	m.UpdateItems()
}

// Blur blurs the table, preventing selection or movement.
func (m *List) Blur() {
	m.focus = false
	m.UpdateItems()
}

// VisibleItems returns the current rows.
func (m List) VisibleItems() []Item {
	return m.Props().Matches
}

// SetItems sets a new rows state.
func (m *List) SetItems(r []Item) {
	m.Props().Matches = r
	m.UpdateItems()
}

// SetWidth sets the width of the viewport of the table.
func (m *List) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateItems()
}

// SetHeight sets the height of the viewport of the table.
func (m *List) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateItems()
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
	return m.Cursor
}

// SetCursor sets the cursor position in the table.
func (m *List) SetCursor(n int) {
	//m.Props().SetCursor(clamp(n, 0, len(m.Props().Matches)-1))
	m.Cursor = clamp(n, 0, len(m.Props().Matches)-1)
	m.UpdateItems()
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
