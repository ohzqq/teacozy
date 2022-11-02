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
				cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
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
					cmds = append(cmds, UpdateDisplayedItemsCmd("selected"))
				case key.Matches(msg, m.Keys.SelectAll):
					ToggleAllItemsCmd(m)
					cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
				}
			} else {
				switch {
				case key.Matches(msg, m.Keys.Enter):
					m.ToggleItem(m.List.SelectedItem())
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
				cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
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
		m.Selections = m.AllItems.GetSelected()
		cmds = append(cmds, tea.Quit)
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width-1, msg.Height-2)
	}

	switch focus := m.FocusedView; focus {
	case "list":
		switch msg := msg.(type) {
		case UpdateDisplayedItemsMsg:
			items := m.DisplayItems(string(msg))
			//m.Model.SetHeight(m.GetHeight(items))
			m.List.SetHeight(util.TermHeight() - 2)
			cmds = append(cmds, m.List.SetItems(items))
		case ToggleItemListMsg:
			cur := m.List.SelectedItem().(Item)
			m.Items.ToggleList(cur.id)
			cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
		case toggleItemMsg:
			cur := m.List.SelectedItem().(Item)
			m.Items.ToggleSelected(cur.id)
			cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
		case UpdateMenuContentMsg:
			m.CurrentMenu.Model.SetContent(string(msg))
			m.ShowMenu = false
		case UpdateItemsMsg:
			m.SetItems(Items(msg))
			m.processAllItems()
			cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
			m.ShowMenu = false
		case OSExecCmdMsg:
			menuCmd := msg.cmd(m.AllItems.GetSelected())
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
			item.state = ItemNotSelected
		}
		items = items.Add(item)
	}

	//l.AllItems = items
	l.Items = items
	return items
}

func (l Model) DisplayItems(opt string) Items {
	return l.Items.Display(opt)
}

func (l *Model) ToggleItem(i list.Item) Item {
	return l.Items.ToggleSelected(i.(Item).id)
}

func (l *Model) ToggleSubList(i list.Item) Item {
	return l.Items.ToggleList(i.(Item).id)
}

//func (l *List) SetItem(i list.Item) {
//  l.AllItems[i.(Item).id] = i
//}
