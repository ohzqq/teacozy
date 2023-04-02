package choose

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
)

type KeyMap struct {
	Up               key.Binding
	Down             key.Binding
	Prev             key.Binding
	Next             key.Binding
	ToggleItem       key.Binding
	Help             key.Binding
	Quit             key.Binding
	ReturnSelections key.Binding
	Filter           key.Binding
	Bottom           key.Binding
	Top              key.Binding
	Edit             key.Binding
}

var Keys = keys.KeyMap{
	keys.ShowHelp(),
	keys.Quit().WithKeys("ctrl+c", "q"),
	keys.ToggleItem().WithKeys("tab", " "),
	keys.NewBinding("e").
		WithHelp("edit field").
		Cmd(message.StartEditing()),
	keys.NewBinding("enter").
		WithHelp("return selections").
		Cmd(message.ReturnSelections()),
	keys.NewBinding("/").
		WithHelp("filter list").
		Cmd(message.StartFiltering()),
}

var Key = KeyMap{
	Next: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("right/l", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("left/h", "prev page"),
	),
	Help: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "help"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit form"),
	),
	//key.NewBinding(
	//  key.WithKeys("V"),
	//  key.WithHelp("V", "deselect all"),
	//),
	//key.NewBinding(
	//  key.WithKeys("v"),
	//  key.WithHelp("v", "select all"),
	//),
	ToggleItem: key.NewBinding(
		key.WithKeys(" ", "tab"),
		key.WithHelp("space", "select item"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "move cursor down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "move cursor up"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "q", "ctrl+c"),
		key.WithHelp("esc/q", "quit"),
	),
	ReturnSelections: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "return selections"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter items"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "last item"),
	),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "first item"),
	),
}
