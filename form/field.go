package form

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

var fieldStyle = Style{
	Key:   lipgloss.NewStyle().Foreground(style.DefaultColors().Blue),
	Value: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
}

type Style struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

type FormData interface {
	Get(string) Field
	Keys() []string
}

type Field interface {
	FilterValue() string
	Value() string
	Key() string
	Set(string)
}

type Fields struct {
	Model    viewport.Model
	HideKeys bool
	Style    Style
	Data     FormData
}

func NewFields() *Fields {
	return &Fields{
		Style: fieldStyle,
	}
}

func (f *Fields) Edit() *Model {
	return NewForm(f.Data)
}

func (f *Fields) Render() *Fields {
	content := f.String()
	height := lipgloss.Height(content)
	f.Model = viewport.New(util.TermWidth(), height)
	f.Model.SetContent(content)
	return f
}

func (f *Fields) SetData(data FormData) *Fields {
	f.Data = data
	return f
}

func (f Fields) AllFields() []Field {
	var fields []Field
	for _, key := range f.Data.Keys() {
		field := f.Data.Get(key)
		fields = append(fields, field)
	}
	return fields
}

func (f *Fields) NoKeys() *Fields {
	f.HideKeys = true
	return f
}

func (m *Fields) Update(msg tea.Msg) (*Fields, tea.Cmd) {
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
		case key.Matches(msg, urkey.SaveAndExit):
			cmds = append(cmds, SaveAsHashCmd())
		case key.Matches(msg, urkey.EditField):
			cmds = append(cmds, EditInfoCmd())
		}
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	case EditInfoMsg:
		m.Render()
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (i Fields) String() string {
	var info []string
	for _, field := range i.AllFields() {
		var line []string
		if !i.HideKeys {
			k := i.Style.Key.Render(field.Key())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(field.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (m *Fields) View() string {
	m.Model.SetContent(m.String())
	return m.Model.View()
}

func (i *Fields) Init() tea.Cmd {
	return nil
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
