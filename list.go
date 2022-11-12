package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	List             list.Model
	Title            string
	MultiSelect      bool
	ShowKeys         bool
	ShowSelectedOnly bool
	Keys             KeyMap
	Items            Items
	Width            int
	Height           int
	Style            list.Styles
}

func NewList() *List {
	w, h := TermSize()
	p := List{
		Items:  NewItems(),
		Width:  w,
		Height: h,
		Keys:   DefaultKeys(),
	}
	return &p
}

func (m *List) InitList() list.Model {
	l := m.Items.List()
	l.SetSize(m.Width, m.Height)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()
	l.Title = m.Title

	if m.Title == "" {
		l.Title = ""
	}
	m.List = l
	return l
}

func (m *List) SetItems(items Items) *List {
	m.Items = items
	return m
}

func (m *List) SetMultiSelect() *List {
	m.MultiSelect = true
	m.Items.SetMultiSelect()
	return m
}

func (m *List) SetShowKeys() *List {
	m.ShowKeys = true
	m.Items.SetShowKeys()
	return m
}

func (m *List) SetSize(w, h int) *List {
	m.Width = w
	m.Height = h
	m.List.SetSize(w, h)
	return m
}

func (m *List) Update(msg tea.Msg) (*List, tea.Cmd) {
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
		case key.Matches(msg, m.Keys.Prev):
			m.ShowSelectedOnly = false
			cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
		}
		if m.MultiSelect {
			switch {
			case key.Matches(msg, Enter):
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
				cur := m.List.SelectedItem().(*Item)
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
	case ToggleSelectedMsg:
		m.Items.ToggleSelectedItem(msg.Index())
	case ReturnSelectionsMsg:
		cmds = append(cmds, tea.Quit)
	case ToggleListMsg:
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

func (m *List) Init() tea.Cmd {
	m.List = m.InitList()
	return nil
}

func (m List) View() string {
	return m.List.View()
}