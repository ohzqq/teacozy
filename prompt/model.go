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

type Model struct {
	List             list.Model
	Title            string
	MultiSelect      bool
	ShowSelectedOnly bool
	Keys             urkey.KeyMap
	Items            item.Items
	width            int
	height           int
	Style            list.Styles
}

func New() Model {
	w, h := util.TermSize()
	p := Model{
		Items:  item.NewItems(),
		width:  w,
		height: h,
		Keys:   urkey.DefaultKeys(),
	}
	return p
}

func (m *Model) Start() *Model {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *Model) SetItems(items item.Items) *Model {
	m.Items = items
	return m
}

func (m *Model) SetMultiSelect() *Model {
	m.MultiSelect = true
	m.Items.SetMultiSelect()
	return m
}

func (m *Model) SetSize(w, h int) *Model {
	m.width = w
	m.height = h
	return m
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.ExitScreen):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, m.Keys.Prev):
			m.ShowSelectedOnly = false
			cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
		}
		if m.MultiSelect {
			switch {
			case key.Matches(msg, urkey.Enter):
				if m.ShowSelectedOnly {
					cmds = append(cmds, ReturnSelectionsCmd())
				}
				m.ShowSelectedOnly = true
				cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
			case key.Matches(msg, m.Keys.SelectAll):
				m.Items.ToggleAllSelectedItems()
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
		} else {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cur := m.List.SelectedItem().(*item.Item)
				m.Items.ToggleSelectedItem(cur.Index())
				cmds = append(cmds, ReturnSelectionsCmd())
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
	case ReturnSelectionsMsg:
		cmds = append(cmds, tea.Quit)
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

func (m *Model) Init() tea.Cmd {
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

func (m Model) View() string {
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
