package choose

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

var (
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

// Options is the customization options for the choose command.
type Options struct {
	Options           []map[string]string
	Limit             int
	NoLimit           bool
	Ordered           bool
	Height            int
	Cursor            string
	CursorPrefix      string
	SelectedPrefix    string
	UnselectedPrefix  string
	CursorStyle       lipgloss.Style
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
}

func New(o Options) *Model {
	tm := Model{
		Options: o,
		KeyMap:  KeyMap,
	}
	tm.Cursor = "> "
	tm.SelectedPrefix = "◉ "
	tm.UnselectedPrefix = "○ "
	tm.CursorPrefix = "○ "

	tm.CursorStyle = lipgloss.NewStyle().Foreground(color.Green)
	tm.ItemStyle = lipgloss.NewStyle().Foreground(color.Foreground)
	tm.SelectedItemStyle = lipgloss.NewStyle().Foreground(color.Grey)

	if tm.Height == 0 {
		tm.Height = 10
	}

	tm.Items = make([]Item, len(o.Options))

	for i, thing := range o.Options {
		for k, option := range thing {
			tm.Items[i] = Item{
				Key:      k,
				Text:     option,
				Selected: false,
				Order:    i,
			}
		}
	}

	if len(tm.Items) == 1 {
		tm.Limit = 1
	}

	// We don't need to display prefixes if we are only picking one option.
	// Simply displaying the cursor is enough.
	if tm.Limit == 1 && !o.NoLimit {
		tm.SelectedPrefix = ""
		tm.UnselectedPrefix = ""
		tm.CursorPrefix = ""
	}

	// If we've set no limit then we can simply select as many options as there
	// are so let's set the limit to the number of options.
	if o.NoLimit {
		tm.Limit = len(o.Options)
	}

	// Use the pagination model to display the current and total number of
	// pages.
	pager := paginator.New()
	pager.SetTotalPages((len(tm.Items) + tm.Height - 1) / tm.Height)
	pager.PerPage = tm.Height
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")

	// Disable Keybindings since we will control it ourselves.
	pager.UseHLKeys = false
	pager.UseLeftRightKeys = false
	pager.UseJKKeys = false
	pager.UsePgUpPgDownKeys = false

	tm.paginator = pager

	return &tm
}
