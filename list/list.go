package list

import (
	"bytes"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	Model            list.Model
	Items            Items
	items            []Item
	all              []list.Item
	Keys             urkey.KeyMap
	IsPrompt         bool
	IsMultiSelect    bool
	ShowSelectedOnly bool
}

func NewList() List {
	return List{
		Keys: urkey.DefaultKeys(),
	}
}

func (l *List) NewItem(content string) Item {
	i := Item{Content: content}
	l.items = append(l.items, i)
	return i
}

func (l *List) ProcessItems() *List {
	l.all = FlattenItems(l.items)
	for idx, item := range l.all {
		i := item.(Item)
		if l.IsMulti() {
			i.IsMulti = true
		}
		i.id = idx
		l.all[idx] = i
	}
	return l
}

func (l List) AllItems() []list.Item {
	l.ProcessItems()
	return l.all
}

func (l *List) AppendItem(item Item) *List {
	//item := NewItem(i)
	if l.IsMulti() {
		item.IsMulti = true
	}
	l.Items = l.Items.Add(item)
	return l
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.IsMulti() {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				if m.ShowSelectedOnly {
					cmds = append(cmds, ReturnSelectionsCmd())
				}
				m.ShowSelectedOnly = true
				cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
			case key.Matches(msg, m.Keys.SelectAll):
				//ToggleAllItemsCmd(m)
				cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			}
		} else {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cur := m.Model.SelectedItem().(Item)
				m.SetItem(m.Model.Index(), cur.ToggleSelected())
				cmds = append(cmds, ReturnSelectionsCmd())
			}
		}

		switch {
		case key.Matches(msg, m.Keys.ExitScreen):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, m.Keys.Prev):
			m.ShowSelectedOnly = false
			cmds = append(cmds, UpdateVisibleItemsCmd("all"))
		}
	case UpdateStatusMsg:
		cmds = append(cmds, m.Model.NewStatusMessage(string(msg)))
	case UpdateVisibleItemsMsg:
		items := m.Items.Display(string(msg))
		//m.Model.SetHeight(m.GetHeight(items))
		m.Model.SetHeight(util.TermHeight() - 2)
		cmds = append(cmds, m.Model.SetItems(items))
	case ToggleItemListMsg:
		cur := m.Items.Get(int(msg))
		m.SetItem(m.Model.Index(), cur.ToggleList())
		cmds = append(cmds, UpdateVisibleItemsCmd("all"))
	case ToggleSelectedItemMsg:
		cur := m.Items.Get(int(msg))
		m.SetItem(m.Model.Index(), cur.ToggleSelected())
		cmds = append(cmds, UpdateVisibleItemsCmd("all"))
	case SetSizeMsg:
		if size := []int(msg); len(size) == 2 {
			m.Model.SetSize(size[0], size[1])
		}
	case SetItemsMsg:
		m.SetItems(Items(msg))
		m.processAllItems()
		cmds = append(cmds, UpdateVisibleItemsCmd("all"))
	case OSExecCmdMsg:
		menuCmd := msg.cmd(m.Items.Selected())
		var (
			stderr bytes.Buffer
			stdout bytes.Buffer
		)
		menuCmd.Stderr = &stderr
		menuCmd.Stdout = &stdout
		err := menuCmd.Run()
		if err != nil {
			fmt.Println(menuCmd.String())
			fmt.Println(stderr.String())
			log.Fatal(err)
		}
	}
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (l List) Init() tea.Cmd {
	return SetItemsCmd(l.Items)
}

func (l List) View() string {
	return l.Model.View()
}

func (l List) IsMulti() bool {
	return l.IsMultiSelect
}

func (l *List) SetItems(items Items) *List {
	l.Items = items
	return l
}

func (m *List) SetItem(modelIndex int, item Item) {
	m.Model.SetItem(modelIndex, item)
	m.Items.Set(item)
}

func (l *List) processAllItems() Items {
	var items Items
	for _, i := range l.Items {
		item := i.(Item)
		if l.IsMulti() {
			item.IsMulti = true
		}
		items = items.Add(item)
	}
	l.Items = items
	return items
}

func (l *List) ToggleSubList(i list.Item) Item {
	return l.Items.ToggleList(i.(Item).id)
}
