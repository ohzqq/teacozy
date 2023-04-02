package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/list"
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
	KeyMap    keys.KeyMap
	list      *list.List
	Style     style.List
}

type Props struct {
	*props.Items
}

func NewChoice() *Choose {
	tm := Choose{
		Style:  style.ListDefaults(),
		KeyMap: Keys,
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
	var cmd tea.Cmd
	var cmds []tea.Cmd

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	switch msg := msg.(type) {
	case message.ShowHelpMsg:
		k := keys.Global
		k = append(k, list.Keys...)
		m.Props().Help(k)
		cmds = append(cmds, message.ChangeRoute("help"))
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
		return message.ChangeRoute("editField")

	case message.ToggleItemMsg:
		idx := m.Props().Visible()[m.Cursor].Index

		if m.Props().NumSelected == 0 && m.quitting {
			cmds = append(cmds, message.ReturnSelections())
		}

		m.Props().ToggleSelection(idx)

		if m.Props().Limit == 1 {
			return message.ReturnSelections()
		}

		cmds = append(cmds, message.Down())

	case message.StartFilteringMsg:
		return message.ChangeRoute("filter")

	case message.ReturnSelectionsMsg:
		if m.Props().Limit == 1 {
			return message.ToggleItem()
		}
		if m.Props().NumSelected == 0 {
			m.quitting = true
			return message.ToggleItem()
		}
		return message.ReturnSelections()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Key.Quit):
			m.quitting = true
			cmds = append(cmds, message.ReturnSelections())
		}
		for _, k := range m.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
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

	tm.list = list.New(props.Items)
	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(tm.Props().Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height
	return nil
}
