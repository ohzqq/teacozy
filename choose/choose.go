package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
	Cursor    int
	Paginator paginator.Model
	quitting  bool
	header    string
	Style     style.List
}

type Props struct {
	*props.Items
}

func NewChoice() *Choose {
	tm := Choose{
		Style: style.ListDefaults(),
	}
	return &tm
}

func RouteInitializer(props Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewChoice()
		return component, component.Init(props)
	}
}

func NewProps(items *props.Items) Props {
	return Props{
		Items: items,
	}
}

func (c Choose) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewChoice()
		return component, component.Init(NewProps(props))
	}
}

func (c Choose) Name() string {
	return "choose"
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	switch msg := msg.(type) {
	case message.NextMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor+m.Props().Height)
		m.Props().Items.SetCurrent(m.Cursor)
		m.Paginator.NextPage()

	case message.PrevMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor-m.Props().Height)
		m.Props().Items.SetCurrent(m.Cursor)
		m.Paginator.PrevPage()

	case message.TopMsg:
		m.Cursor = 0
		m.Paginator.Page = 0
		m.Props().SetCurrent(m.Cursor)

	case message.BottomMsg:
		m.Cursor = len(m.Props().Visible()) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
		m.Props().SetCurrent(m.Cursor)

	case message.UpMsg:
		m.Cursor--
		if m.Cursor < 0 {
			m.Cursor = len(m.Props().Visible()) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		}
		if m.Cursor < start {
			m.Paginator.PrevPage()
		}
		m.Props().SetCurrent(m.Cursor)

	case message.DownMsg:
		m.Cursor++
		if m.Cursor >= len(m.Props().Visible()) {
			m.Cursor = 0
			m.Paginator.Page = 0
		}
		if m.Cursor >= end {
			m.Paginator.NextPage()
		}
		m.Props().SetCurrent(m.Cursor)

	case message.StartEditingMsg:
		cur := m.Props().Visible()[m.Cursor]
		m.Props().SetCurrent(cur.Index)
		return message.ChangeRouteCmd("editField")

	case message.ToggleItemMsg:
		idx := m.Props().Visible()[m.Cursor].Index

		if m.Props().NumSelected == 0 && m.quitting {
			cmds = append(cmds, message.ReturnSelectionsCmd())
		}

		m.Props().ToggleSelection(idx)

		if m.Props().Limit == 1 {
			return message.ReturnSelectionsCmd()
		}

		cmds = append(cmds, message.DownCmd())

	case message.StartFilteringMsg:
		return message.ChangeRouteCmd("filter")

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Key.Up):
			cmds = append(cmds, message.UpCmd())
		case key.Matches(msg, Key.Down):
			cmds = append(cmds, message.DownCmd())
		case key.Matches(msg, Key.Prev):
			cmds = append(cmds, message.PrevCmd())
		case key.Matches(msg, Key.Next):
			cmds = append(cmds, message.NextCmd())
		case key.Matches(msg, Key.ToggleItem):
			cmds = append(cmds, message.ToggleItemCmd())
		case key.Matches(msg, Key.Edit):
			cmds = append(cmds, message.StartEditingCmd())
		case key.Matches(msg, Key.Filter):
			cmds = append(cmds, message.StartFilteringCmd())
		case key.Matches(msg, Key.Bottom):
			cmds = append(cmds, message.BottomCmd())
		case key.Matches(msg, Key.Top):
			cmds = append(cmds, message.TopCmd())
		case key.Matches(msg, Key.Quit):
			m.quitting = true
			cmds = append(cmds, message.ReturnSelectionsCmd())
		case key.Matches(msg, Key.ReturnSelections):
			if m.Props().Limit == 1 {
				return message.ToggleItemCmd()
			}
			if m.Props().NumSelected == 0 {
				m.quitting = true
				return message.ToggleItemCmd()
			}
			cmds = append(cmds, message.ReturnSelectionsCmd())
		}
	}

	return tea.Batch(cmds...)
}

func (m *Choose) Render(w, h int) string {
	var s strings.Builder

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	items := m.Props().RenderItems(
		m.Cursor%m.Props().Height,
		m.Props().Visible()[start:end],
	)
	s.WriteString(items)

	var view string
	view = s.String()
	if m.Paginator.TotalPages <= 1 {
		//view = s.String()
	} else if m.Paginator.TotalPages > 1 {
		p := style.Footer.Render(m.Paginator.View())
		//view = lipgloss.JoinVertical(lipgloss.Left, view, p)
		m.Props().Footer(p)
	}

	return view
}

func (tm *Choose) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)
	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(tm.Props().Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height
	return nil
}
