package item

type Items struct {
	all []Item
}

func NewItems() Items {
	return Items{}
}

func (i *Items) Add(item Item) *Items {
	i.all = append(i.all, item)
	return i
}

func (i Items) All() []Item {
	var items []Item
	for _, item := range i.all {
		items = append(items, item)
		items = append(items, item.Flatten()...)
	}
	i.all = items
	return i.all
}
