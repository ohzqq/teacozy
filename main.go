package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/form"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/list/item"
	"github.com/ohzqq/teacozy/menu"
	"github.com/ohzqq/teacozy/prompt"
	"github.com/ohzqq/teacozy/util"
)

var (
	width  = util.TermWidth()
	height = util.TermHeight()
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//m := testMenu()
	//m := testInfo()

	//m.Start()
	//p := tea.NewProgram(m)
	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}

	//fmt.Printf("%+V\n", m.Hash)
	testPrompt()
}

func testInfo() *form.Model {
	f := &form.DefaultFields{}
	//f := info.NewFields()
	f.Add("one", "poot")
	f.Add("kjl", "toot")
	f.Add("kjl", "toot")
	i := form.NewInfo(f)
	//i.Info.NoKeys()
	//fmt.Println(i.String())
	return i
}

func testPrompt() {
	items := newItems()
	m := prompt.NewPrompt()
	//m.MultiSelect = false
	m.SetItems(items)
	m.SetMultiSelect()
	m.Start()

	for _, i := range m.Items.Selections() {
		fmt.Printf("%v\n", i.Content)
	}
}

func newItems() item.Items {
	items := item.NewItems()
	sub3 := item.NewDefaultItem("sub3")
	sub3.List = subList()
	subsub3 := item.NewDefaultItem("subsub3")
	subsub3.List = subList()
	sub3.List.Add(subsub3)
	items.Add(sub3)
	for l, c := range testData {
		i := item.NewDefaultItem(l)
		//i.ToggleSelected()
		i.SetLabel(c)
		items.Add(i)
	}
	return items
}

func subList() item.Items {
	var items item.Items
	for key, _ := range testSubList {
		i := item.NewDefaultItem(key)
		i.SetLevel(1)
		items.Add(i)
	}
	return items
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
	l := list.NewList("test")
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

	//l.AddMenu(testMenu())
	l.SetMulti()
	l.List = testList()
	//l.showMenu = true

	//il := itemWithList("test sub list")
	//il := newItemWithList()
	//il.Info = infoWidget()
	//l.AppendItem(il)
	//l.Items = append(l.Items, il)
	//for _, i := range TestItems() {
	//for key, _ := range testData {
	//i := list.Item{Content: key}
	//l.AppendItem(i)
	//l.NewItem(key)
	//}

	//l.List.Model = l.BuildModel()

	return l
}

func testMenu() *menu.Menu {
	t := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "deselect all"),
	)
	testHelpKeys := []menu.Item{
		menu.NewItem("t", "select item", TestKeyAction),
		menu.NewItem("o", "deselect item", TestKeyAction),
	}
	m := menu.NewMenu("test", t, testHelpKeys...)
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

func TestKeyAction() tea.Cmd {
	//return list.UpdateStatusCmd(fmt.Sprintf("%v", m.IsMultiSelect))
	return menu.UpdateMenuContentCmd("update")
}
