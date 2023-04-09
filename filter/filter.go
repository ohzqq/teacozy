package filter

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Filter struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor      int
	Matches     []props.Item
	Input       textinput.Model
	Viewport    *viewport.Model
	end         int
	start       int
	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List
}

type Props struct {
	*props.Items
}

func New() *Filter {
	tm := Filter{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func (c Filter) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (c Filter) Name() string {
	return "filter"
}

func (m *Filter) KeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Up(),
		keys.Down(),
		keys.ToggleItem(),
		keys.Quit(),
		keys.ShowHelp(),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.NewBinding("esc").
			WithHelp("stop filtering").
			Cmd(StopFiltering()),
	}
}

func (m *Filter) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.UpMsg:
		n := 1
		m.Cursor = clamp(m.Cursor-n, 0, len(m.Props().Visible())-1)
		switch {
		case m.start == 0:
			m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
		case m.start < m.Viewport.Height:
			m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
		case m.Viewport.YOffset >= 1:
			m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Viewport.Height)
		}
		m.UpdateViewport()
		//if len(m.Matches) > 0 {
		//  m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor-1)
		//  m.SetCurrent(m.Cursor)
		//  h := m.Matches[m.Cursor].LineHeight()
		//  if m.Cursor < m.Viewport.YOffset {
		//    m.Viewport.LineUp(h)
		//  }
		//}

	case message.DownMsg:
		n := 1
		m.Cursor = clamp(m.Cursor+n, 0, len(m.Props().Visible())-1)
		m.UpdateViewport()

		switch {
		case m.end == len(m.Props().Visible()):
			m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Viewport.Height))
		case m.Cursor > (m.end-m.start)/2:
			m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
		case m.Viewport.YOffset > 1:
		case m.Cursor > m.Viewport.YOffset+m.Viewport.Height-1:
			m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
		}
		//if len(m.Matches) > 0 {
		//  m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor+1)
		//  m.SetCurrent(m.Cursor)

		//  offset := m.Viewport.YOffset
		//  h := m.Matches[m.Cursor].LineHeight()
		//  if o := h - m.Viewport.Height; o > 0 {
		//    m.Viewport.LineDown(o)
		//  } else if m.Cursor <= offset+m.Viewport.Height {
		//    m.Viewport.LineDown(h)
		//  }
		//}

	case message.ToggleItemMsg:
		if len(m.Matches) > 0 {
			m.SetCurrent(m.Matches[m.Cursor].Index)
			if m.Props().NumSelected == 0 && m.quitting {
				cmds = append(cmds, m.ReturnSelections())
			}
			m.Props().ToggleSelection()
			if m.Props().Limit == 1 {
				return m.ReturnSelections()
			}
			cmds = append(cmds, message.Down())
		}

	case StopFilteringMsg:
		if m.Props().Limit == 1 {
			cmds = append(cmds, message.ToggleItem())
		}
		m.Input.Reset()
		m.Input.Blur()
		return message.ChangeRoute("prev")

	case message.ShowHelpMsg:
		k := m.KeyMap()
		m.Props().SetHelp(k)
		cmds = append(cmds, message.ChangeRoute("help"))

	case tea.KeyMsg:
		for _, k := range m.KeyMap() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		m.Input, cmd = m.Input.Update(msg)
		if v := m.Input.Value(); v != "" {
			//m.Matches = m.Props().Visible(m.Input.Value())
			m.Matches = m.Props().Visible(v)
		} else {
			//m.Matches = []props.Item{}
			m.Matches = m.Props().Visible()
		}
		cmds = append(cmds, message.Top())
		cmds = append(cmds, cmd)
	}

	//m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor)
	return tea.Batch(cmds...)
}

func (m *Filter) SetCurrent(idx int) {
	m.Props().SetCurrent(idx)
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *Filter) UpdateViewport() {
	//renderedRows := make([]string, 0, len(m.Matches))
	var renderedRows string

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Viewport.Height, 0, m.Cursor)
	} else {
		m.start = 0
	}
	m.end = clamp(m.Cursor+m.Viewport.Height, m.Cursor, len(m.Props().Visible()))
	//for i := m.start; i < m.end; i++ {
	renderedRows = m.Props().RenderItems(m.Matches[m.start:m.end])
	//renderedRows = append(renderedRows, strconv.Itoa(i))
	//}

	m.Viewport.SetContent(renderedRows)
}

func (m *Filter) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height
	m.Viewport.Width = m.Props().Width
	m.UpdateViewport()
	//m.Viewport.SetContent(m.Props().RenderItems(m.Matches))

	view := m.Input.View() + "\n" + m.Viewport.View()
	return view
}

func (tm *Filter) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)

	tm.Input = textinput.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.PromptStyle = tm.Style.Prompt
	tm.Input.Placeholder = tm.Placeholder
	tm.Input.Width = tm.Props().Width

	v := viewport.New(0, 0)
	tm.Viewport = &v
	tm.Matches = props.Visible()
	tm.UpdateViewport()

	tm.Input.Focus()

	return nil
}

func (m *Filter) ReturnSelections() tea.Cmd {
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}

type StartFilteringMsg struct{}

func StartFiltering() tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}

type StopFilteringMsg struct{}

func StopFiltering() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
