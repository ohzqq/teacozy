package prompt

import (
	"fmt"
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

func New() Prompt {
	w, h := util.TermSize()
	p := Prompt{
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

func (m *Prompt) SetMultiSelect() *Prompt {
	m.MultiSelect = true
	m.Items.SetMultiSelect()
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
			case key.Matches(msg, urkey.Enter):
				t := fmt.Sprintf("%v", m.MultiSelect)
				cmds = append(cmds, UpdateStatusCmd(t))
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
				t := fmt.Sprintf("%v", m.MultiSelect)
				cmds = append(cmds, UpdateStatusCmd(t))
				//cur := m.Model.SelectedItem().(Item)
				//m.SetItem(m.Model.Index(), cur.ToggleSelected())
				//cmds = append(cmds, ReturnSelectionsCmd())
			}
		}
	case UpdateStatusMsg:
		cmds = append(cmds, m.List.NewStatusMessage(msg.Msg))
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width-1, msg.Height-2)
	case UpdateVisibleItemsMsg:
		items := m.Items.Display(string(msg))
		m.List.SetItems(items)
	case item.ToggleSelectedMsg:
		m.Items.ToggleSelectedItem(msg.Index())
	case item.ToggleListMsg:
		switch msg.ListOpen {
		case true:
			m.Items.CloseItemList(msg.Index())
		default:
			m.Items.OpenItemList(msg.Index())
		}
		cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Prompt) Init() tea.Cmd {
	l := m.Items.List()
	l.SetSize(m.width, m.height)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()

	if m.Title == "" {
		l.Title = ""
		//l.SetShowTitle(true)
	}
	m.List = l
	return nil
}

func (m Prompt) View() string {
	return m.List.View()
}

func ListKeyMap() list.KeyMap {
	km := list.DefaultKeyMap()
	km.NextPage = key.NewBinding(
		key.WithKeys("right", "l", "pgdown"),
		key.WithHelp("l/pgdn", "next page"),
	)
	km.Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	)
	return km
}
