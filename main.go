package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/util"
)

var (
	width  = util.TermWidth()
	height = util.TermHeight()
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := listModel()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	for _, s := range m.Items.Selections() {
		println(s.Prefix())
		println(s.Content)
	}
}

var testHelpKeys = []list.MenuItem{
	list.NewMenuItem("t", "select item", TestKeyAction),
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

func TestItems() list.Items {
	var items list.Items
	for key, _ := range testData {
		i := list.NewDefaultItem(key, key)
		items = append(items, list.NewListItem(i))
	}
	return items
}

func itemWithList(t string) list.Item {
	i := list.NewDefaultItem(t, t)
	//i.hasList = true
	l := New(t)
	for key, _ := range testSubList {
		i := list.NewDefaultItem(key, key)
		i.isSub = true
		i.level = 1
		l.AppendItem(i)
		i.items = append(i.items, list.NewListItem(i))
	}
	l.Model = l.BuildModel()
	i.list = l
	i.items = i.list.Items
	return i
}

func TestList() *List {
	l := list.New("test poot toot")
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
		i := list.NewDefaultItem(key, key)
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
		i := list.NewDefaultItem(key, key)
		items = append(items, i)
	}
	return list.UpdateItemsCmd(items)
}
func TestKeyAction(m *list.List) tea.Cmd {
	return m.Model.NewStatusMessage(fmt.Sprintf("%v", m.IsMultiSelect))
}
