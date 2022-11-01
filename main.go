package main

import (
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

func menuModel() {
}
func listModel() list.List {

	var testHelpKeys = []list.MenuItem{
		list.NewMenuItem("t", "select item", TestKeyAction),
	}
	items := list.NewItems()
	items.IsMultiSelect = false
	items.Add(list.NewItem("poot"))
	toot := list.NewItem("toot")
	toot.Items = list.NewItems()
	toot.Items.Add(list.NewItem("moot"))
	items.Add(toot)
	t := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "deselect all"),
	)
	m := items.NewList("test", false)
	menu := list.NewMenu("test", t, testHelpKeys)
	m.NewWidget(menu)
	m.Action = list.SingleSelectAction(m)
	//m.Action = list.MultiSelectAction(m)
	return m
}

func TestKeyAction(m *list.List) tea.Cmd {
	return m.Model.NewStatusMessage("poot")
}
