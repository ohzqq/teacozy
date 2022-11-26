package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/list"
)

type TUI struct {
	Main   tea.Model
	Alt    tea.Model
	Input  textarea.Model
	view   viewport.Model
	prompt textinput.Model
	//info              *info.Info
	Title             string
	FocusedView       string
	fullScreen        bool
	actionConfirmed   bool
	showMenu          bool
	showInfo          bool
	showHelp          bool
	showConfirm       bool
	currentListItem   int
	currentItemFields *teacozy.FormData
	Style             Style
	width             int
	height            int
	Hash              map[string]string
	//Help              *info.Info
	//MainMenu          *menu.Menu
	//ActionMenu        *menu.Menu
	//Menus             menu.Menus
	//CurrentMenu       *menu.Menu
	//ShortHelp         Help
}

func New(main *list.List) TUI {
	ui := TUI{
		Main:        main,
		FocusedView: "list",
		Style:       DefaultStyle(),
		//Menus:       make(menu.Menus),
		//MainMenu:    menu.New("m", "menu", key.NewKeyMap()),
		//ActionMenu:  ActionMenu(),
		showHelp: true,
	}
	//ui.SetHelp(Keys.SortList, Keys.Menu, Keys.Help)
	//ui.AddMenu(SortListMenu())
	return ui
}

//func (l *TUI) AddMenu(menu *Menu) {
//  k := key.NewKey(menu.Toggle.Help().Key, menu.Toggle.Help().Desc).
//    SetCmd(GoToMenuCmd(menu))
//  l.MainMenu.AddKey(k)
//  l.Menus[menu.Label] = menu
//}

func (ui *TUI) SetSize(w, h int) *TUI {
	switch main := ui.Main.(type) {
	case *list.List:
		main.SetSize(w, h)
	}
	return ui
}

func (ui TUI) Width() int {
	return ui.Style.Frame.Width()
}

func (ui TUI) Height() int {
	return ui.Style.Frame.Height()
}
