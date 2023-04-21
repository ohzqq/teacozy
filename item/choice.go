package item

import (
	"fmt"

	"github.com/sahilm/fuzzy"
	"golang.org/x/exp/maps"
)

// Choices is a slice of string maps to satisfy the fuzzy.Source interface
type Choices []Choice

// Choice is a map[string]string for a single choice
type Choice map[string]string

// String satisfies the fuzzy.Source interface
func (c Choices) String(i int) string {
	return c[i].String()
}

func (c Choices) Label(i int) string {
	return c[i].Label()
}

// Len satisfies the fuzzy.Source interface
func (c Choices) Len() int {
	return len(c)
}

// Filter fuzzy matches items in the list
func (c Choices) Filter(s string) []Item {
	matches := []Item{}
	m := fuzzy.FindFrom(s, c)
	if len(m) == 0 {
		return ChoicesToItems(c)
	}
	for _, match := range m {
		item := New()
		item.SetMatch(match).SetLabel(c[match.Index].Label())
		matches = append(matches, item)
	}
	return matches
}

func (c Choices) Find(s string) fuzzy.Matches {
	m := fuzzy.FindFrom(s, c)
	return m
}

// Set sets the value for an item in the slice
func (c Choices) Set(idx int, val string) {
	c[idx] = c[idx].Set(val)
}

// Label returns the key or label
func (c Choice) Label() string {
	return maps.Keys(c)[0]
}

// String returns the string value
func (c Choice) String() string {
	return maps.Values(c)[0]
}

// Set sets the value
func (c Choice) Set(v string) Choice {
	for k, _ := range c {
		c[k] = v
		break
	}
	return c
}

// MapToChoices converts a slice of maps to Choices. Only the first item in the map
// is collected.
func MapToChoices[K comparable, V any, M ~map[K]V](c []M) Choices {
	choices := make(Choices, len(c))
	for i, ch := range c {
		choices[i] = stringifyMap(ch)
	}
	return choices
}

func stringifyMap[K comparable, V any, M ~map[K]V](c M) Choice {
	k := fmt.Sprint(maps.Keys(c)[0])
	v := fmt.Sprint(maps.Values(c)[0])
	return Choice{k: v}
}

// SliceToChoices converts a generic slice to Choices. Values are converted to
// a string using fmt.Sprint and the key is the zero value.
func SliceToChoices[E any](c []E) Choices {
	choices := make([]Choice, len(c))
	for i, val := range c {
		choices[i] = Choice{"": fmt.Sprint(val)}
	}
	return choices
}
