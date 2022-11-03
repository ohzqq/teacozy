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

	//m := TestList()
	//m := newTestList()

	//p := tea.NewProgram(m)
	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}

	//fmt.Printf("%+V\n", m.AllItems)

	//for _, s := range m.Items {
	//  item := s.(list.Item)
	//  println(item.Content)
	//}
	//for _, s := range m.Items.Selected() {
	//  item := s.(list.Item)
	//  println(item.Content)
	//}

	l := testList()
	for _, item := range l.AllItems() {
		fmt.Printf("%+V\n", item.(list.Item).Content)
		fmt.Printf("%+V\n", item.(list.Item).Index())
	}
	//items := testItems().Flatten()

	//for _, i := range items {
	//  item := i.(list.Item)
	//if item.HasList() {
	//items = list.Flatten(item.Items)
	//items = append(items, list.Flatten(item.Items)...)
	//}
	//}
	//for _, i := range items {
	//  item := i.(list.Item)
	//  fmt.Printf("%+V\n", item.Content)
	//}
	//fmt.Println(len(items))
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

func testList() list.List {
	l := list.NewList()
	sub3 := l.NewItem("sub3")
	sub3.Info = infoWidget()
	for key, _ := range testSubList {
		i := l.NewItem(key)
		i.SetLevel(1)
		sub3.Li = append(sub3.Li, i)
	}
	for key, _ := range testData {
		l.NewItem(key)
	}

	return l
}

func newItemWithList() list.Item {
	item := list.NewItem(list.Item{Content: "sub3"})
	for key, _ := range testSubList {
		i := list.NewItem(list.Item{Content: key})
		i.SetLevel(1)
		item.Items = item.Items.Add(i)
		//item.Items = append(item.Items, i)
	}
	return item
}

func testItems() list.Items {
	var items list.Items
	il := newItemWithList()
	il.Info = infoWidget()
	items = append(items, il)
	for key, _ := range testData {
		i := list.Item{Content: key}
		items = append(items, i)
	}
	return items
}

func TestList() *list.Model {
	l := list.New("test poot toot")
	//l.isPrompt = true

	l.AddMenu(testMenu())
	l.SetMulti()
	//l.showMenu = true

	//il := itemWithList("test sub list")
	il := newItemWithList()
	il.Info = infoWidget()
	l.AppendItem(il)
	//l.Items = append(l.Items, il)
	//for _, i := range TestItems() {
	for key, _ := range testData {
		//i := list.Item{Content: key}
		//l.AppendItem(i)
		l.NewItem(key)
	}

	l.List.Model = l.BuildModel()

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

func infoWidget() *list.InfoWidget {
	info := list.NewInfoWidget()
	info.AddString("poot", "toot")
	info.AddString("newt", "root")
	for key, val := range testSubList {
		info.AddString(key, val)
	}
	for key, val := range testData {
		info.AddString(key, val)
	}
	return info
}

var testHelpKeys = []list.MenuItem{
	list.NewMenuItem("t", "select item", TestKeyAction),
}

func TestKeyAction(m *list.Model) tea.Cmd {
	return list.UpdateStatusCmd(fmt.Sprintf("%v", m.IsMultiSelect))
}
