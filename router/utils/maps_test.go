package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mapStringToInt = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

func TestKeys(t *testing.T) {
	keys := Keys(mapStringToInt)
	assert.ElementsMatch(t, []string{"a", "b", "c", "d"}, keys)
}

func TestValues(t *testing.T) {
	values := Values(mapStringToInt)
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, values)
}

func TestEntries(t *testing.T) {
	entries := Entries(mapStringToInt)
	assert.ElementsMatch(t, []Entry[string, int]{{"a", 1}, {"b", 2}, {"c", 3}, {"d", 4}}, entries)
}

func TestFromEntries(t *testing.T) {
	entries := Entries(mapStringToInt)
	entries2 := []Entry[string, int]{{"a", 1}, {"b", 2}, {"c", 3}, {"d", 4}}
	assert.ElementsMatch(t, entries, entries2)
	assert.Equal(t, FromEntries(entries), FromEntries(entries2))
	assert.Equal(t, mapStringToInt, FromEntries(entries))
	assert.Equal(t, mapStringToInt, FromEntries(entries2))
}

func TestClone(t *testing.T) {
	clone := Clone(mapStringToInt)
	clone2 := Clone(mapStringToInt)
	assert.Equal(t, mapStringToInt, clone)
	assert.Equal(t, mapStringToInt, clone2)
	assert.Equal(t, clone, clone2)
	assert.NotSame(t, mapStringToInt, clone)
	assert.NotSame(t, mapStringToInt, clone2)
	assert.NotSame(t, clone, clone2)
}
