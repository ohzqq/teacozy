package choose

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/ohzqq/teacozy/keys"
)

func FilterKeyMap(m *Filter) keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("esc"),
			keys.WithHelp("esc", "stop filtering"),
			keys.WithCmd(FStopFilteringCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("enter"),
			keys.WithHelp("enter", "return selections"),
			keys.WithCmd(FReturnSelectionsCmd(m)),
		),
	}
	return km
}

func GlobalsKeyMap(m *Filter) keys.KeyMap {
	return keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("down"),
			keys.WithHelp("down", "move cursor down"),
			keys.WithCmd(FDownCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("up"),
			keys.WithHelp("up", "move cursor up"),
			keys.WithCmd(FUpCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("ctrl+c"),
			keys.WithHelp("ctrl+c", "quit"),
			keys.WithCmd(FQuitCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("tab"),
			keys.WithHelp("tab", "select item"),
			keys.WithCmd(FSelectItemCmd(m)),
		),
	}
}

func GlobalKeyMap(m *Choose) keys.KeyMap {
	return keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("down"),
			keys.WithHelp("down", "move cursor down"),
		),
		keys.NewBinding(
			keys.WithKeys("up"),
			keys.WithHelp("up", "move cursor up"),
		),
		keys.NewBinding(
			keys.WithKeys("ctrl+c"),
			keys.WithHelp("ctrl+c", "quit"),
		),
		keys.NewBinding(
			keys.WithKeys("tab"),
			keys.WithHelp("tab", "select item"),
		),
	}
}

var chooseKey = ChooseKeys{
	Next: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("right/l", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("left/h", "prev page"),
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

func ListKeyMap(m *Choose) keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("right", "l"),
			keys.WithHelp("right/l", "next page"),
			//keys.WithCmd(NextPageCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("left", "h"),
			keys.WithHelp("left/h", "prev page"),
			//keys.WithCmd(PrevPageCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("V"),
			keys.WithHelp("V", "deselect all"),
			//keys.WithCmd(DeselectAllItemsCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("v"),
			keys.WithHelp("v", "select all"),
			//keys.WithCmd(SelectAllItemsCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys(" "),
			keys.WithHelp("space", "select item"),
			//keys.WithCmd(SelectItemCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("j"),
			keys.WithHelp("j", "move cursor down"),
			//keys.WithCmd(CursorDownCmd(m.CursorDown)),
		),
		keys.NewBinding(
			keys.WithKeys("k"),
			keys.WithHelp("k", "move cursor up"),
			//keys.WithCmd(UpCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("esc", "q"),
			keys.WithHelp("esc/q", "quit"),
			//keys.WithCmd(QuitCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("enter"),
			keys.WithHelp("enter", "return selections"),
			//keys.WithCmd(ReturnSelectionsCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("/"),
			keys.WithHelp("/", "filter items"),
			//keys.WithCmd(StartFilteringCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("G"),
			keys.WithHelp("G", "last item"),
			//keys.WithCmd(BottomCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("g"),
			keys.WithHelp("g", "first item"),
			//keys.WithCmd(TopCmd(m)),
		),
	}
	return km
}
