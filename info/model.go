package info

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/util"
)

type state int

const (
	form state = iota
	view
	edit
)

type Model struct {
	state state
	*Fields
	*Form
}

func New(data FormData) *Form {
	fields := NewFields().SetData(data)
	m := Form{
		Fields: fields,
		//Form: &Form{
		//  Fields: fields,
		//},
	}
	m.Fields.Render()
	m.Edit()
	return &m
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
				cmds = append(cmds, EditInfoCmd())
			}
		case form:
			switch {
			case key.Matches(msg, urkey.ExitScreen):
				m.state = view
			}
		}
	case EditInfoMsg:
		m.Edit()
		m.state = form
	}

	switch m.state {
	case view:
		m.Fields, cmd = m.Fields.Update(msg)
		cmds = append(cmds, cmd)
	case form:
		//m.Form, cmd = m.Form.Update(msg)
		//cmds = append(cmds, cmd)
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
		v := m.Fields.View()
		sections = append(sections, v)
	case form:
		v := m.Form.View()
		sections = append(sections, v)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
