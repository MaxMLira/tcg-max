package domain

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

// handOf builds a Hand holding n instances with ids "h#1".."h#n".
func handOf(t *testing.T, n int) *Hand {
	t.Helper()

	hand := newHand()
	for i := 1; i <= n; i++ {
		inst := instanceOf(t, InstanceID(fmt.Sprintf("h#%d", i)), deckCard(t, i))
		if err := hand.add(inst); err != nil {
			t.Fatalf("add() returned unexpected error: %v", err)
		}
	}
	return hand
}

// ids returns the instance ids in order, to compare zone contents compactly.
func ids(cards []*CardInstance) []InstanceID {
	out := make([]InstanceID, 0, len(cards))
	for _, c := range cards {
		out = append(out, c.ID())
	}
	return out
}

func TestHand_HoldsUpToTheHandLimit(t *testing.T) {
	// Arrange
	hand := newHand()

	// Act
	for i := 1; i <= HandLimit; i++ {
		err := hand.add(instanceOf(t, InstanceID(fmt.Sprintf("h#%d", i)), deckCard(t, i)))
		if err != nil {
			t.Fatalf("add() #%d returned unexpected error: %v", i, err)
		}
	}

	// Assert
	if HandLimit != 4 {
		t.Errorf("HandLimit = %d, want 4 (game-rules §5.2)", HandLimit)
	}
	if hand.Len() != HandLimit {
		t.Errorf("Len() = %d, want %d", hand.Len(), HandLimit)
	}
	if !hand.Full() {
		t.Errorf("Full() = false, want true")
	}
}

func TestHand_RejectsTheFifthCardAndStaysUnchanged(t *testing.T) {
	// Arrange
	hand := handOf(t, HandLimit)
	before := ids(hand.Cards())
	extra := instanceOf(t, "h#extra", deckCard(t, 9))

	// Act
	err := hand.add(extra)

	// Assert
	if !errors.Is(err, ErrHandFull) {
		t.Fatalf("add() error = %v, want ErrHandFull", err)
	}
	if got := ids(hand.Cards()); !slices.Equal(got, before) {
		t.Errorf("Cards() = %v, want %v (unchanged)", got, before)
	}
	if _, found := hand.Get("h#extra"); found {
		t.Errorf("rejected card must not be in the hand")
	}
}

func TestHand_KeepsInsertionOrder(t *testing.T) {
	// Arrange
	hand := handOf(t, 3)

	// Act
	got := ids(hand.Cards())

	// Assert
	want := []InstanceID{"h#1", "h#2", "h#3"}
	if !slices.Equal(got, want) {
		t.Errorf("Cards() = %v, want %v", got, want)
	}
}

func TestHand_GetFindsOnlyCardsInTheHand(t *testing.T) {
	// Arrange
	hand := handOf(t, 2)

	// Act
	found, ok := hand.Get("h#2")
	_, missing := hand.Get("h#9")

	// Assert
	if !ok || found.ID() != "h#2" {
		t.Errorf("Get(h#2) = %v, %v, want the instance", found, ok)
	}
	if missing {
		t.Errorf("Get(h#9) found a card that is not in the hand")
	}
}

func TestHand_RemoveTakesTheCardAndKeepsTheOrder(t *testing.T) {
	// Arrange
	hand := handOf(t, 3)

	// Act
	removed, err := hand.remove("h#2")

	// Assert
	if err != nil {
		t.Fatalf("remove() returned unexpected error: %v", err)
	}
	if removed.ID() != "h#2" {
		t.Errorf("remove() = %q, want h#2", removed.ID())
	}
	if got, want := ids(hand.Cards()), []InstanceID{"h#1", "h#3"}; !slices.Equal(got, want) {
		t.Errorf("Cards() = %v, want %v", got, want)
	}
}

func TestHand_RemoveOfAbsentCardFailsAndStaysUnchanged(t *testing.T) {
	// Arrange
	hand := handOf(t, 2)
	before := ids(hand.Cards())

	// Act
	removed, err := hand.remove("h#9")

	// Assert
	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("remove() error = %v, want ErrCardNotFound", err)
	}
	if removed != nil {
		t.Errorf("remove() = %v, want nil on error", removed)
	}
	if got := ids(hand.Cards()); !slices.Equal(got, before) {
		t.Errorf("Cards() = %v, want %v (unchanged)", got, before)
	}
}

func TestHand_CardsReturnsACopy(t *testing.T) {
	// Arrange
	hand := handOf(t, 2)
	before := ids(hand.Cards())

	// Act
	cards := hand.Cards()
	cards[0] = nil
	cards = append(cards[:0], cards[1:]...)
	_ = cards

	// Assert
	if got := ids(hand.Cards()); !slices.Equal(got, before) {
		t.Errorf("Cards() = %v, want %v (mutating the copy must not affect the hand)", got, before)
	}
}
