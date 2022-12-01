package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

type Info struct {
	Model    viewport.Model
	hideKeys bool
	Visible  bool
	Editable bool
	content  string
	Sections []*Section
	Title    string
	Style    Style
	Toggle   *key.Key
}

type Style struct {
	style.Field
	style.Frame
	Title lipgloss.Style
}

func DefaultStyles() Style {
	s := Style{
		Field: style.DefaultFieldStyles(),
		Title: lipgloss.NewStyle().Foreground(style.Color.Pink),
		Frame: style.DefaultFrameStyle(),
	}
	return s
}

func New() *Info {
	info := &Info{
		Style:   DefaultStyles(),
		Toggle:  key.NewKey("i", "info"),
		Visible: false,
	}
	info.Model = viewport.New(info.Style.Width(), info.Style.Height())
	return info
}

func (m *Info) SetToggle(toggle, help string) *Info {
	m.Toggle = key.NewKey(toggle, help)
	return m
}

func (i *Info) SetContent(c string) *Info {
	i.content = c
	return i
}

func (i *Info) NewSection() *Section {
	s := &Section{}
	i.Sections = append(i.Sections, s)
	return s
}

func (i *Info) Render() string {
	var sections []string
	for _, section := range i.Sections {
		sections = append(sections, section.Render(i.Style, i.hideKeys))
	}
	content := strings.Join(sections, "\n")

	if i.content != "" {
		content = i.content
	}

	return content
}

func (i *Info) SetHeight(h int) *Info {
	i.SetSize(i.Style.Width(), h)
	return i
}

func (i *Info) SetSize(w, h int) *Info {
	i.Style.SetSize(w, h)
	i.Model = viewport.New(i.Style.Width(), i.Style.Height())
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

func (i *Info) ToggleVisible() {
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
		case m.Toggle.Matches(msg):
			cmds = append(cmds, HideInfoCmd())
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width-2, msg.Height-2)
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
	c := m.Model.View()
	return c
}
