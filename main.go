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

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	for _, s := range m.AllItems.GetSelected() {
		item := s.(list.Item)
		println(item.Title())
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
	i.HasList = true
	l := list.New(t)
	for key, _ := range testSubList {
		i := list.NewDefaultItem(key, key)
		i.IsSub = true
		i.Level = 1
		l.AppendItem(i)
		i.Items = append(i.Items, list.NewListItem(i))
	}

	l.Model = l.BuildModel()
	i.List = l
	i.Items = i.List.Items
	return i
}

func TestList() *list.List {
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

func TestKeyAction(m *list.List) tea.Cmd {
	return m.Model.NewStatusMessage(fmt.Sprintf("%v", m.IsMultiSelect))
}
