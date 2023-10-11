package pager

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/bubbles/viewport"
	"github.com/ohzqq/teacozy/util"
)

type Model struct {
	viewport.Model
	render   Renderer
	text     string
	rendered string
	width    int
	height   int
	focused  bool
	KeyMap   KeyMap
}

type Renderer func(text string, width int) string

func New(fn Renderer) *Model {
	w, h := util.TermSize()
	m := &Model{
		render: fn,
		width:  w,
		height: h,
		Model:  viewport.New(w, h),
		KeyMap: DefaultKeyMap(),
	}
	m.Model.KeyMap = m.KeyMap.KeyMap

	return m
}

func (m *Model) SetText(text string) *Model {
	m.text = text
	return m
}

func RenderText(text string, width int) string {
	s := lipgloss.NewStyle().Width(width)
	return s.Render(text)
}

func RenderHTML(html string, width int) string {
	conv := md.NewConverter("", true, nil)
	mark, err := conv.ConvertString(html)
	if err != nil {
		return html
	}
	return RenderMarkdown(mark, width)
}

func RenderMarkdown(mark string, width int) string {
	gr, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return mark
	}

	r, err := gr.Render(mark)
	if err != nil {
		return mark
	}

	//println(r)
	return r
}

type FocusMsg struct{}
type UnfocusMsg struct{}

func (m *Model) Focus() tea.Cmd {
	return func() tea.Msg {
		m.focused = true
		m.KeyMap.PageDown.SetEnabled(true)
		m.KeyMap.PageUp.SetEnabled(true)
		m.KeyMap.HalfPageDown.SetEnabled(true)
		m.KeyMap.HalfPageUp.SetEnabled(true)
		m.KeyMap.Down.SetEnabled(true)
		m.KeyMap.Up.SetEnabled(true)
		return FocusMsg{}
	}
}

func (m *Model) Unfocus() tea.Cmd {
	return func() tea.Msg {
		m.focused = false
		m.KeyMap.PageDown.SetEnabled(false)
		m.KeyMap.PageUp.SetEnabled(false)
		m.KeyMap.HalfPageDown.SetEnabled(false)
		m.KeyMap.HalfPageUp.SetEnabled(false)
		m.KeyMap.Down.SetEnabled(false)
		m.KeyMap.Up.SetEnabled(false)
		return UnfocusMsg{}
	}
}

func (m Model) Focused() bool {
	return m.focused
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.width == 0 {
			m.width = msg.Width
		}
		if m.height == 0 {
			m.height = msg.Height
		}
		if msg.Width < m.width {
			m.Model.Width = msg.Width
		}
		if msg.Height < m.height {
			m.Model.Height = msg.Height
		}
		if msg.Width > m.width {
			m.Model.Width = m.width
		}
		if msg.Height > m.height {
			m.Model.Height = m.height
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		}
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	m.Model.SetContent(m.Render())

	return m, tea.Batch(cmds...)
}

func (m Model) Render() string {
	return m.render(m.text, m.Model.Width)
}

func (m Model) View() string {
	m.Model.SetContent(m.Render())
	return m.Model.View()
}

func (m *Model) SetSize(w, h int) *Model {
	m.width = w
	m.height = h
	m.Model = viewport.New(w, h)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

type KeyMap struct {
	viewport.KeyMap
	Quit key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		KeyMap: defaultKeyMap(),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

func defaultKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdn", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "½ page down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
	}
}
