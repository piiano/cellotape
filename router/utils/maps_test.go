package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
