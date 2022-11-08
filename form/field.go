package form

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/slices"
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
	Set(string, string)
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

type DefaultFields struct {
	data []Field
}

func NewFields() *Fields {
	return &Fields{
		Style: fieldStyle,
	}
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

func (f DefaultFields) Get(key string) Field {
	for _, field := range f.data {
		if field.Key() == key {
			return field
		}
	}
	return &DefaultField{}
}

func (f *DefaultFields) Set(key, val string) {
	if f.Has(key) {
		ff := f.Get(key)
		ff.Set(val)
	} else {
		f.Add(key, val)
	}
}

func (f DefaultFields) Keys() []string {
	var keys []string
	for _, field := range f.data {
		keys = append(keys, field.Key())
	}
	return keys
}

func (f DefaultFields) Has(key string) bool {
	return slices.Contains(f.Keys(), key)
}

func (f *DefaultFields) Add(key, val string) error {
	if f.Has(key) {
		return fmt.Errorf("keys must be unique")
	}
	field := NewDefaultField(key, val)
	f.data = append(f.data, field)
	return nil
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
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (i Fields) String() string {
	var info []string
	for _, key := range i.Data.Keys() {
		var line []string
		field := i.Data.Get(key)
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

type DefaultField struct {
	key   string
	value string
}

func NewDefaultField(key, val string) *DefaultField {
	return &DefaultField{
		key:   key,
		value: val,
	}
}

func (f DefaultField) FilterValue() string {
	return f.value
}

func (f DefaultField) Value() string {
	return f.value
}

func (f *DefaultField) Set(val string) {
	f.value = val
}

func (f DefaultField) Key() string {
	return f.key
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
