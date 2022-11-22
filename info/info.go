package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keybind"
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
	Style     style.Field
}

func New(data teacozy.FormData) *Info {
	return &Info{
		Data:  data,
		Style: style.DefaultFieldStyles(),
	}
}

func (i Info) RenderData() string {
	var info []string
	for _, key := range i.Data.Keys() {
		fd := i.Data.Get(key)

		var line []string
		if !i.HideKeys {
			k := i.Style.Key.Render(fd.Key())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(fd.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}

	return strings.Join(info, "\n")
}

func (i *Info) SetHeight(h int) *Info {
	i.SetSize(i.Frame.Width(), h)
	return i
}

func (i *Info) SetSize(w, h int) *Info {
	i.Frame.SetSize(w, h)
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

//func (m *Info) Update(msg tea.Msg) (*Info, tea.Cmd) {
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
		switch {
		case key.Matches(msg, keybind.SaveAndExit):
			if m.Editable {
				cmds = append(cmds, EditInfoCmd())
			}
		case key.Matches(msg, keybind.HelpKey):
			fallthrough
		case key.Matches(msg, keybind.ExitScreen):
			fallthrough
		case key.Matches(msg, keybind.PrevScreen):
			cmds = append(cmds, HideInfoCmd())
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
	content := m.RenderData()
	if c := m.content; len(m.content) > 0 {
		content = strings.Join(c, "\n")
	}
	m.Model.SetContent(content)
	return m.Model.View()
}
