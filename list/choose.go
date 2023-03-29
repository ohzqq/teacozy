package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/style"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ListProps]

	Viewport  *viewport.Model
	Paginator paginator.Model
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

func (tm *Choose) Init(props ListProps) tea.Cmd {
	tm.UpdateProps(props)
	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Dots
	tm.Paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	tm.Paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)
	return nil
}

func (tm *Choose) Render(w, h int) string {
	v := viewport.New(w, h+4)
	tm.Viewport = &v
	tm.Paginator.SetTotalPages((len(tm.Props().Matches) + h - 1) / h)
	tm.Paginator.PerPage = h

	var s strings.Builder

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Matches))

	for _, item := range m.Props().Matches[start:end] {
		s.WriteString(item.Str)
		s.WriteRune('\n')
	}
	//items := item.RenderItems(m.Props().Cursor, m.Props().Matches[start:end])
	//s.WriteString(items)

	var view string
	if m.Paginator.TotalPages <= 1 {
		view = s.String()
	} else if m.Paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", h-m.Paginator.ItemsOnPage(len(m.Props().Matches))+1))
		s.WriteString("  " + m.Paginator.View())
	}

	view = s.String()
	if m.header != "" {
		header := m.Style.Header.Render(m.header + strings.Repeat(" ", m.Width))
		view = lipgloss.JoinVertical(lipgloss.Left, header, view)
	}

	m.Viewport.SetContent(view)
	view = m.Viewport.View()

	return view
}
