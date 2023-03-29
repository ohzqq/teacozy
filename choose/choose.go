package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ChooseProps]
	item.Items
	Choices     []string
	choiceMap   []map[string]string
	Viewport    *viewport.Model
	Paginator   paginator.Model
	ListKeys    func(m *Choose) keys.KeyMap
	numSelected int
	limit       int
	aborted     bool
	quitting    bool
	header      string
	Placeholder string
	Prompt      string
	Width       int
	Height      int
	Style       style.List
}

func New(choices ...string) *Choose {
	tm := Choose{
		Choices:  choices,
		ListKeys: ListKeyMap,
		Style:    DefaultStyle(),
		Prompt:   style.PromptPrefix,
	}
	return &tm
}

func (m *Choose) Run() []int {
	//p := tea.NewProgram(m)
	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}

	//if m.quitting {
	//return []int{}
	//}
	return m.Chosen()
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Make place in the view port if header is set
		if m.header != "" {
			m.Viewport.Height = m.Viewport.Height - lipgloss.Height(m.Style.Header.Render(m.header))
		}
		m.Viewport.Width = msg.Width
	case ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case StartFilteringMsg:
		reactea.SetCurrentRoute("filter")
		return nil
	case tea.KeyMsg:
		for _, k := range GlobalKeyMap(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
		for _, k := range m.ListKeys(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *Choose) CursorUp() int {
	start, _ := m.Paginator.GetSliceBounds(len(m.Props().Visible()))
	m.Cursor--
	if m.Cursor < 0 {
		m.Cursor = len(m.Props().Visible()) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
	}
	if m.Cursor < start {
		m.Paginator.PrevPage()
	}
	return m.Cursor
}

func (m *Choose) CursorDown() int {
	_, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))
	m.Cursor++
	if m.Cursor >= len(m.Props().Visible()) {
		m.Cursor = 0
		m.Paginator.Page = 0
	}
	if m.Cursor >= end {
		m.Paginator.NextPage()
	}
	return m.Cursor
}

func (m *Choose) ToggleSelection() {
	idx := m.Props().Visible()[m.Cursor].Index
	m.Props().ToggleItem(idx)
	m.CursorDown()
}

func (m *Choose) Render(w, h int) string {
	m.Viewport.Height = h
	if m.Paginator.TotalPages > 1 {
		m.Viewport.Height = m.Viewport.Height + 4
	}
	m.Viewport.Width = w
	return m.View()
}

func (m *Choose) View() string {
	var s strings.Builder

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	for i, match := range m.Props().Visible()[start:end] {
		pre := "x"

		if match.Label != "" {
			pre = match.Label
		}

		switch {
		case i == m.Cursor%m.Props().Height:
			pre = match.Style.Cursor.Render(pre)
		default:
			if _, ok := m.Props().Selected[match.Index]; ok {
				pre = match.Style.Selected.Render(pre)
			} else if match.Label == "" {
				pre = strings.Repeat(" ", lipgloss.Width(pre))
			} else {
				pre = match.Style.Label.Render(pre)
			}
		}

		s.WriteString("[")
		s.WriteString(pre)
		s.WriteString("]")
		s.WriteString(match.RenderText())
		s.WriteRune('\n')
	}

	var view string
	if m.Paginator.TotalPages <= 1 {
		view = s.String()
	} else if m.Paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.Props().Height-m.Paginator.ItemsOnPage(len(m.Props().Visible()))+1))
		s.WriteString("  " + m.Paginator.View())
	}

	view = s.String()

	m.Viewport.SetContent(view)
	view = m.Viewport.View()

	return view
}

//nolint:unparam
func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (tm *Choose) Init(props ChooseProps) tea.Cmd {
	tm.UpdateProps(props)
	v := viewport.New(0, 0)
	tm.Viewport = &v
	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Dots
	tm.Paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	tm.Paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)
	tm.Paginator.SetTotalPages((len(tm.Props().Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height
	return nil
}
