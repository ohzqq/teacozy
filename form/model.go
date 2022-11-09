package form

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
)

type Model struct {
	*Form
}

func NewForm(data info.FormData) *Model {
	return &Model{Form: New(data)}
}

func NewInfo(data info.FormData) *Model {
	form := New(data)
	form.state = view
	return &Model{
		Form: form,
	}
}

func (m *Model) Start() *info.Fields {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m.Fields
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
	}

	m.Form, cmd = m.Form.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return m.Form.View()
}
