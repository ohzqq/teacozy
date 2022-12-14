//go:build ignore

package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/form"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type TUI struct {
	Main              tea.Model
	Alt               tea.Model
	Input             textarea.Model
	view              viewport.Model
	prompt            textinput.Model
	Info              *info.Info
	info              *Info
	state             state
	Title             string
	FocusedView       string
	fullScreen        bool
	actionConfirmed   bool
	showMenu          bool
	showInfo          bool
	showHelp          bool
	showFullHelp      bool
	showConfirm       bool
	currentListItem   int
	currentItemFields *teacozy.Fields
	Style             Style
	width             int
	height            int
	Hash              map[string]string
	Help              Help
	MainMenu          *Menu
	Menus             Menus
	CurrentMenu       *Menu
	//ActionMenu        *menu.Menu
	//ShortHelp         Help
}

func New(main *list.List) TUI {
	mk := key.NewKey("m", "menu")
	ui := TUI{
		Main:        main,
		FocusedView: "list",
		state:       mainModel,
		Style:       DefaultStyle(),
		Menus:       make(Menus),
		MainMenu:    NewMenu(mk),
		Info:        info.New(),
		Help:        NewHelp(),
		info:        NewInfo(),
		//ActionMenu:  ActionMenu(),
		showHelp: true,
	}
	ui.MainMenu.AddKey(ui.Help.Toggle, GoToHelp)
	ui.CurrentMenu = ui.MainMenu
	//ui.SetHelp(Keys.SortList, Keys.Menu, Keys.Help)
	//ui.AddMenu(SortListMenu())
	return ui
}

func (ui *TUI) AddMenu(menu *Menu) {
	k := key.NewKey(menu.Toggle.Name(), menu.Toggle.Content())
	ui.Help.KeyMap.Add(k)
	//SetCmd(GoToMenuCmd(menu))
	//l.MainMenu.Add(k)
	ui.Menus[menu.Label] = menu
}

func (m *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		focus = m.FocusedView
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		//cmds = append(cmds, list.UpdateStatusCmd(msg.String()))
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		//if _, ok := m.Main.(*Form); ok {
		//  switch {
		//  case Keys.ExitScreen.Matches(msg):
		//    m.Main = m.Alt
		//  case Keys.PrevScreen.Matches(msg):
		//    m.Main = m.Alt
		//    //case Keys.SaveAndExit.Matches(msg):
		//    //m.HideInfo()
		//  }
		//}

		//if m.showMenu {
		//cmds = append(cmds, m.UpdateMenu(msg))
		//}

		switch {
		case key.Matches(msg, key.PrevScreen):
			if main := m.Main.(*list.List); main.SelectionList {
				main.SelectionList = false
				m.Main = main
				cmds = append(cmds, list.UpdateVisibleItemsCmd("visible"))
			}
		case key.Matches(msg, key.FullScreen):
			m.fullScreen = !m.fullScreen
			cmds = append(cmds, m.ToggleFullScreenCmd())
		//case Keys.Info.Matches(msg):
		//cmds = append(cmds, HideInfoCmd())
		case key.Matches(msg, key.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, key.HelpKey):
			m.Info = m.info.Info
			cmds = append(cmds, m.ShowHelp())
			//m.Info = m.Help.Info
			//cmds = append(cmds, info.ToggleVisibleCmd())
			//cmds = append(cmds, info.ShowInfoCmd())
			//m.showFullHelp = true
			//cmds = append(cmds, list.UpdateStatusCmd("porot"))
			//cmds = append(cmds, SetFocusedViewCmd("info"))
		case m.MainMenu.Toggle.Matches(msg):
			if focus == "menu" {
				cmds = append(cmds, HideMenuCmd())
			} else {
				cmds = append(cmds, ChangeMenuCmd(m.MainMenu))
			}
			//default:
			//  for _, menu := range m.Menus {
			//    if menu.Toggle.Matches(msg) && len(menu.Keys()) > 0 {
			//      cmds = append(cmds, ChangeMenuCmd(menu))
			//    }
			//  }
		}
	case tea.WindowSizeMsg:
		w := msg.Width - 1
		h := msg.Height - 2
		m.Style.Frame.SetSize(w, h)
		m.SetSize(w, h)
	case EditInfoMsg:
		//if main := m.Main.(*List); main.SelectionList {
		//  m.Alt = main
		//}
		//m.Alt = m.Main
		//form := NewForm(m.info.Fields)
		//m.Main = form
		//m.HideInfo()
		//cmds = append(cmds, SetFocusedViewCmd("list"))
	case ShowItemInfoMsg:
		//m.info = NewInfoForm().SetData(msg.Fields)
		//if main, ok := m.Main.(*List); ok {
		//m.currentListItem = main.Model.Index()
		//item := NewFields()
		//copier.CopyWithOption(item, main.SelectedItem().Fields, copier.Option{DeepCopy: true})
		//m.currentItemFields = item
		//}
		//cmds = append(cmds, ShowInfoCmd())
	case ToggleHelpMsg:
		if m.showFullHelp {
			//m.info.showHelp = false
			cmds = append(cmds, m.HideHelp())
			//cmds = append(cmds, info.HideInfoCmd())
		} else {
			//m.info.showHelp = true
			//m.info.Help = m.Help
			cmds = append(cmds, m.ShowHelp())
		}
	case info.ToggleVisibleMsg:
		if m.Info.Visible {
			cmds = append(cmds, info.HideInfoCmd())
		} else {
			cmds = append(cmds, info.ShowInfoCmd())
		}
	case info.HideInfoMsg:
		m.Info.Hide()
		m.info.Hide()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case info.ShowInfoMsg:
		m.Info.Show()
		m.info.Show()
		cmds = append(cmds, SetFocusedViewCmd("info"))
		m.HideMenu()
	case ChangeMenuMsg:
		m.CurrentMenu = msg.Menu
		cmds = append(cmds, ShowMenuCmd(m.CurrentMenu))
	case ShowMenuMsg:
		m.showMenu = true
		m.CurrentMenu.Show()
		m.Info.Hide()
		cmds = append(cmds, SetFocusedViewCmd(m.CurrentMenu.Label))
	case HideMenuMsg:
		m.HideMenu()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
		//cmds = append(cmds, list.UpdateStatusCmd(string(msg)))
	case ActionMenuMsg:
		//m.CurrentMenu = m.ActionMenu
		//m.ShowMenu()
		//m.HideInfo()
		//cmds = append(cmds, SetFocusedViewCmd("action"))
	case form.ExitFormMsg:
		//m.Main = m.Alt
		//m.Alt = m.Main
		//if form, ok := m.Main.(*Form); ok {
		//m.Main = m.Alt
		//if !form.Changed {
		//form.Model.Select(m.currentListItem)
		//list := m.Main.(*List)
		//cur := list.SelectedItem()
		//cur.Fields = m.currentItemFields
		//cur := m.currentItemFields
		//list.Set(cur.Index(), cur)
		//}
		//}
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case form.FormChangedMsg:
		//var cur list.Item
		//if _, ok := m.Main.(*Form); ok {
		//main.Model.Select(m.currentListItem)
		//cur = main.Model.SelectedItem()
		//}
		//cmds = append(cmds, UpdateStatusCmd("saved"))
		//cmds = append(cmds, ExitFormCmd())
		//cmds = append(cmds, ItemChangedCmd(cur))
		//case SaveAndExitFormMsg:
		//  cmds = append(cmds, msg.Exit(m.Main))
	}

	switch focus {
	case "help":
		//m.Help.Info, cmd = m.Help.Info.Update(msg)
		//cmds = append(cmds, cmd)
	case "info":
		var model tea.Model
		model, cmd = m.Info.Update(msg)
		m.Info = model.(*info.Info)
		//m.Info, cmd = m.Info.Update(msg)
		cmds = append(cmds, cmd)

		//m.info, cmd = m.info.Update(msg)
		//m.view = m.Info.Model
		//m.view.SetContent(m.Help.Render())
		//m.view, cmd = m.view.Update(msg)
		cmds = append(cmds, cmd)
	case "list":
		switch main := m.Main.(type) {
		case *list.List:
			if main.SelectionList {
				cmds = append(cmds, ActionMenuCmd())
			}
			//    var li tea.Model
			//    li, cmd = main.Update(msg)
			//    cmds = append(cmds, cmd)
			//    m.Main = li
			m.Main, cmd = m.Main.Update(msg)
			cmds = append(cmds, cmd)
		}
	default:
		for label, _ := range m.Menus {
			if focus == label {
				cmds = append(cmds, m.CurrentMenu.Update(m, msg))
			}
		}
	}

	//cmds = append(cmds, list.UpdateStatusCmd(focus))
	return m, tea.Batch(cmds...)
}

func (m *TUI) Init() tea.Cmd {
	m.MainMenu.NewSection().SetFields(m.MainMenu.KeyMap)
	m.Help.Render()
	return nil
}

func (m *TUI) View() string {
	var (
		sections     []string
		availHeight  = m.Height()
		widgetWidth  = m.Style.Widget.Width()
		widgetHeight = m.Style.Widget.Height()
	)

	m.SetSize(m.Width(), availHeight)

	var widget string

	if m.showMenu {
		m.CurrentMenu.SetSize(widgetWidth, widgetHeight)
		widget = m.CurrentMenu.View()
		availHeight -= widgetHeight
	}

	if m.showInfo {
		if m.showConfirm {
			//m.info.SetHeight(2)
		}

		m.Info.SetSize(widgetWidth, widgetHeight)
		widget = m.view.View()
		//widget = m.Info.View()
		availHeight -= widgetHeight
	}

	if m.showFullHelp {
		m.info.Model = viewport.New(widgetWidth, widgetHeight)
		m.view = viewport.New(widgetWidth, widgetHeight)
		m.view.SetContent(m.Help.Render())
		//m.info.Model = viewport.New(m.Width(), m.Height())
		m.info.showHelp = true
		//widget = m.info.View()
		widget = m.view.View()
		//widget = m.Help.View()
		availHeight -= widgetHeight
	}

	var status string
	if m.showHelp {
		//status = m.ShortHelp.View()
		availHeight -= lipgloss.Height(status)
	}

	if m.showConfirm {
		//status = m.CurrentMenu.View()
		availHeight -= lipgloss.Height(status)
	}

	content := m.Main.View()
	sections = append(sections, content)

	if m.showMenu || m.showInfo || m.showFullHelp {
		sections = append(sections, widget)
	}

	if m.showHelp || m.showConfirm {
		sections = append(sections, status)
	}

	if main, ok := m.Main.(*list.List); ok {
		if min := main.Frame.MinHeight; min > availHeight {
			if m.showMenu || m.showInfo {
				return lipgloss.NewStyle().Height(availHeight).Render(widget)
			}
			return lipgloss.NewStyle().Height(availHeight).Render(content)
		}
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m *TUI) Start() *TUI {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

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

func (l *TUI) ShowMenu() {
	l.showMenu = true
	l.CurrentMenu.Show()
}

func (l *TUI) HideMenu() {
	l.showMenu = false
	l.CurrentMenu.Hide()
}

func (ui *TUI) ShowHelp() tea.Cmd {
	ui.Info = ui.Help.Info
	ui.showFullHelp = true
	//ui.Style.Widget = style.DefaultFrameStyle()
	ui.info.showHelp = true
	return info.ShowInfoCmd()
}

func (ui *TUI) HideHelp() tea.Cmd {
	ui.showFullHelp = false
	ui.Style.Widget = DefaultWidgetStyle()
	return info.HideInfoCmd()
}
