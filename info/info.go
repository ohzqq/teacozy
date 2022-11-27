package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Info struct {
	Model    viewport.Model
	hideKeys bool
	Visible  bool
	Editable bool
	Content  []string
	Sections []*Section
	Title    string
	Frame    style.Frame
	Style    Style
}

type Style struct {
	style.Field
	Title lipgloss.Style
}

func DefaultStyles() Style {
	s := Style{
		Field: style.DefaultFieldStyles(),
		Title: lipgloss.NewStyle().Foreground(style.Color.Pink),
	}
	return s
}

func New() *Info {
	info := &Info{
		Style: DefaultStyles(),
		Model: viewport.New(util.TermWidth(), util.TermHeight()),
	}
	return info
}

func (i *Info) NewSection() *Section {
	s := &Section{}
	i.Sections = append(i.Sections, s)
	return s
}

func (i Info) Render() string {
	var content []string
	for _, section := range i.Sections {
		content = append(content, section.Render(i.Style, i.hideKeys))
	}
	c := strings.Join(content, "\n")
	return c
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

func (i *Info) HideKeys() *Info {
	i.hideKeys = true
	return i
}

func (i *Info) Show() {
	i.Visible = true
}

func (i *Info) Hide() {
	i.Visible = false
}

func (i *Info) Toggle() {
	i.Visible = !i.Visible
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
		case key.Matches(msg, key.EditField):
			if m.Editable {
				cmds = append(cmds, EditInfoCmd())
			}
		case key.Matches(msg, key.InfoKey):
			m.Toggle()
		case key.Matches(msg, key.HelpKey):
			fallthrough
		case key.Matches(msg, key.ExitScreen):
			fallthrough
		case key.Matches(msg, key.PrevScreen):
			cmds = append(cmds, HideInfoCmd())
		}
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	case HideInfoMsg:
		m.Hide()
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Info) Init() tea.Cmd {
	m.Show()
	return nil
}

func (m *Info) View() string {
	m.Model.SetContent(m.Render())
	if m.Visible {
		return m.Model.View()
	}
	return ""
}
