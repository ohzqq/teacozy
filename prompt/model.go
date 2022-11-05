package prompt

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list/item"
	"github.com/ohzqq/teacozy/util"
)

type Prompt struct {
	List        list.Model
	Title       string
	MultiSelect bool
	Keys        urkey.KeyMap
	Items       item.Items
	width       int
	height      int
	Style       list.Styles
}

func New(title string) Prompt {
	w, h := util.TermSize()
	p := Prompt{
		Title:  title,
		Items:  item.NewItems(),
		width:  w,
		height: h,
		Keys:   urkey.DefaultKeys(),
	}
	return p
}

func (m *Prompt) Start() *Prompt {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *Prompt) SetItems(items item.Items) *Prompt {
	m.Items = items
	return m
}

func (m *Prompt) SetSize(w, h int) *Prompt {
	m.width = w
	m.height = h
	//if m.List != nil {
	//  m.List.SetSize(w, h)
	//}
	return m
}

func (m *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.Keys.Quit) {
			cmds = append(cmds, tea.Quit)
		}
		if m.MultiSelect {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cmds = append(cmds, UpdateStatusCmd("multi"))
				//if m.ShowSelectedOnly {
				//  cmds = append(cmds, ReturnSelectionsCmd())
				//}
				//m.ShowSelectedOnly = true
				//cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
				//case key.Matches(msg, m.Keys.SelectAll):
				//ToggleAllItemsCmd(m)
				//cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			}
		} else {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cmds = append(cmds, UpdateStatusCmd("single"))
				//cur := m.Model.SelectedItem().(Item)
				//m.SetItem(m.Model.Index(), cur.ToggleSelected())
				//cmds = append(cmds, ReturnSelectionsCmd())
			}
		}
	case UpdateStatusMsg:
		cmds = append(cmds, m.List.NewStatusMessage(msg.Msg))
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Prompt) Init() tea.Cmd {
	l := m.Items.List()
	l.SetSize(m.width, m.height)
	m.List = l
	return nil
}

func (m Prompt) View() string {
	return m.List.View()
}
