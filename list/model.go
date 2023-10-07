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

var (
	defaultListKeyMap = list.DefaultKeyMap()
)

type Model struct {
	*list.Model
	width         int
	height        int
	editable      bool
	shortHelpKeys []key.Binding
	fullHelpKeys  []key.Binding

	toggledItems map[int]struct{}
	limit        int

	items Items
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
func New(items Items, opts ...Option) *Model {
	w, h := util.TermSize()
	m := &Model{
		width:        w,
		height:       h,
		items:        items,
		state:        Browsing,
		limit:        0,
		toggledItems: make(map[int]struct{}),
	}

	for _, opt := range opts {
		opt(m)
	}

	m.Model = m.NewListModel(items)

	if m.editable {
		m.AddShortHelpKeys(m.KeyMap.InsertItem, m.KeyMap.RemoveItem)
		m.AddFullHelpKeys(m.KeyMap.InsertItem, m.KeyMap.RemoveItem)
	}

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
func (m Model) NewListModel(items Items) *list.Model {
	var li []list.Item
	for _, i := range items.ParseFunc() {
		li = append(li, i)
	}
	l := list.New(li, m.items, m.width, m.height)
	return &l
}

// ChooseOne configures a list to return a single choice.
func ChooseOne(items Items, opts ...Option) *Model {
	m := New(items, opts...)
	m.ConfigureList(WithLimit(1))
	return m
}

// ChooseAny configures a list for multiple selections.
func ChooseAny(items Items, opts ...Option) *Model {
	m := New(items, opts...)
	m.ConfigureList(WithLimit(-1))
	return m
}

// ChooseSome configures a list for limited multiple selections.
func ChooseSome(items Items, limit int, opts ...Option) *Model {
	m := New(items, opts...)
	m.ConfigureList(WithLimit(limit))
	return m
}

// Edit configures an editable list: items are not selectable but can be
// removed from the list or new items entered with a prompt.
func Edit(items Items, opts ...Option) *Model {
	opts = append(opts, Editable())
	m := New(items, opts...)
	m.SetInput("Insert Item: ", InsertItem)

	m.ConfigureList(WithFiltering(false), WithLimit(0))

	return m
}

// Editable marks a list as editable
func Editable() Option {
	return func(m *Model) {
		m.editable = true
	}
}

// WithFiltering sets filtering on list.Model.
func WithFiltering(f bool) ListOption {
	return func(m *list.Model) {
		m.SetFilteringEnabled(f)
	}
}

// WithLimit sets the limit of choices for a selectable list.
func WithLimit(n int) ListOption {
	return func(m *list.Model) {
		m.SetLimit(n)
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Filter):
			//m.state = Input
		}
		if m.Browsing() && m.Selectable() {
			switch msg.Type {
			case tea.KeyEnter:
				if !m.Model.MultiSelectable() {
					m.Model.ToggleItem()
				}
				return m, tea.Quit
			}
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

// InsertItemMsg holds the title of the item to be inserted.
type InsertItemMsg struct {
	Value string
}

// InsertItem returns a tea.Cmd to insert an item into a list.
func InsertItem(val string) tea.Cmd {
	return func() tea.Msg {
		return InsertItemMsg{
			Value: val,
		}
	}
}

// RemoveItemMsg is a struct for the index to be removed.
type RemoveItemMsg struct {
	Index int
}

// RemoveItem returns a tea.Cmd for removing the item at index n.
func RemoveItem(idx int) tea.Cmd {
	return func() tea.Msg {
		return RemoveItemMsg{Index: idx}
	}
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
		if m.input.Focused() {
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
