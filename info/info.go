package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keybind"
	"github.com/ohzqq/teacozy/style"
)

type Info struct {
	Model    viewport.Model
	HideKeys bool
	visible  bool
	Editable bool
	content  []string
	title    string
	Frame    style.Frame
	Data     teacozy.FormData
	Style    Style
}

type Style struct {
	style.Field
	Title lipgloss.Style
}

func New(data teacozy.FormData) *Info {
	s := Style{
		Field: style.DefaultFieldStyles(),
		Title: lipgloss.NewStyle().Foreground(style.Color.Pink),
	}
	return &Info{
		Data:  data,
		Style: s,
	}
}

func (i *Info) SetData(data teacozy.FormData) *Info {
	i.Data = data
	return i
}

func (i *Info) SetTitle(title string) *Info {
	i.title = title
	return i
}

func (i Info) RenderData() string {
	var info []string

	if i.title != "" {
		t := i.Style.Title.Render(i.title)
		info = append(info, t)
	}

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

func (i *Info) Show() {
	i.visible = true
}

func (i *Info) Hide() {
	i.visible = false
}

func (i *Info) Toggle() {
	i.visible = !i.visible
}

func (i Info) Visible() bool {
	return i.visible
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
		case key.Matches(msg, keybind.EditField):
			if m.Editable {
				cmds = append(cmds, EditInfoCmd())
			}
		case key.Matches(msg, keybind.InfoKey):
			m.Toggle()
		case key.Matches(msg, keybind.HelpKey):
			fallthrough
		case key.Matches(msg, keybind.ExitScreen):
			fallthrough
		case key.Matches(msg, keybind.PrevScreen):
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
	if m.Visible() {
		content := m.RenderData()
		if c := m.content; len(m.content) > 0 {
			content = strings.Join(c, "\n")
		}
		m.Model.SetContent(content)
		return m.Model.View()
	}
	return ""
}
