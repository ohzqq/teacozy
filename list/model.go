package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/util"
)

type State int

const (
	Browsing State = iota
	Input
)

type Layout int

const (
	Horizontal Layout = iota
	Vertical
)

type Model struct {
	*list.Model
	width         int
	height        int
	editable      bool
	shortHelpKeys []key.Binding
	fullHelpKeys  []key.Binding

	KeyMap KeyMap

	toggledItems map[int]struct{}
	limit        int

	items *Items
	state State

	DelegateUpdateFuncs []func(tea.Msg, *list.Model) tea.Cmd

	// input
	input    *input.Model
	hasInput bool
}

// Option configures a Model.
type Option func(*Model)

// ListOption configures list.Model.
type ListOption func(*list.Model)

// New initializes a Model.
func New(items *Items, opts ...Option) *Model {
	w, h := util.TermSize()
	m := &Model{
		width:        w,
		height:       h,
		items:        items,
		state:        Browsing,
		limit:        0,
		KeyMap:       DefaultKeyMap(),
		toggledItems: make(map[int]struct{}),
	}

	m.Model = m.NewListModel(items)

	m.Model.SetSelectNone() //remove later

	for _, opt := range opts {
		opt(m)
	}
	println(m.items.limit)

	m.AdditionalShortHelpKeys = m.shortHelp
	m.AdditionalFullHelpKeys = m.fullHelp

	return m
}

func (m *Model) Run() (*Model, error) {
	p := tea.NewProgram(m)

	mod, err := p.Run()
	if err != nil {
		return m, err
	}
	rm := mod.(*Model)

	return rm, nil
}

func (m *Model) ConfigureList(opts ...ListOption) {
	for _, opt := range opts {
		opt(m.Model)
	}
}

// NewListModel returns a *list.Model.
func (m *Model) NewListModel(items *Items) *list.Model {
	var li []list.Item
	for _, i := range items.ParseFunc() {
		li = append(li, i)
	}
	l := list.New(li, m.items, m.width, m.height)
	l.KeyMap = m.KeyMap.KeyMap
	return &l
}

// ChooseOne configures a list to return a single choice.
func ChooseOne(items *Items, opts ...Option) *Model {
	opts = append(opts, WithLimit(1))
	m := New(items, opts...)
	return m
}

// ChooseAny configures a list for multiple selections.
func ChooseAny(items *Items, opts ...Option) *Model {
	opts = append(opts, WithLimit(-1))
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
	o := []Option{Editable()}
	o = append(o, opts...)
	m := New(items, o...)
	//m.ConfigureList(WithFiltering(false))
	return m
}

// Editable marks a list as editable
func Editable() Option {
	return func(m *Model) {
		m.items.editable = true
		m.items.SetSelectNone()
		m.SetInput("Insert Item: ", InsertItem)
		m.SetFilteringEnabled(false)
		m.items.ShortHelpFunc = func() []key.Binding {
			return []key.Binding{
				keyMap.InsertItem,
				keyMap.RemoveItem,
			}
		}
	}
}

// WithFiltering sets filtering on list.Model.
func WithFiltering(f bool) Option {
	return func(m *Model) {
		m.SetFilteringEnabled(f)
	}
}

// WithLimit sets the limit of choices for a selectable list.
func WithLimit(n int) Option {
	return func(m *Model) {
		m.items.SetLimit(n)
	}
}

// OrderedList sets the list.DefaultDelegate ListType.
func OrderedList() Option {
	return func(m *Model) {
		m.items.ListType = Ol
	}
}

// UnrderedList sets the list.DefaultDelegate ListType.
func UnorderedList() Option {
	return func(m *Model) {
		m.items.ListType = Ul
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

// AdditionalShortHelpKeys adds key.Binding to list.Model's short help.
func AdditionalShortHelpKeys(keys ...key.Binding) Option {
	return func(m *Model) {
		m.AddShortHelpKeys(keys...)
	}
}

// AdditionalFullHelpKeys adds key.Binding to list.Model's full help.
func AdditionalFullHelpKeys(keys ...key.Binding) Option {
	return func(m *Model) {
		m.AddFullHelpKeys(keys...)
	}
}

func (m Model) fullHelp() []key.Binding {
	return m.fullHelpKeys
}

func (m Model) shortHelp() []key.Binding {
	return m.shortHelpKeys
}

// State returns the current list state.
func (m Model) State() State {
	return m.state
}

// SetInput configures an input.Model with the default list.Model styles.
func (m *Model) SetInput(prompt string, enter input.EnterInput) {
	m.hasInput = true
	m.input = input.New()
	m.input.Prompt = prompt
	m.input.Enter = enter
	//m.input.PromptStyle = m.Styles.FilterPrompt
	//m.input.Cursor.Style = m.Styles.FilterCursor
}

func (m Model) UpdateItems(msg tea.Msg, li *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	for _, update := range m.DelegateUpdateFuncs {
		cmds = append(cmds, update(msg, li))
	}
	return tea.Batch(cmds...)
}

// SetBrowsing sets the state to Browsing
func (m *Model) SetBrowsing() {
	m.state = Browsing
}

func (m Model) ToggledItems() []int {
	return m.items.ToggledItems()
}

// IsBrowsing returns whether or not the list state is Browsing.
func (m Model) Browsing() bool {
	return !m.Model.SettingFilter() || m.state == Browsing
}

// Update is the tea.Model update loop.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case input.FocusInputMsg:
		if m.hasInput {
			m.SetShowInput(true)
			cmds = append(cmds, m.input.Focus())
		}
	case input.ResetInputMsg:
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
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	default:
		li, cmd := m.Model.Update(msg)
		m.Model = &li
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
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
	m.SetBrowsing()
}

// ResetInput resets the current input state.
func (m *Model) ResetInput() {
	m.resetInput()
}

func (m *Model) resetInput() {
	if m.state == Browsing {
		return
	}
	m.input.Reset()
	m.input.Blur()
	m.SetShowInput(false)
}

// View satisfies the tea.Model view method.
func (m *Model) View() string {
	var views []string

	if m.hasInput {
		m.SetShowFilter(true)
		if m.input.Focused() {
			m.SetShowFilter(false)
			in := m.input.View()
			views = append(views, in)
		}
	}

	li := m.Model.View()
	views = append(views, li)

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (m *Model) Init() tea.Cmd {
	return nil
}
