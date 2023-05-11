package state

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

func InputValue() string {
	return inputValue
}

func SetInputValue(val string) {
	inputValue = val
}
