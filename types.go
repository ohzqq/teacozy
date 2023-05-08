package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Page interface {
	Header() reactea.SomeComponent
	Main() reactea.SomeComponent
	Footer() reactea.SomeComponent
}

type Route func(map[string]string) (Page, tea.Cmd)

type Routes map[string]Route
