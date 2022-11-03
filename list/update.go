package list

import (
	"bytes"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	urkey "github.com/ohzqq/teacozy/key"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		focus = m.FocusedView
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.Keys.Quit) {
			cmds = append(cmds, tea.Quit)
		}
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
			switch focus {
			case "info":
				cmds = append(cmds, UpdateInfoWidget(m, msg))
			case "list":
				cmds = append(cmds, UpdateList(m, msg))
			}
			for label, menu := range m.Menus {
				if key.Matches(msg, menu.Toggle) {
					m.CurrentMenu = menu
					m.ShowMenu()
					cmds = append(cmds, SetFocusedViewCmd(label))
				}
			}
		}
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case EditItemMsg:
		cur := m.List.SelectedItem().(Item)
		m.area = cur.Edit()
		m.area.Focus()
	case ReturnSelectionsMsg:
		m.Selections = m.Items.Selected()
		cmds = append(cmds, tea.Quit)
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width-1, msg.Height-2)
		m.info = viewport.New(msg.Width-2, msg.Height/3)
	case UpdateInfoContentMsg:
		m.ShowInfo()
		m.info.SetContent(string(msg))
		cmds = append(cmds, SetFocusedViewCmd("info"))
	case UpdateMenuContentMsg:
		m.CurrentMenu.Model.SetContent(string(msg))
		m.HideMenu()
	}

	switch focus {
	case "info":
		//cmds = append(cmds, UpdateInfoWidget(m, msg))
	case "list":
		switch msg := msg.(type) {
		case SetItemsMsg:
			m.SetItems(Items(msg))
			m.processAllItems()
			cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			m.showMenu = false
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
