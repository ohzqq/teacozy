package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/menu"
)

type TUI struct {
	Main              tea.Model
	Alt               tea.Model
	Input             textarea.Model
	view              viewport.Model
	prompt            textinput.Model
	info              *info.Info
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
	Style             TUIStyle
	width             int
	height            int
	Hash              map[string]string
	ShortHelp         Help
	Help              *info.Info
	MainMenu          *menu.Menu
	ActionMenu        *menu.Menu
	Menus             menu.Menus
	CurrentMenu       *menu.Menu
}
