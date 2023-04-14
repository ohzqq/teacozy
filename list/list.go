package list

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	Cursor    int
	Paginator paginator.Model
	Viewport  *viewport.Model
	quitting  bool
	Style     style.List
	items     *props.Items
	end       int
	start     int
	Footer    string
}

func New(props *props.Items) *List {
	tm := &List{
		Style: style.ListDefaults(),
		items: props,
	}

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(props.Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height

	v := viewport.New(props.Width, props.Height)
	tm.Viewport = &v

	return tm
}

func (m List) KeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Up().WithKeys("k", "up"),
		keys.Down().WithKeys("j", "down"),
		keys.Next().WithKeys("right", "l"),
		keys.Prev().WithKeys("left", "h"),
		keys.NewBinding("G").
			WithHelp("list bottom").
			Cmd(message.Bottom),
		keys.NewBinding("g").
			WithHelp("list top").
			Cmd(message.Top),
	}
	return km
}

func (m *List) Props() *props.Items {
	return m.items
}

func (m *List) Update(msg tea.Msg) (*List, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))
	switch msg := msg.(type) {
	case message.NextMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor+m.Props().Height)
		m.SetCurrent()
		m.Paginator.NextPage()
		m.Viewport.GotoTop()

	case message.PrevMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor-m.Props().Height)
		m.SetCurrent()
		m.Paginator.PrevPage()
		m.Viewport.GotoBottom()

	case message.TopMsg:
		m.Cursor = 0
		m.SetCurrent()
		m.Paginator.Page = 0
		m.Viewport.GotoTop()

	case message.BottomMsg:
		m.Cursor = len(m.Props().Visible()) - 1
		m.SetCurrent()
		m.Paginator.Page = m.Paginator.TotalPages - 1
		m.Viewport.GotoBottom()

	case message.LineUpMsg:
		m.Cursor--
		if m.Cursor < start {
			if m.Paginator.Page > 0 {
				m.Cursor = len(m.Props().Visible()) - 1
				m.Paginator.PrevPage()
				m.Viewport.GotoBottom()
				cmds = append(cmds, message.PrevPage)
			} else {
				m.Cursor = 0
			}
		}
		m.SetCurrent()
		h := m.Props().CurrentItem().LineHeight()
		if m.Props().Lines > m.Props().Height {
			m.Viewport.LineUp(h)
		}

	case message.LineDownMsg:
		m.Cursor++
		if m.Cursor >= end {
			if m.Paginator.OnLastPage() {
				m.Cursor = len(m.Props().Visible()) - 1
			} else {
				m.Cursor = 0
				cmds = append(cmds, message.NextPage)
			}
		}
		m.SetCurrent()
		h := m.Props().CurrentItem().LineHeight()
		if m.Props().Lines > m.Props().Height {
			m.Viewport.LineDown(h)
		}

	case tea.KeyMsg:
		for _, k := range m.KeyMap() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *List) SetCurrent() {
	m.Props().SetCurrent(m.Cursor % m.Props().Height)
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *List) MoveUp(n int) {
	m.Cursor = clamp(m.Cursor-n, 0, len(m.Props().Visible())-1)
	switch {
	case m.start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
	case m.start < m.Viewport.Height:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Viewport.Height)
	}
	m.UpdateViewport()
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *List) MoveDown(n int) {
	m.Cursor = clamp(m.Cursor+n, 0, len(m.Props().Visible())-1)
	m.UpdateViewport()

	switch {
	case m.end == len(m.Props().Visible()):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Viewport.Height))
	case m.Cursor > (m.end-m.start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Viewport.Height-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *List) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.Props().Visible()))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Viewport.Height, 0, m.Cursor)
	} else {
		m.start = 0
	}
	m.end = clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.Props().Visible()))
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, strconv.Itoa(i))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *List) View() string {
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	items := m.Props().RenderItems(
		m.Props().Visible()[start:end],
	)
	m.Viewport.SetContent(items)

	view := m.Viewport.View()

	if m.Paginator.TotalPages > 1 {
		p := style.Footer.Render(m.Paginator.View())
		m.Footer = p
	}

	return view
}

func (m List) Init() tea.Cmd { return nil }

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
