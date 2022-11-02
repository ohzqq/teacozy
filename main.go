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

	m := TestList()
	//m := newTestList()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	for _, s := range m.AllItems.GetSelected() {
		item := s.(list.Item)
		println(item.Title())
	}
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

func newTestItems() []list.Item {
	var items []list.Item
	for key, _ := range testData {
		i := list.Item{Content: key}
		items = append(items, i)
	}
	return items
}

func newItemWithList() list.Item {
	item := list.NewItem(list.Item{Content: "sub3"})
	for key, _ := range testSubList {
		i := list.NewItem(list.Item{Content: key})
		i.IsSub = true
		i.Level = 1
		item.Items = append(item.Items, i)
	}
	return item
}

func newTestList() *list.Model {
	l := list.New("test poot toot")
	//l.isPrompt = true

	l.AddMenu(testMenu())
	l.SetMulti()

	il := newItemWithList()
	l.AllItems.Add(il)
	for _, i := range newTestItems() {
		l.AllItems.Add(i)
	}

	l.List = l.BuildModel()

	return l
}
func TestItems() []list.Item {
	var items []list.Item
	for key, _ := range testData {
		i := list.NewDefaultItem(key, key)
		items = append(items, list.NewListItem(i))
	}
	return items
}

func itemWithList(t string) list.Item {
	i := list.NewDefaultItem(t, t)
	i.HasList = true
	l := list.New(t)
	for key, _ := range testSubList {
		i := list.NewDefaultItem(key, key)
		i.IsSub = true
		i.Level = 1
		l.AppendItem(i)
		i.Items = append(i.Items, list.NewListItem(i))
	}

	l.List = l.BuildModel()
	i.List = l
	i.Items = i.List.Items
	return i
}

func TestList() *list.Model {
	l := list.New("test poot toot")
	//l.isPrompt = true

	l.AddMenu(testMenu())
	l.SetMulti()
	//l.showMenu = true

	il := itemWithList("test sub list")
	l.Items = append(l.Items, il)
	for _, i := range TestItems() {
		l.AppendItem(i)
	}
	l.Items = append(l.Items, itemWithList("another sub list"))

	l.List = l.BuildModel()

	return l
}

func testMenu() *list.Menu {
	t := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "deselect all"),
	)
	m := list.NewMenu("test", t)
	m.SetKeys(testHelpKeys)
	return m
}

var testHelpKeys = []list.MenuItem{
	list.NewMenuItem("t", "select item", TestKeyAction),
}

func TestKeyAction(m *list.Model) tea.Cmd {
	return m.List.NewStatusMessage(fmt.Sprintf("%v", m.IsMultiSelect))
}
