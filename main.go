package main

import (
	"fmt"
	"log"

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

	fmt.Printf("%+V\n", m.Items.All)
}

func menuModel() {
}
func listModel() list.List {
	items := list.NewItems()
	items.MultiSelect = false
	items.Add(list.NewItem("poot"))
	toot := list.NewItem("toot")
	toot.Items = list.NewItems()
	toot.Items.Add(list.NewItem("moot"))
	items.Add(toot)
	m := items.NewList("test", false)
	m.Action = list.SingleSelectAction(m)
	//m.Action = list.MultiSelectAction(m)
	return m
}
