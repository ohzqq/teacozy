package match

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/sahilm/fuzzy"
	"golang.org/x/exp/maps"
)

type Component struct {
	reactea.BasicComponent // It implements AfterUpdate() for us, so we don't have to care!
	reactea.BasicPropfulComponent[Props]

	Matches fuzzy.Matches
	Style   style.ListItem
	*Prefix
}

type Props struct {
	Search   string
	Choices  []map[string]string
	Selected map[int]struct{}
	Cursor   int
	Matches  func([]string)
}

type Prefix struct {
	Cursor     string
	Selected   string
	Unselected string
}

func New() *Component {
	c := new(Component)
	c.Style = DefaultItemStyle()
	c.Prefix = DefaultPrefix()
	return c
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.Filter(props.Search)
	return nil
}

func (c *Component) Render(w, h int) string {
	var matches []string
	for _, match := range c.Matches {
		var s strings.Builder
		item := c.Props().Choices[match.Index]
		label := maps.Keys(item)[0]

		pre := "x"

		switch {
		case match.Index == c.Props().Cursor:
			pre = c.Style.Cursor.Render(pre)
		default:
			if _, ok := c.Props().Selected[match.Index]; ok {
				pre = c.Style.Selected.Render(pre)
			} else if label == "" {
				pre = strings.Repeat(" ", lipgloss.Width(pre))
			} else {
				pre = c.Style.Label.Render(pre)
			}
		}

		s.WriteString("[")
		s.WriteString(pre)
		s.WriteString("]")

		text := lipgloss.StyleRunes(
			match.Str,
			match.MatchedIndexes,
			c.Style.Match,
			c.Style.Text,
		)
		s.WriteString(lipgloss.NewStyle().Render(text))
		matches = append(matches, s.String())
	}
	c.Props().Matches(matches)

	return lipgloss.JoinVertical(lipgloss.Left, matches...)
}

func (i Props) String(idx int) string {
	var str string
	item := i.Choices[idx]
	for _, v := range item {
		str = v
	}
	return str
}

func (i Props) Len() int {
	return len(i.Choices)
}

func (c *Component) Filter(search string) {
	c.Matches = fuzzy.FindFrom(search, c.Props())
	if len(c.Matches) == 0 {
		for idx, i := range c.Props().Choices {
			for _, v := range i {
				m := fuzzy.Match{
					Str:   v,
					Index: idx,
				}
				c.Matches = append(c.Matches, m)
			}
		}
	}
}

func DefaultItemStyle() style.ListItem {
	var s style.ListItem
	s.Cursor = style.Cursor
	s.Selected = style.Selected
	s.Unselected = style.Unselected
	s.Text = style.Foreground
	s.Label = style.Label
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())
	return s
}

func DefaultPrefix() *Prefix {
	return &Prefix{
		Cursor:     CursorPrefix,
		Selected:   SelectedPrefix,
		Unselected: UnselectedPrefix,
	}
}

const (
	PromptPrefix     = "> "
	CursorPrefix     = "x"
	SelectedPrefix   = "◉ "
	UnselectedPrefix = " "
)
