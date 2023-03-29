package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/style"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ChooseProps]
	Cursor    int
	Choices   []string
	choiceMap []map[string]string
	Viewport  *viewport.Model
	Paginator paginator.Model
	aborted   bool
	quitting  bool
	header    string
	Style     style.List
}

type ChooseKeys struct {
	Up               key.Binding
	Down             key.Binding
	Prev             key.Binding
	Next             key.Binding
	ToggleItem       key.Binding
	Quit             key.Binding
	ReturnSelections key.Binding
	Filter           key.Binding
	Bottom           key.Binding
	Top              key.Binding
}

func New(choices ...string) *Choose {
	tm := Choose{
		Choices: choices,
		Style:   DefaultStyle(),
	}
	return &tm
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Make place in the view port if header is set
		if m.header != "" {
			m.Viewport.Height = m.Viewport.Height - lipgloss.Height(m.Style.Header.Render(m.header))
		}
		m.Viewport.Width = msg.Width
	case UpMsg:
		m.Cursor--
		if m.Cursor < 0 {
			m.Cursor = len(m.Props().Visible()) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		}
		if m.Cursor < start {
			m.Paginator.PrevPage()
		}
	case DownMsg:
		m.Cursor++
		if m.Cursor >= len(m.Props().Visible()) {
			m.Cursor = 0
			m.Paginator.Page = 0
		}
		if m.Cursor >= end {
			m.Paginator.NextPage()
		}
	case StartFilteringMsg:
		reactea.SetCurrentRoute("filter")
		return nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, chooseKey.Up):
			cmds = append(cmds, UpCmd())
		case key.Matches(msg, chooseKey.Down):
			cmds = append(cmds, DownCmd())
		case key.Matches(msg, chooseKey.Prev):
			m.Cursor = clamp(0, len(m.Props().Items.Items)-1, m.Cursor-m.Props().Height)
			m.Paginator.PrevPage()
		case key.Matches(msg, chooseKey.Next):
			m.Cursor = clamp(0, len(m.Props().Items.Items)-1, m.Cursor+m.Props().Height)
			m.Paginator.NextPage()
		case key.Matches(msg, chooseKey.ToggleItem):
			if m.Props().Limit == 1 {
				return nil
			}
			idx := m.Props().Visible()[m.Cursor].Index
			m.Props().ToggleItem(idx)
			cmds = append(cmds, DownCmd())
		case key.Matches(msg, chooseKey.Filter):
			reactea.SetCurrentRoute("filter")
			return nil
		case key.Matches(msg, chooseKey.Bottom):
			m.Cursor = len(m.Props().Items.Items) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		case key.Matches(msg, chooseKey.Top):
			m.Cursor = 0
			m.Paginator.Page = 0
		case key.Matches(msg, chooseKey.Quit):
			m.quitting = true
			cmds = append(cmds, ReturnSelectionsCmd())
		case key.Matches(msg, chooseKey.ReturnSelections):
			cmds = append(cmds, ReturnSelectionsCmd())
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

	items := m.Props().RenderItems(
		m.Cursor%m.Props().Height,
		m.Props().Visible()[start:end],
	)
	s.WriteString(items)

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

func (m *Choose) Run() []int {
	//p := tea.NewProgram(m)
	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}

	//if m.quitting {
	return []int{}
	//}
	//return m.Chosen()

}
