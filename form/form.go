package form

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/util"
)

type Form struct {
	List   list.Model
	Title  string
	Keys   urkey.KeyMap
	fields []Field
}

func New(title string, fields ...Field) Form {
	w, h := util.TermSize()
	p := Form{
		Title:  title,
		Keys:   urkey.DefaultKeys(),
		fields: fields,
	}
	p.List = list.New(p.items(), NewDelegate(), w, h)
	return p
}

func (m Form) items() []list.Item {
	var items []list.Item
	for _, f := range m.fields {
		items = append(items, f)
	}
	return items
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) View() string {
	return m.List.View()
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		//switch {
		//}
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

type Field struct {
	Label   string
	Content string
}

func NewField(label, content string) Field {
	return Field{
		Label:   label,
		Content: content,
	}
}

func (f Field) FilterValue() string {
	return f.Label
}

func (f Field) Title() string {
	return f.Label
}

func (f Field) Description() string {
	return f.Content
}
