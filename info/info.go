package info

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

type infoStr string

func (i infoStr) String() string {
	return string(i)
}

type Model struct {
	view viewport.Model
	*Info
}

func New() *Model {
	return &Model{
		view: viewport.New(1, 1),
		Info: &Info{
			KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
			ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		},
	}
}

type FormData interface {
	Get(string) string
	Set(string, string)
	Keys() []string
}

type Info struct {
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
	content    []map[fmt.Stringer]fmt.Stringer
	HideKeys   bool
}

func (m *Model) Init() tea.Cmd {
	m.view.SetContent(m.String())
	//return UpdateContentCmd(m.String())
	return nil
}

func (m *Model) View() string {
	return m.view.View()
}

type UpdateContentMsg struct {
	Content string
}

func UpdateContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Content: content}
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
			cmds = append(cmds, UpdateContentCmd("edit"))
		}
	case UpdateContentMsg:
		m.view.SetContent(msg.Content)
	case tea.WindowSizeMsg:
		m.view = viewport.New(msg.Width-2, msg.Height-2)
	}
	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (i *Info) NoKeys() *Info {
	i.HideKeys = true
	return i
}

func (i *Info) AddString(key, val string) {
	i.Add(infoStr(key), infoStr(val))
}

func (i *Info) Add(key, val fmt.Stringer) {
	content := make(map[fmt.Stringer]fmt.Stringer)
	content[key] = val
	i.content = append(i.content, content)
}

func (i *Info) Set(content ...map[fmt.Stringer]fmt.Stringer) *Info {
	i.content = content
	return i
}

func (i Info) String() string {
	var info []string
	for _, pair := range i.content {
		var line []string
		for key, val := range pair {
			if !i.HideKeys {
				if str := key.String(); str != "" {
					k := i.KeyStyle.Render(str)
					line = append(line, k, ": ")
				}
			}
			if str := val.String(); str != "" {
				v := i.KeyStyle.Render(str)
				line = append(line, v)
			}
		}
		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}
