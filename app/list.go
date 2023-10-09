package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/list"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	list *list.Model
}

type Props struct {
	SetCurrentItem func(*list.Item)
	Items          *list.Items
	Opts           []list.Option
}

func NewList() *List {
	return &List{
		//list: list.New(items, opts...),
	}
}

func (l *List) Init(props Props) tea.Cmd {
	l.UpdateProps(props)
	l.list = list.New(l.Props().Items, l.Props().Opts...)
	//return l.list.NewStatusMessage(fmt.Sprint(l.list.Editable()))
	return nil
}

func (l *List) Update(msg tea.Msg) tea.Cmd {
	m, cmd := l.list.Update(msg)
	l.list = m.(*list.Model)
	l.Props().SetCurrentItem(l.list.CurrentItem())
	return cmd
}

func (l *List) Render(w, h int) string {
	if l.list.HasInput() {
		h--
	}
	l.list.SetSize(w, h)
	return l.list.View()
}
