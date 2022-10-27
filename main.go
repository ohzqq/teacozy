package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := frameModel()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

}

func menuModel() {
}
func listModel() {
}
func frameModel() frame.Model {
}
