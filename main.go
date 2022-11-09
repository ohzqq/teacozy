package main

import (
	"fmt"
	"log"

	"github.com/ohzqq/teacozy/form"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/prompt"
	"github.com/ohzqq/teacozy/util"
)

var (
	width  = util.TermWidth()
	height = util.TermHeight()
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	testPrompt()
}

func otherInfo() *info.Fields {
	f := &info.DefaultFields{}
	f.Add("two", "poot")
	f.Add("three", "toot")
	return form.NewInfo(f).Fields
}

func testInfo() *form.Model {
	f := &info.DefaultFields{}
	f.Add("one", "poot")
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
