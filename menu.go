package teacozy

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Menus map[string]*Menu

func (m Menus) Get(key string) *Menu {
	return m[key]
}

type Menu struct {
	Model     viewport.Model
	width     int
	Toggle    key.Binding
	height    int
	Label     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	Keys      []MenuItem
	Frame
}

func NewMenu(l string, toggle key.Binding, items ...MenuItem) *Menu {
	m := DefaultMenu().SetKeys(items...)
	m.Label = l
	m.Toggle = toggle
	return m
}

func DefaultMenu() *Menu {
	m := Menu{
		Frame: DefaultWidgetStyle(),
	}
	m.Model = viewport.New(m.Width(), m.Height())
	return &m
}

func (m *Menu) SetKeys(keys ...MenuItem) *Menu {
	m.Keys = keys
	return m
}

func (m *Menu) NewKey(k, h string, cmd MenuFunc) *Menu {
	key := NewMenuItem(k, h, cmd)
	m.AddKey(key)
	return m
}

func (m *Menu) AddKey(key MenuItem) *Menu {
	m.Keys = append(m.Keys, key)
	return m
}

func (m *Menu) SetLabel(l string) *Menu {
	m.Label = l
	return m
}

func (m *Menu) SetToggle(toggle, help string) *Menu {
	m.Toggle = key.NewBinding(
		key.WithKeys(toggle),
		key.WithHelp(toggle, help),
	)
	return m
}

func (m *Menu) View() string {
	m.Model.SetContent(m.Render())
	return m.Model.View()
}

func (m Menu) Render() string {
	var kh []string
	for _, k := range m.Keys {
		kh = append(kh, k.String())
	}
	return m.Style.Render(strings.Join(kh, "\n"))
}

type MenuItem struct {
	Key key.Binding
	Cmd MenuFunc
}

func NewMenuItem(k, h string, cmd MenuFunc) MenuItem {
	return MenuItem{
		Key: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
		Cmd: cmd,
	}
}

func (i MenuItem) String() string {
	return i.Key.Help().Key + ": " + i.Key.Help().Desc
}
