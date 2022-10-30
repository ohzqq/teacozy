package main

import (
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

}

func menuModel() {
}
func listModel() tea.Model {
	items := list.Items{}
	items.Add(list.NewItem("poot"))
	items.Add(list.NewItem("toot"))
	return items.NewList("test")
}
