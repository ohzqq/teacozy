package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/confirm"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/view"
)

type Option func(*Component)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor int
	KeyMap keys.KeyMap

	view  *view.Model
	start int
	end   int
}

type Props struct {
	view.Props
	Selectable  bool
	ToggleItems func(...int)
	ShowHelp    func([]map[string]string)
}

func New() *Component {
	m := Component{
		Cursor: 0,
	}
	m.DefaultKeyMap()

	return &m
}

func (m *Component) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	props.Selectable = true
	m.view = view.NewModel(props.Props)
	m.view.UpdateItems()
	return nil
}

func (m *Component) AfterUpdate() tea.Cmd {
	m.view.UpdateItems()
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case keys.ReturnSelectionsMsg:
		if len(m.Props().Selected) == 0 {
			return nil
		}
		if reactea.CurrentRoute() == "list" {
			return confirm.GetConfirmation("Accept selected?", AcceptChoices)
		}

	case keys.ShowHelpMsg:
		h := m.KeyMap.Map()
		h = append(h, map[string]string{"Filtering Keys": "\n"})
		h = append(h, keys.TextInput().Map()...)
		m.Props().ShowHelp(h)
		cmds = append(cmds, keys.ChangeRoute("help"))

	case keys.ToggleAllItemsMsg:
		for _, i := range m.Props().Matches() {
			m.Props().ToggleItems(i.Index)
		}
	case keys.ToggleItemMsg:
		m.Props().ToggleItems(m.view.CurrentItem())
		cmds = append(cmds, keys.LineDown)

	case tea.KeyMsg:
		if reactea.CurrentRoute() == "list" {
			if m.Props().Editable {
				if k := keys.Edit(); key.Matches(msg, k.Binding) {
					return k.TeaCmd
				}
			}
			if m.Props().Filterable {
				if k := keys.Filter(); key.Matches(msg, k.Binding) {
					return k.TeaCmd
				}
			}
		}
		for _, k := range m.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Component) Render(w, h int) string {
	m.view.SetWidth(w)
	m.view.SetHeight(h)
	return m.view.View()
}

// CurrentItem returns the selected row.
// You can cast it to your own implementation.
func (m Component) CurrentItem() int {
	return m.view.CurrentItem()
}

func (m *Component) commonKeys() keys.KeyMap {
	var km = []*keys.Binding{
		keys.PgUp(),
		keys.PgDown(),
		keys.Enter().
			WithHelp("return selections").
			Cmd(keys.ReturnSelections),
	}
	return keys.NewKeyMap(km...)
}

// SetKeyMap sets the keymap for the list.
func (m *Component) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = m.commonKeys()
	m.KeyMap.AddBinds(km.Keys()...)
}

func (m *Component) VimKeyMap() *Component {
	m.SetKeyMap(VimKeyMap())
	m.KeyMap.AddBinds(keys.Help().AddKeys("h"))
	return m
}

func (m *Component) DefaultKeyMap() *Component {
	m.SetKeyMap(DefaultKeyMap())
	return m
}

// AcceptChoices returns a confirmation dialogue.
func AcceptChoices(accept bool) tea.Cmd {
	if accept {
		return reactea.Destroy
	}
	return keys.ReturnToList
}

func DefaultKeyMap() keys.KeyMap {
	km := []*keys.Binding{
		keys.Toggle(),
		keys.Up(),
		keys.Down(),
		keys.HalfPgUp(),
		keys.HalfPgDown(),
		keys.Home(),
		keys.End(),
		keys.Quit(),
		keys.New("ctrl+a").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
	}
	return keys.NewKeyMap(km...)
}

func VimKeyMap() keys.KeyMap {
	km := []*keys.Binding{
		keys.Toggle().AddKeys(" "),
		keys.Up().WithKeys("k"),
		keys.Down().WithKeys("j"),
		keys.HalfPgUp().WithKeys("K"),
		keys.HalfPgDown().WithKeys("J"),
		keys.Home().WithKeys("g"),
		keys.End().WithKeys("G"),
		keys.Quit().AddKeys("q"),
		keys.New("ctrl+a", "v").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
	}
	return keys.NewKeyMap(km...)
}
