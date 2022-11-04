package list

import (
	"bytes"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/slices"
)

type List struct {
	Model            list.Model
	Items            Items
	OGitems          []Item
	Title            string
	Keys             urkey.KeyMap
	IsPrompt         bool
	IsMultiSelect    bool
	ShowSelectedOnly bool
}

func NewList(title string) List {
	return List{
		Keys:  urkey.DefaultKeys(),
		Title: title,
	}
}

func (l *List) SetModel() {
	w := util.TermWidth()
	h := util.TermHeight()
	del := NewItemDelegate(l.IsMulti())
	l.AllItems()
	l.Model = list.New(l.VisibleItems(), del, w, h)
	l.Model.Title = l.Title
	l.Model.Styles = style.ListStyles()
	l.Model.SetShowStatusBar(false)
	l.Model.SetShowHelp(false)
}

func (l *List) NewItem(content string) Item {
	i := Item{Content: content}
	l.OGitems = append(l.OGitems, i)
	return i
}

func (l *List) GetItemIndex(i list.Item) int {
	content := i.(Item).Content
	fn := func(item list.Item) bool {
		c := item.(Item).Content
		return content == c
	}
	return slices.IndexFunc(l.Items, fn)
}

func (l *List) AllItems() Items {
	l.Items = FlattenItems(l.OGitems)
	//l.items = l.FlatItems()
	l.ProcessItems()
	return l.Items
}

func (l *List) FlatItems() []Item {
	var items []Item
	for _, item := range l.OGitems {
		sub := item.Flatten()
		items = append(items, sub...)
	}
	return items
}

func FlattenItems(li []Item) Items {
	var items Items
	for _, item := range li {
		items = append(items, item)
		if item.HasList() {
			subList := FlattenItems(item.Li)
			for _, i := range subList {
				sub := i.(Item)
				sub.IsHidden = true
				items = append(items, sub)
			}
		}
	}
	return items
}

func (l *List) ProcessItems() {
	for _, item := range l.Items {
		i := item.(Item)
		if l.IsMulti() {
			i.IsMulti = true
		}
		i.id = l.GetItemIndex(item)
		l.Items[i.id] = i
	}
}

func (l *List) AppendItem(item Item) *List {
	if l.IsMulti() {
		item.IsMulti = true
	}
	l.OGitems = append(l.OGitems, item)
	return l
}

func (l List) Get(idx int) Item {
	if idx < len(l.Items) {
		return l.Items[idx].(Item)
	}
	return Item{}
}

func (l List) DisplayItems(opt string) Items {
	switch opt {
	case "selected":
		return l.Items.Selected()
	default:
		return l.Items.Visible()
	}
}

func (l List) VisibleItems() Items {
	var items Items
	level := 0
	for _, item := range l.Items {
		i := item.(Item)
		if !i.IsHidden {
			items = append(items, i)
		}
		if i.HasList() && i.listOpen {
			level++
			for _, sub := range l.GetItemSubList(i) {
				s := sub.(Item)
				s.IsHidden = false
				s.SetLevel(level)
				items = append(items, s)
			}
		}
	}
	return items
}

func (l List) SelectedItems() Items {
	var items Items
	for _, item := range l.Items {
		if i, ok := item.(Item); ok && i.IsSelected() {
			items = append(items, i)
		}
	}
	return items
}

func (l List) SelectAllItems() Items {
	var items Items
	for _, i := range l.Items {
		item := i.(Item)
		item.isSelected = true
		items = append(items, item)
	}
	return items
}

func (l List) GetItemSubList(i list.Item) Items {
	item := i.(Item)
	if item.HasList() {
		t := len(item.Items)
		return l.Items[item.id+1 : item.id+t+1]
	}
	return Items{}
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
		//m.processAllItems()
		cmds = append(cmds, m.Model.NewStatusMessage("set items"))
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
	return SetItemsCmd(l.AllItems())
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
