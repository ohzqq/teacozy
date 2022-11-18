package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Info struct {
	Model     viewport.Model
	HideKeys  bool
	IsVisible bool
	Editable  bool
	id        int
	Style     FieldStyle
	Frame     Frame
	Data      FormData
	Fields    *Fields
}

func NewInfo() *Info {
	info := Info{
		Fields: NewFields(),
		Frame:  DefaultWidgetStyle(),
	}
	info.Model = viewport.New(info.Frame.Width(), info.Frame.Height())
	return &info
}

func (i *Info) SetData(data FormData) *Info {
	fields := NewFields().SetData(data)
	if f, ok := data.(*Fields); ok {
		fields = f
	}
	i.Data = data
	i.Fields = fields
	return i
}

func (i *Info) SetSize(w, h int) *Info {
	i.Model = viewport.New(w, h)
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
		case key.Matches(msg, Keys.ExitScreen):
			cmds = append(cmds, HideInfoCmd())
		case key.Matches(msg, Keys.EditField):
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

func (m *Info) Init() tea.Cmd {
	return nil
}

func (m *Info) View() string {
	m.Model.SetContent(m.Fields.String())
	return m.Model.View()
}
