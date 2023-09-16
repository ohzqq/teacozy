package list

import (
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
)

// Choices is a slice of string maps to satisfy the fuzzy.Source interface
type Choices []Choice

// Choice is a map[string]string for a single choice
type Choice map[string]string

func NewChoices(c any) (Choices, error) {
	return SliceToChoices(c)
	// var choices []map[string]string
}

func SliceToChoices[S ~[]E, E any](c S) (Choices, error) {
	choices := make([]Choice, len(c))
	for i, val := range c {
		m, err := cast.ToStringMapStringE(val)
		if err != nil {
			return choices, err
		}
		choices[i] = m
	}
	return choices, nil
}

// String satisfies the fuzzy.Source interface
func (c Choices) String(i int) string {
	return c[i].Value()
}

// Len satisfies the fuzzy.Source interface
func (c Choices) Len() int {
	return len(c)
}

// Filter fuzzy matches items in the list
func (c Choices) Filter(s string) fuzzy.Matches {
	var matches fuzzy.Matches

	m := fuzzy.FindFrom(s, c)
	if len(m) == 0 {
		//return ChoicesToItems(c)
		return matches
	}
	//for _, match := range m {
	//item := New()
	//item.SetMatch(match).SetLabel(c[match.Index].Key())
	//matches = append(matches, item)
	//}
	return matches
}

// Key returns the key or label
func (c Choice) Key() string {
	return maps.Keys(c)[0]
}

// Value returns the string value
func (c Choice) Value() string {
	return maps.Values(c)[0]
}
