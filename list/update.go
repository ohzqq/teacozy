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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.area.Focused() {
			if key.Matches(msg, urkey.SaveAndExit) {
				cur := m.List.SelectedItem().(Item)
				val := m.area.Value()
				cur.SetContent(val)
				m.SetItem(m.List.Index(), cur)
				m.area.Blur()
				cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			}
			m.area, cmd = m.area.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			if m.IsMulti() {
				switch {
				case key.Matches(msg, m.Keys.Enter):
					if m.ShowSelectedOnly {
						cmds = append(cmds, ReturnSelectionsCmd())
					}
					m.ShowSelectedOnly = true
					cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
				case key.Matches(msg, m.Keys.SelectAll):
					ToggleAllItemsCmd(m)
					cmds = append(cmds, UpdateVisibleItemsCmd("all"))
				}
			} else {
				switch {
				case key.Matches(msg, m.Keys.Enter):
					cur := m.List.SelectedItem().(Item)
					m.SetItem(m.List.Index(), cur.ToggleSelected())
					cmds = append(cmds, ReturnSelectionsCmd())
				}
			}

			switch {
			case key.Matches(msg, m.Keys.ExitScreen):
				cmds = append(cmds, tea.Quit)
			case key.Matches(msg, m.Keys.Quit):
				cmds = append(cmds, tea.Quit)
			case key.Matches(msg, m.Keys.Prev):
				m.ShowSelectedOnly = false
				cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			default:
				for label, menu := range m.Menus {
					if key.Matches(msg, menu.Toggle) {
						m.CurrentMenu = menu
						m.ShowMenu = !m.ShowMenu
						cmds = append(cmds, SetFocusedViewCmd(label))
					}
				}
			}
			m.List, cmd = m.List.Update(msg)
			cmds = append(cmds, cmd)
		}
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case EditItemMsg:
		cur := m.List.SelectedItem().(Item)
		m.area = cur.Edit()
		m.area.Focus()
	case ReturnSelectionsMsg:
		m.Selections = m.AllItems.Selected()
		cmds = append(cmds, tea.Quit)
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width-1, msg.Height-2)
	}

	switch focus := m.FocusedView; focus {
	case "list":
		switch msg := msg.(type) {
		case UpdateVisibleItemsMsg:
			items := m.DisplayItems(string(msg))
			//m.Model.SetHeight(m.GetHeight(items))
			m.List.SetHeight(util.TermHeight() - 2)
			cmds = append(cmds, m.List.SetItems(items))
		case ToggleItemListMsg:
			cur := m.Items.Get(int(msg))
			m.SetItem(m.List.Index(), cur.ToggleList())
			cmds = append(cmds, UpdateVisibleItemsCmd("all"))
		case ToggleSelectedItemMsg:
			cur := m.Items.Get(int(msg))
			m.SetItem(m.List.Index(), cur.ToggleSelected())
			cmds = append(cmds, UpdateVisibleItemsCmd("all"))
		case UpdateMenuContentMsg:
			m.CurrentMenu.Model.SetContent(string(msg))
			m.ShowMenu = false
		case SetItemsMsg:
			m.SetItems(Items(msg))
			m.processAllItems()
			cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			m.ShowMenu = false
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
	default:
		for label, _ := range m.Menus {
			if focus == label {
				cmds = append(cmds, UpdateMenu(m, msg))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (l *Model) processAllItems() Items {
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

func (l Model) DisplayItems(opt string) Items {
	return l.Items.Display(opt)
}

func (l *Model) ToggleSubList(i list.Item) Item {
	return l.Items.ToggleList(i.(Item).id)
}
