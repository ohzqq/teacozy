package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

var fieldStyle Style

type Style struct {
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

type Info struct {
	Model    viewport.Model
	Fields   *Fields
	HideKeys bool
}

func NewInfo(data FormData) *Info {
	fieldStyle = Style{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}

	fields := NewFields().SetData(data)
	info := Info{Fields: fields}
	return &info
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
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
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

func (m *Info) View() string {
	m.Model.SetContent(m.String())
	return m.Model.View()
}

func (i *Info) Init() tea.Cmd {
	return nil
}
