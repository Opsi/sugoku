package backtrack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func TestRemovePossibilities(t *testing.T) {
	t.Run("PossibilitesCantGetBigger", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			possibilites := rapid.SliceOf(rapid.Int()).Draw(t, "possibilities")
			cs := &cellState{Possibilities: possibilites}
			removed := rapid.SliceOf(rapid.Int()).Draw(t, "removed")
			cs.RemovePossibilities(removed...)
			assert.True(t, len(cs.Possibilities) <= len(possibilites))
		})
	})

	t.Run("HasValueReturnsTrue", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			value := rapid.Int().Draw(t, "value")
			cs := &cellState{HasValue: true, Value: value}
			removed := rapid.SliceOf(rapid.Int()).Draw(t, "removed")
			assert.True(t, cs.RemovePossibilities(removed...))
		})
	})

	t.Run("RemoveTheLastValueReturnsFalse", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			value := rapid.Int().Draw(t, "value")
			cs := &cellState{Possibilities: []int{value}}
			assert.False(t, cs.RemovePossibilities(value))
		})
	})
}
