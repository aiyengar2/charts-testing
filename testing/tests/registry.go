package tests

import (
	"fmt"
	"sort"
)

var (
	allTests = make(map[string]Test)
)

// Register registers a template with the given name.
// Intended to be called at program init time.
func Register(t Test) {
	if _, ok := allTests[t.ID]; ok {
		panic(fmt.Sprintf("duplicate test: %v", t.ID))
	}
	allTests[t.ID] = t
}

// Get gets a template by name, returning a boolean indicating whether it was found.
func Get(name string) (Test, bool) {
	t, ok := allTests[name]
	return t, ok
}

// List returns all known templates, sorted by name.
func List() []Test {
	out := make([]Test, 0, len(allTests))
	for _, t := range allTests {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}
