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

func menuModel() {
}
func listModel() list.List {

	var testHelpKeys = []list.MenuItem{
		list.NewMenuItem("t", "select item", TestKeyAction),
	}
	//items := list.NewItems()
	l := list.NewMultiSelect("test")
	l.Add(list.NewItem("poot"))
	toot := list.NewItem("toot")
	//toot.Items = list.NewItems()
	//toot.AllItems = append(toot.AllItems, list.NewItem("moot"))
	l.Add(toot)
	t := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "deselect all"),
	)
	//m := items.NewList("test", true)
	menu := list.NewMenu("test", t, testHelpKeys)
	l.NewWidget(menu)
	l.Action = list.SingleSelectAction(l)
	l.Build()
	//m.Action = list.MultiSelectAction(m)
	return l
}

func TestKeyAction(m *list.List) tea.Cmd {
	return m.Model.NewStatusMessage(fmt.Sprintf("%v", m.IsMultiSelect))
}
