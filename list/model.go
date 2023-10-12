package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/pager"
	"github.com/ohzqq/teacozy/util"
)

type State int

const (
	Browsing State = iota
	Input
	Paging
)

func (s State) String() string {
	switch s {
	case Browsing:
		return "browsing"
	case Input:
		return "input"
	case Paging:
		return "paging"
	}
	return ""
}

type Layout int

const (
	Vertical Layout = iota
	Horizontal
)

type Model struct {
	*list.Model
	*Items

	KeyMap        KeyMap
	shortHelpKeys []key.Binding
	fullHelpKeys  []key.Binding

	state  State
	layout Layout

	editable bool
	focused  bool

	showDescription bool

	// Input
	Input    *input.Model
	hasInput bool

	// view
	Pager    *pager.Model
	hasPager bool
}

// Option configures a Model.
type Option func(*Model)

// ListOption configures list.Model.
type ListOption func(*list.Model)

// New initializes a Model.
func New(items *Items, opts ...Option) *Model {
	m := &Model{
		Items:  items,
		state:  Browsing,
		KeyMap: DefaultKeyMap(),
	}
	m.Model = m.NewListModel(items)

	for _, opt := range opts {
		opt(m)
	}

	m.hasInput = true
	m.Input = input.New()
	m.Input.Prompt = "Insert Item: "
	m.Input.Enter = InsertItem
	m.Input.PromptStyle = m.Styles.FilterPrompt
	m.Input.Cursor.Style = m.Styles.FilterCursor
	m.AddFullHelpKeys(m.Items.KeyMap.InsertItem, m.Items.KeyMap.RemoveItem)
	m.AddShortHelpKeys(m.Items.KeyMap.InsertItem, m.Items.KeyMap.RemoveItem)

	m.AdditionalShortHelpKeys = func() []key.Binding {
		return m.shortHelpKeys
	}
	m.AdditionalFullHelpKeys = func() []key.Binding {
		return m.fullHelpKeys
	}

	return m
}

//func (m *Model) Run() (*Model, error) {
//  p := tea.NewProgram(m)

//  mod, err := p.Run()
//  if err != nil {
//    return m, err
//  }
//  rm := mod.(*Model)

//  return rm, nil
//}

// NewListModel returns a *list.Model.
func (m *Model) NewListModel(items *Items) *list.Model {
	del := items.NewDelegate()
	del.ShowDescription = false
	m.Items.DefaultDelegate = del

	w, h := util.TermSize()
	l := list.New(m.Items.li, m.Items, w, h)

	l.KeyMap = m.KeyMap.KeyMap
	l.Title = ""
	l.Styles = DefaultStyles()
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)

	// Update paginator style
	l.Paginator.ActiveDot = l.Styles.ActivePaginationDot.String()
	l.Paginator.InactiveDot = l.Styles.InactivePaginationDot.String()

	return &l
}

//func (m Model) Items()

// ChooseOne configures a list to return a single choice.
func ChooseOne(items *Items, opts ...Option) *Model {
	opts = append(opts, WithLimit(SelectOne))
	m := New(items, opts...)
	return m
}

// ChooseAny configures a list for multiple selections.
func ChooseAny(items *Items, opts ...Option) *Model {
	opts = append(opts, WithLimit(SelectAll))
	m := New(items, opts...)
	return m
}

// ChooseSome configures a list for limited multiple selections.
func ChooseSome(items *Items, limit int, opts ...Option) *Model {
	opts = append(opts, WithLimit(limit))
	m := New(items, opts...)
	return m
}

// Edit configures an editable list: items are not selectable but can be
// removed from the list or new items entered with a prompt.
func Edit(items *Items, opts ...Option) *Model {
	o := []Option{Editable(true)}
	o = append(o, opts...)
	m := New(items, o...)
	return m
}

// Editable marks a list as editable
func Editable(edit bool) Option {
	return func(m *Model) {
		m.SetEditable(edit)
	}
}

// WithFiltering sets filtering on list.Model.
func WithFiltering(f bool) Option {
	return func(m *Model) {
		m.SetFilteringEnabled(f)
	}
}

// WithPager enables the pager bubble.
func WithPager(p *pager.Model) Option {
	return func(m *Model) {
		m.hasPager = true
		switch m.layout {
		case Horizontal:
			m.width = m.width / 2
			p.SetSize(m.width/2, m.height)
		case Vertical:
			m.height = m.height / 2
			p.SetSize(m.width, m.height/2)
		}
		//p.Blur()
		m.Pager = p
	}
}

// WithDescription sets the list to show an item's description.
func WithDescription(desc bool) Option {
	return func(m *Model) {
		del := m.Items.NewDelegate()
		del.ShowDescription = desc
		if desc {
			del.SetHeight(2)
		}
		m.Items.DefaultDelegate = del
		m.SetDelegate(m.Items)
	}
}

// WithLimit sets the limit of choices for a selectable list.
func WithLimit(n int) Option {
	return func(m *Model) {
		if n != SelectNone {
			m.Items.SetLimit(n)
			m.AddFullHelpKeys(m.Items.KeyMap.ToggleItem)
			m.AddShortHelpKeys(m.Items.KeyMap.ToggleItem)
		}
	}
}

// OrderedList sets the list.DefaultDelegate ListType.
func OrderedList() Option {
	return func(m *Model) {
		m.Items.ListType = Ol
	}
}

// UnrderedList sets the list.DefaultDelegate ListType.
func UnorderedList() Option {
	return func(m *Model) {
		m.Items.ListType = Ul
	}
}

// AddShortHelpKeys adds key.Binding to list.Model's short help.
func (m *Model) AddShortHelpKeys(keys ...key.Binding) {
	m.shortHelpKeys = append(m.shortHelpKeys, keys...)
}

// AddFullHelpKeys adds key.Binding to list.Model's full help.
func (m *Model) AddFullHelpKeys(keys ...key.Binding) {
	m.fullHelpKeys = append(m.fullHelpKeys, keys...)
}

func (m *Model) SetEditable(edit bool) {
	m.Items.SetEditable(edit)
}

func (m Model) Focused() bool {
	return m.focused
}

type FocusMsg struct{}
type UnfocusMsg struct{}

func (m *Model) Focus() tea.Cmd {
	return func() tea.Msg {
		m.focused = true
		return FocusMsg{}
	}
}

func (m *Model) Unfocus() tea.Cmd {
	return func() tea.Msg {
		m.focused = false
		return UnfocusMsg{}
	}
}

// State returns the current list state.
func (m Model) State() State {
	return m.state
}

// SetBrowsing sets the state to Browsing
func (m *Model) SetBrowsing() {
	m.state = Browsing
}

// IsBrowsing returns whether or not the list state is Browsing.
func (m Model) Browsing() bool {
	return !m.Model.SettingFilter() || m.state == Browsing
}

// HasInput returns whether or not the model has input.
func (m Model) HasInput() bool {
	return m.hasInput
}

// CurrentItem returns the selected item.
func (m Model) CurrentItem() *Item {
	li := m.Model.SelectedItem()
	return li.(*Item)
}

// Update is the tea.Model update loop.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case InputItemMsg:
		m.SetShowInput(true)
		cmds = append(cmds, m.Input.Focus())
	case ResetInputMsg:
		m.ResetInput()

	case ItemsChosenMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Filter):
			//m.state = Input
		}
	}

	switch m.State() {
	case Input:
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	default:
		li, cmd := m.Model.Update(msg)
		m.Model = &li
		cmds = append(cmds, cmd)
		//if m.showDescription {
		//  m.Pager.SetText(m.CurrentItem().Description())
		//}

	}

	return m, tea.Batch(cmds...)
}

// InputItemMsg is a tea.Msg to focus the item input.Model.
type InputItemMsg struct{}

// InputItem is a tea.Cmd to input an item.
func InputItem() tea.Msg {
	return InputItemMsg{}
}

type ResetInputMsg struct{}

func ResetInput() tea.Msg {
	return ResetInputMsg{}
}

// SetShowInput shows or hides the input model.
func (m *Model) SetShowInput(show bool) {
	m.SetShowTitle(!show)
	if show {
		m.SetHeight(m.Height() - 1)
		m.state = Input
		return
	}
	m.SetHeight(m.Height() + 1)
	m.state = Browsing
}

// ResetInput resets the current input state.
func (m *Model) ResetInput() {
	m.resetInput()
}

func (m *Model) resetInput() {
	if m.state == Browsing {
		return
	}
	m.Input.Reset()
	m.Input.Blur()
	m.SetShowInput(false)
}

func (m *Model) updateBindings() {
	switch m.state {
	case Paging:
		m.KeyMap.CursorUp.SetEnabled(false)
		m.KeyMap.CursorDown.SetEnabled(false)
		m.KeyMap.PrevPage.SetEnabled(false)
		m.KeyMap.NextPage.SetEnabled(false)
		m.KeyMap.GoToStart.SetEnabled(false)
		m.KeyMap.GoToEnd.SetEnabled(false)
		m.KeyMap.Filter.SetEnabled(false)
	case Browsing:
		m.KeyMap.CursorUp.SetEnabled(true)
		m.KeyMap.CursorDown.SetEnabled(true)
		m.KeyMap.PrevPage.SetEnabled(true)
		m.KeyMap.NextPage.SetEnabled(true)
		m.KeyMap.GoToStart.SetEnabled(true)
		m.KeyMap.GoToEnd.SetEnabled(true)
		m.KeyMap.Filter.SetEnabled(true)
	}
}

// View satisfies the tea.Model view method.
func (m *Model) View() string {
	var views []string
	if m.HasInput() {
		m.SetShowFilter(true)
		if m.Input.Focused() {
			m.SetShowFilter(false)
			in := m.Input.View()
			views = append(views, in)
		}
	}
	li := m.Model.View()
	views = append(views, li)
	view := lipgloss.JoinVertical(lipgloss.Left, views...)
	return view
}

func (m *Model) Init() tea.Cmd {
	return nil
}
