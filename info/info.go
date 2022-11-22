package info

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/style"
)

type Info struct {
	Model     viewport.Model
	HideKeys  bool
	IsVisible bool
	Editable  bool
	content   []string
	Frame     style.Frame
	Data      teacozy.FormData
	//Style     FieldStyle
}

func New(data teacozy.FormData) *Info {
	return &Info{
		Data: data,
	}
}

func (i *Info) SetHeight(h int) *Info {
	i.Model = viewport.New(i.Frame.Width(), h)
	return i
}

func (i *Info) SetSize(w, h int) *Info {
	i.Model = viewport.New(w, h)
	return i
}

func (i *Info) SetContent(content string) *Info {
	i.content = append(i.content, content)
	return i
}

func (i *Info) AddContent(content ...string) *Info {
	i.content = append(i.content, content...)
	return i
}

func (m *Info) Update(msg tea.Msg) (*Info, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case Keys.Help.Matches(msg):
			cmds = append(cmds, HideInfoCmd())
		case Keys.ExitScreen.Matches(msg):
			cmds = append(cmds, HideInfoCmd())
		case Keys.PrevScreen.Matches(msg):
			cmds = append(cmds, HideInfoCmd())
		case Keys.EditField.Matches(msg):
			if m.Editable {
				cmds = append(cmds, EditInfoCmd(m.Fields))
			}
		}
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
