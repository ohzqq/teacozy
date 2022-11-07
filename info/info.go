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

type Style struct {
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

type Model struct {
	view viewport.Model
	*Info
}

func New(data FormData) *Model {
	fieldStyle = Style{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}
	fields := NewFields().SetData(data)
	m := Model{
		Info: &Info{
			Data: fields,
		},
	}
	height := lipgloss.Height(m.String())
	m.view = viewport.New(util.TermWidth(), height)
	m.view.SetContent(m.String())
	return &m
}

type FormData interface {
	Get(string) string
	Set(string, string)
	Keys() []string
}

type Info struct {
	Data     *Fields
	HideKeys bool
}

func (i *Info) NoKeys() *Info {
	i.HideKeys = true
	return i
}

func (i Info) String() string {
	var info []string
	for _, key := range i.Data.Keys() {
		var line []string
		if !i.HideKeys {
			k := fieldStyle.KeyStyle.Render(key)
			line = append(line, k, ": ")
		}

		val := i.Data.Get(key)
		v := fieldStyle.ValueStyle.Render(val)
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (i *Info) Set(f ...map[string]string) *Info {
	var fields Fields
	for _, field := range f {
		for k, v := range field {
			fields.data = append(fields.data, NewField(k, v))
		}
	}
	i.Data = &fields
	return i
}

type UpdateContentMsg struct {
	Field
}

func UpdateContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Field: NewField(key, val)}
	}
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
		switch {
		case key.Matches(msg, urkey.EditField):
			cmds = append(cmds, UpdateContentCmd("one", "edit"))
		}
	case EditInfoMsg:
	case UpdateContentMsg:
		m.Data.Set(msg.Key, msg.Value)
	case tea.WindowSizeMsg:
		m.view = viewport.New(msg.Width-2, msg.Height-2)
	}
	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

type EditInfoMsg struct{}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	m.view.SetContent(m.String())
	return m.view.View()
}
