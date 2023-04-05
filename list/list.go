package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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
	KeyMap    keys.KeyMap
	items     *props.Items
	Footer    string
}

var Keys = keys.KeyMap{
	keys.Up().WithKeys("k", "up"),
	keys.Down().WithKeys("j", "down"),
	keys.Next().WithKeys("right", "l"),
	keys.Prev().WithKeys("left", "h"),
	keys.NewBinding("G").
		WithHelp("list bottom").
		Cmd(message.Bottom()),
	keys.NewBinding("g").
		WithHelp("list top").
		Cmd(message.Top()),
}

func (m *List) Props() *props.Items {
	return m.items
}

func New(props *props.Items) *List {
	tm := &List{
		Style:  style.ListDefaults(),
		KeyMap: Keys,
		items:  props,
	}

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(props.Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height

	v := viewport.New(props.Width, props.Height)
	tm.Viewport = &v

	return tm
}

func (m *List) Update(msg tea.Msg) (*List, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))
	switch msg := msg.(type) {
	case message.NextMsg:
		//m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor+m.Props().Height)
		m.Cursor = 0
		m.SetCurrent()
		m.Paginator.NextPage()
		m.Viewport.GotoTop()

	case message.PrevMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor-m.Props().Height)
		//m.Cursor = 0
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

	case message.UpMsg:
		m.Cursor--
		if m.Cursor < start {
			if m.Paginator.Page > 0 {
				m.Cursor = len(m.Props().Visible()) - 1
				m.Paginator.PrevPage()
				m.Viewport.GotoBottom()
				//cmds = append(cmds, message.Prev())
			} else {
				m.Cursor = 0
			}
		}
		m.SetCurrent()
		h := m.Props().CurrentItem().LineHeight()
		if m.Props().Lines > m.Props().Height {
			m.Viewport.LineUp(h)
		}

	case message.DownMsg:
		m.Cursor++
		if m.Cursor >= end {
			if m.Paginator.OnLastPage() {
				m.Cursor = len(m.Props().Visible()) - 1
			} else {
				m.Cursor = 0
				cmds = append(cmds, message.Next())
			}
		}
		m.SetCurrent()
		h := m.Props().CurrentItem().LineHeight()
		if m.Props().Lines > m.Props().Height {
			m.Viewport.LineDown(h)
		}

	case tea.KeyMsg:
		for _, k := range m.KeyMap {
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
