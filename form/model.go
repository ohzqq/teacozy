package form

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
)

type Info struct {
	*Form
}

func NewForm(data info.FormData) *Info {
	return &Info{Form: New(data)}
}

func NewInfo(data info.FormData) *Info {
	form := New(data)
	form.state = view
	return &Info{
		Form: form,
	}
}

func (m *Info) Start() *info.Fields {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m.Fields
}

func (m *Info) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Info) Init() tea.Cmd {
	return nil
}

func (m *Info) View() string {
	return m.Form.View()
}
