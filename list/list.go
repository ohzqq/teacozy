package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
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
	return tm
}

func (m *List) Update(msg tea.Msg) (*List, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))
	switch msg := msg.(type) {
	case message.NextMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor+m.Props().Height)
		m.Props().SetCurrent(m.Cursor % m.Props().Height)
		m.Paginator.NextPage()

	case message.PrevMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor-m.Props().Height)
		m.Props().SetCurrent(m.Cursor % m.Props().Height)
		m.Paginator.PrevPage()

	case message.TopMsg:
		m.Cursor = 0
		m.Paginator.Page = 0
		m.Props().SetCurrent(m.Cursor % m.Props().Height)

	case message.BottomMsg:
		m.Cursor = len(m.Props().Visible()) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
		m.Props().SetCurrent(m.Cursor % m.Props().Height)

	case message.UpMsg:
		m.Cursor--
		if m.Cursor < 0 {
			m.Cursor = len(m.Props().Visible()) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		}
		if m.Cursor < start {
			m.Paginator.PrevPage()
		}
		m.Props().SetCurrent(m.Cursor % m.Props().Height)

	case message.DownMsg:
		m.Cursor++
		if m.Cursor >= len(m.Props().Visible()) {
			m.Cursor = 0
			m.Paginator.Page = 0
		}
		if m.Cursor >= end {
			m.Paginator.NextPage()
		}
		m.Props().SetCurrent(m.Cursor % m.Props().Height)

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

func (m *List) View() string {
	var s strings.Builder
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	items := m.Props().RenderItems(
		m.Cursor%m.Props().Height,
		m.Props().Visible()[start:end],
	)
	s.WriteString(items)

	view := s.String()

	if m.Paginator.TotalPages > 1 {
		p := style.Footer.Render(m.Paginator.View())
		m.Footer = p
	}

	return view
}

func (m List) Init() tea.Cmd { return nil }
