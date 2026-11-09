package domain

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// The standard generator satisfies Shuffler, so the composition root can inject
// it without an adapter.
var _ Shuffler = (*rand.Rand)(nil)

func TestShuffler_StandardGeneratorPermutesThroughSwap(t *testing.T) {
	// Arrange
	var shuffler Shuffler = rand.New(rand.NewPCG(1, 2))
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// Act
	shuffler.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })

	// Assert
	sorted := slices.Clone(items)
	slices.Sort(sorted)
	if want := []int{1, 2, 3, 4, 5, 6, 7, 8}; !slices.Equal(sorted, want) {
		t.Errorf("shuffled items = %v, want a permutation of %v", items, want)
	}
}
