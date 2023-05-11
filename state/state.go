package state

import tea "github.com/charmbracelet/bubbletea"

var (
	currentItem   int
	selectedItems = make(map[int]struct{})
	inputValue    string
)

func CurrentItem() int {
	return currentItem
}

func SetCurrentItem(c int) {
	currentItem = c
}

func SelectedItems() map[int]struct{} {
	return selectedItems
}

type InputValueMsg struct {
	Value string
}

func InputValue(val string) tea.Cmd {
	//return inputValue
	return func() tea.Msg {
		return InputValueMsg{Value: string}
	}
}

func SetInputValue(val string) {
	inputValue = val
}
