package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var testHelpKeys = []MenuItem{
	NewMenuItem("t", "select item", TestKeyAction),
}

var testData = map[string]string{
	"one":   "poot",
	"two":   "toot",
	"three": "scoot",
}

var testSubList = map[string]string{
	"sub1": "poot",
	"sub2": "toot",
}

func TestItems() Items {
	var items Items
	for key, _ := range testData {
		i := NewDefaultItem(key, key)
		items = append(items, NewListItem(i))
	}
	return items
}

func itemWithList(t string) Item {
	i := NewDefaultItem(t, t)
	i.hasList = true
	l := New(t)
	for key, _ := range testSubList {
		i := NewDefaultItem(key, key)
		i.isSub = true
		i.level = 1
		l.AppendItem(i)
		i.items = append(i.items, NewListItem(i))
	}
	l.Model = l.BuildModel()
	i.list = l
	i.items = i.list.Items
	return i
}

func TestList() *List {
	l := New("test poot toot")
	//l.isPrompt = true

	t := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "deselect all"),
	)
	l.NewMenu("test", t, testHelpKeys)
	l.SetMulti()
	//l.showMenu = true

	il := itemWithList("test sub list")
	l.Items = append(l.Items, il)
	var items []list.Item
	for key, _ := range testData {
		i := NewDefaultItem(key, key)
		l.AppendItem(i)
		items = append(items, i)
	}
	l.Items = append(l.Items, itemWithList("another sub list"))

	l.Model = l.BuildModel()

	return l
}

func TestKeyAction(m *List) tea.Cmd {
	var items Items
	for key, _ := range testData {
		i := NewDefaultItem(key, key)
		items = append(items, i)
	}
	return UpdateItemsCmd(items)
}
