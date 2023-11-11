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
			ori := VariableCellState(possibilites)
			removed := rapid.SliceOf(rapid.Int()).Draw(t, "removed")
			updated, _ := ori.WithRemovedPossibilities(removed...)
			assert.True(t, len(updated.Possibilities) <= len(ori.Possibilities))
		})
	})

	t.Run("HasValueReturnsTrue", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			value := rapid.Int().Draw(t, "value")
			cs := FixedCellState(value)
			removed := rapid.SliceOf(rapid.Int()).Draw(t, "removed")
			_, ok := cs.WithRemovedPossibilities(removed...)
			assert.True(t, ok)
		})
	})

	t.Run("RemoveTheLastValueReturnsFalse", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			value := rapid.Int().Draw(t, "value")
			cs := VariableCellState([]int{value})
			_, ok := cs.WithRemovedPossibilities(value)
			assert.False(t, ok)
		})
	})
}
