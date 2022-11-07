package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

var fieldStyle Style

type state int

const (
	view state = iota
	form
	edit
)

type Style struct {
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

type Model struct {
	state state
	*Info
	*Form
}

func New(data FormData) *Model {
	fieldStyle = Style{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}
	fields := NewFields().SetData(data)
	m := Model{
		Info: &Info{
			Fields: fields,
		},
		Form: &Form{
			Fields: fields,
		},
	}
	m.Display()
	m.Edit()
	return &m
}

type Info struct {
	Model    viewport.Model
	Fields   *Fields
	HideKeys bool
}

func (i *Info) Display() *Info {
	content := i.String()
	height := lipgloss.Height(content)
	i.Model = viewport.New(util.TermWidth(), height)
	i.Model.SetContent(content)
	return i
}

func (i *Info) NoKeys() *Info {
	i.HideKeys = true
	return i
}

func (i Info) String() string {
	var info []string
	for _, key := range i.Fields.Keys() {
		var line []string
		field := i.Fields.Get(key)
		if !i.HideKeys {
			k := fieldStyle.KeyStyle.Render(field.Key())
			line = append(line, k, ": ")
		}

		v := fieldStyle.ValueStyle.Render(field.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch m.state {
		case view:
			switch {
			case key.Matches(msg, urkey.EditField):
				cmds = append(cmds, UpdateInfoCmd())
			}
			m.Info.Model, cmd = m.Info.Model.Update(msg)
			cmds = append(cmds, cmd)
		case form:
			switch {
			case key.Matches(msg, urkey.ExitScreen):
				m.state = view
			}
		}
	case EditInfoMsg:
		m.Edit()
		m.state = form
	case tea.WindowSizeMsg:
		m.Info.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	switch m.state {
	case view:
	case form:
		m.Form, cmd = m.Form.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var (
		sections []string
		//availHeight = m.form.List.Height()
		availHeight = util.TermHeight()
	)

	switch m.state {
	case view:
		m.Info.Model.SetContent(m.String())
		v := m.Info.Model.View()
		sections = append(sections, v)
	case form:
		v := m.Form.View()
		sections = append(sections, v)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
