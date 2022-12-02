package tui

import (
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/list"
)

type Item struct {
	*list.Item
	Editable bool
}

func ShowItemMeta(curItem *list.Item, items *list.Items, m *bubblelist.Model) tea.Cmd {
	if curItem.HasMeta() {
		return ShowItemInfoCmd(curItem)
	}
	return nil
}
