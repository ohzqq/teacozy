package style

import (
	"bytes"
	"log"
	"text/template"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/urbooks-core/cfg"
)

const (
	Bullet   = "•"
	Ellipsis = "…"
)

type ItemStyle struct {
	NormalItem   lipgloss.Style
	CurrentItem  lipgloss.Style
	SelectedItem lipgloss.Style
	SubItem      lipgloss.Style
}

func ItemStyles() (s ItemStyle) {
	s.NormalItem = lipgloss.NewStyle().Foreground(cfg.Tui().Colors().DefaultFg)
	s.CurrentItem = lipgloss.NewStyle().Foreground(cfg.Tui().Colors().Green).Reverse(true)
	s.SelectedItem = lipgloss.NewStyle().Foreground(cfg.Tui().Colors().Grey)
	s.SubItem = lipgloss.NewStyle().Foreground(cfg.Tui().Colors().Purple)
	return s
}

func FrameStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true).
		MarginRight(0)
	return s
}

func ListStyles() (s list.Styles) {
	verySubduedColor := cfg.Tui().Colors().Grey
	subduedColor := cfg.Tui().Colors().White

	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 0, 0)

	s.Title = lipgloss.NewStyle().
		Background(cfg.Tui().Colors().Purple).
		Foreground(cfg.Tui().Colors().Black).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Cyan)

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Pink)

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Yellow)

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Blue).
		Padding(0, 0, 1, 2)

	s.StatusEmpty = lipgloss.NewStyle().Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Purple)

	s.StatusBarFilterCount = lipgloss.NewStyle().Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Grey)

	s.ArabicPagination = lipgloss.NewStyle().Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(cfg.Tui().Colors().Pink).
		SetString(Bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(Bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + Bullet + " ")

	return s
}

func helpStyles() help.Styles {
	styles := help.Styles{}

	keyStyle := lipgloss.NewStyle().PaddingRight(1).Foreground(cfg.Tui().Colors().Green)
	descStyle := lipgloss.NewStyle().PaddingRight(1).Foreground(cfg.Tui().Colors().Blue)
	sepStyle := lipgloss.NewStyle().Foreground(cfg.Tui().Colors().Pink)

	styles.ShortKey = keyStyle
	styles.ShortDesc = descStyle
	styles.ShortSeparator = sepStyle
	styles.FullKey = keyStyle.Copy()
	styles.FullDesc = descStyle.Copy()
	styles.FullSeparator = sepStyle.Copy()
	styles.Ellipsis = sepStyle.Copy()

	return styles
}

func RenderMarkdown(md string) string {
	var (
		metaStyle = template.Must(template.New("mdStyle").Parse(styleTmpl))
		style     bytes.Buffer
	)

	err := metaStyle.Execute(&style, cfg.Tui().Colors())
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(style.Bytes()),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		log.Fatal(err)
	}

	str, err := renderer.Render(md)
	if err != nil {
		log.Fatal(err)
	}

	return str
}

var styleTmpl = `
{
  "document": {
    "block_prefix": "",
    "block_suffix": "",
    "color": "{{.DefaultFg}}",
    "background_color": "{{.DefaultBg}}",
    "margin": 0
  },
  "block_quote": {
    "indent": 1,
    "indent_token": "│ "
  },
  "paragraph": {
    "block_suffix": ""
  },
  "list": {
    "level_indent": 2
  },
  "heading": {
    "block_suffix": "",
    "color": "{{.Pink}}",
    "bold": true
  },
  "h1": {
    "prefix": " ",
    "suffix": " ",
    "color": "{{.DefaultBg}}",
    "background_color": "{{.Blue}}",
    "bold": true
  },
  "h2": {
    "prefix": "## "
  },
  "h3": {
    "prefix": "### "
  },
  "h4": {
    "prefix": "#### "
  },
  "h5": {
    "prefix": "##### "
  },
  "h6": {
    "prefix": "###### ",
    "bold": false
  },
  "text": {},
  "strikethrough": {
    "crossed_out": true
  },
  "emph": {
    "italic": true
  },
  "strong": {
    "color": "{{.Cyan}}",
    "bold": true
  },
  "hr": {
    "color": "{{.Pink}}",
    "format": "\n--------\n"
  },
  "item": {
    "block_prefix": "• "
  },
  "enumeration": {
    "block_prefix": ". "
  },
  "html_block": {},
  "html_span": {}
}
`
