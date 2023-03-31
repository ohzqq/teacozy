package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
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
	Props
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
		Style: DefaultStyle(),
	}
	return &tm
}

func ChooseRouteInitializer(c *List) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewChoice()
		props := ChooseProps{
			Props:      c.NewProps(),
			ToggleItem: c.ToggleSelection,
		}
		return component, component.Init(props)
	}
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
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
	case StartEditingMsg:
		reactea.SetCurrentRoute("form")
		return nil
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
			cmds = append(cmds, DownCmd())
		case key.Matches(msg, chooseKey.Edit):
			cmds = append(cmds, StartEditingCmd())
		case key.Matches(msg, chooseKey.Filter):
			cmds = append(cmds, StartFilteringCmd())
		case key.Matches(msg, chooseKey.Bottom):
			m.Cursor = len(m.Props().Visible()) - 1
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
		m.Props().Footer(m.Paginator.View())
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
