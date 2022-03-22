package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSet(t *testing.T) {
	baseSet := NewSet(1, 5, 1, 4, 2, 2, 1, 3)
	baseLength := len(baseSet)
	require.Equal(t, baseLength, 5)

	t.Run("clone a set", func(t *testing.T) {
		clone := baseSet.Clone()
		assert.Len(t, clone, baseLength)
		assert.Equal(t, baseSet, clone)
		assert.NotSame(t, baseSet, clone)
	})
	t.Run("adding new element", func(t *testing.T) {
		set := baseSet.Clone()
		added := set.Add(6)
		assert.Len(t, set, baseLength+1)
		assert.True(t, added)
	})
	t.Run("adding existing element", func(t *testing.T) {
		set := baseSet.Clone()
		added := set.Add(4)
		assert.Len(t, set, baseLength)
		assert.False(t, added)
	})
	t.Run("removing an element", func(t *testing.T) {
		set := baseSet.Clone()
		removed := set.Remove(5)
		assert.Len(t, set, baseLength-1)
		assert.True(t, removed)
	})
	t.Run("removing non existing element", func(t *testing.T) {
		set := baseSet.Clone()
		removed := set.Remove(7)
		assert.Len(t, set, baseLength)
		assert.False(t, removed)
	})
	t.Run("has element", func(t *testing.T) {
		set := baseSet.Clone()
		has3 := set.Has(3)
		assert.Len(t, set, baseLength)
		assert.True(t, has3)
	})
	t.Run("has for non existing element", func(t *testing.T) {
		set := baseSet.Clone()
		has8 := set.Has(8)
		assert.Len(t, set, baseLength)
		assert.False(t, has8)
	})

}
