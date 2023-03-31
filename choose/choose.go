package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ChooseProps]
	Cursor    int
	Viewport  *viewport.Model
	Paginator paginator.Model
	quitting  bool
	header    string
	Style     style.List
}

type ChooseProps struct {
	*props.Items
	ToggleItem func(int)
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
	Edit             key.Binding
}

func NewChoice() *Choose {
	tm := Choose{
		Style: style.ListDefaults(),
	}
	return &tm
}

func ChooseRouteInitializer(props ChooseProps) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewChoice()
		//props := ChooseProps{
		//  Props:      c.NewProps(),
		//  ToggleItem: c.ToggleSelection,
		//}
		return component, component.Init(props)
	}
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case message.UpMsg:
		m.Cursor--
		if m.Cursor < 0 {
			m.Cursor = len(m.Props().Visible()) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		}
		if m.Cursor < start {
			m.Paginator.PrevPage()
		}
	case message.DownMsg:
		m.Cursor++
		if m.Cursor >= len(m.Props().Visible()) {
			m.Cursor = 0
			m.Paginator.Page = 0
		}
		if m.Cursor >= end {
			m.Paginator.NextPage()
		}
	case message.StartEditingMsg:
		reactea.SetCurrentRoute("form")
		return nil
	case message.StartFilteringMsg:
		reactea.SetCurrentRoute("filter")
		return nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, chooseKey.Up):
			cmds = append(cmds, message.UpCmd())
		case key.Matches(msg, chooseKey.Down):
			cmds = append(cmds, message.DownCmd())
		case key.Matches(msg, chooseKey.Prev):
			m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor-m.Props().Height)
			m.Paginator.PrevPage()
		case key.Matches(msg, chooseKey.Next):
			m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor+m.Props().Height)
			m.Paginator.NextPage()
		case key.Matches(msg, chooseKey.ToggleItem):
			if m.Props().Limit == 1 {
				return nil
			}
			idx := m.Props().Visible()[m.Cursor].Index
			m.Props().ToggleItem(idx)
			cmds = append(cmds, message.DownCmd())
		case key.Matches(msg, chooseKey.Edit):
			cmds = append(cmds, message.StartEditingCmd())
		case key.Matches(msg, chooseKey.Filter):
			cmds = append(cmds, message.StartFilteringCmd())
		case key.Matches(msg, chooseKey.Bottom):
			m.Cursor = len(m.Props().Visible()) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		case key.Matches(msg, chooseKey.Top):
			m.Cursor = 0
			m.Paginator.Page = 0
		case key.Matches(msg, chooseKey.Quit):
			m.quitting = true
			cmds = append(cmds, message.ReturnSelectionsCmd())
		case key.Matches(msg, chooseKey.ReturnSelections):
			cmds = append(cmds, message.ReturnSelectionsCmd())
		}
	}

	return tea.Batch(cmds...)
}

func (m *Choose) Render(w, h int) string {
	m.Viewport.Height = h
	if m.Paginator.TotalPages > 1 {
		m.Viewport.Height = m.Viewport.Height + 4
	}
	m.Viewport.Width = w

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
		//m.Props().Footer(m.Paginator.View())
	}

	view = s.String()

	m.Viewport.SetContent(view)
	view = m.Viewport.View()

	return view
}

func (tm *Choose) Init(props ChooseProps) tea.Cmd {
	tm.UpdateProps(props)
	v := viewport.New(0, 0)
	tm.Viewport = &v
	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(tm.Props().Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height
	return nil
}

var chooseKey = ChooseKeys{
	Next: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("right/l", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("left/h", "prev page"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit form"),
	),
	//key.NewBinding(
	//  key.WithKeys("V"),
	//  key.WithHelp("V", "deselect all"),
	//),
	//key.NewBinding(
	//  key.WithKeys("v"),
	//  key.WithHelp("v", "select all"),
	//),
	ToggleItem: key.NewBinding(
		key.WithKeys(" ", "tab"),
		key.WithHelp("space", "select item"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "move cursor down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "move cursor up"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "q", "ctrl+c"),
		key.WithHelp("esc/q", "quit"),
	),
	ReturnSelections: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "return selections"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter items"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "last item"),
	),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "first item"),
	),
}

func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
