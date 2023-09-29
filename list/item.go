package list

type Items func() []*Item

type Item struct {
	title       string
	desc        string
	filterValue string
}

type ItemOpt func(*Item)

func ItemsStringSlice(items []string) Items {
	return func() []*Item {
		var li []*Item
		for _, item := range items {
			li = append(li, NewItem(item))
		}
		return li
	}
}

func ItemsMapSlice(items []map[string]string) Items {
	return func() []*Item {
		var li []*Item
		for _, item := range items {
			for k, v := range item {
				li = append(li, NewItem(k, Description(v)))
			}
		}
		return li
	}
}

func ItemsMap(items map[string]string) Items {
	return func() []*Item {
		var li []*Item
		for k, v := range items {
			li = append(li, NewItem(k, Description(v)))
		}
		return li
	}
}

func NewItem(title string, opts ...ItemOpt) *Item {
	item := &Item{
		title:       title,
		desc:        title,
		filterValue: title,
	}

	for _, opt := range opts {
		opt(item)
	}

	return item
}

func Description(desc string) ItemOpt {
	return func(i *Item) {
		i.desc = desc
	}
}

func FilterValue(val string) ItemOpt {
	return func(i *Item) {
		i.filterValue = val
	}
}

func (i Item) FilterValue() string { return i.filterValue }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
