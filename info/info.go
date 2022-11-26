package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

type Info struct {
	Model    viewport.Model
	HideKeys bool
	Visible  bool
	Editable bool
	Content  []string
	Title    string
	Data     teacozy.Fields
	Frame    style.Frame
	Style    Style
}

type Style struct {
	style.Field
	Title lipgloss.Style
}

func New(data teacozy.Fields) *Info {
	s := Style{
		Field: style.DefaultFieldStyles(),
		Title: lipgloss.NewStyle().Foreground(style.Color.Pink),
	}
	info := &Info{
		Style: s,
	}
	info.AddFields(data)
	return info
}

func (i *Info) AddFields(data teacozy.Fields) *Info {
	i.AddContent(i.RenderData(data)...)
	return i
}

func (i *Info) SetTitle(title string) *Info {
	i.Title = title
	return i
}

func (i Info) RenderData(data teacozy.Fields) []string {
	var info []string

	if i.Title != "" {
		t := i.Style.Title.Render(i.Title)
		info = append(info, t)
	}

	for _, key := range data.Keys() {
		fd := data.Get(key)

		var line []string
		if !i.HideKeys {
			k := i.Style.Key.Render(fd.Name())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(fd.Content())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}

	return info
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
	i.Content = []string{content}
	return i
}

func (i *Info) AddContent(content ...string) *Info {
	i.Content = append(i.Content, content...)
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
	if m.Visible {
		m.Model.SetContent(strings.Join(m.Content, "\n"))
		return m.Model.View()
	}
	return ""
}
