package domain

import (
	"fmt"
	"slices"
	"testing"
)

func TestGraveyard_KeepsArrivalOrder(t *testing.T) {
	// Arrange
	graveyard := newGraveyard()
	a := instanceOf(t, "g#A", deckCard(t, 1))
	b := instanceOf(t, "g#B", deckCard(t, 2))

	// Act
	graveyard.add(a)
	graveyard.add(b)

	// Assert
	if got, want := ids(graveyard.Cards()), []InstanceID{"g#A", "g#B"}; !slices.Equal(got, want) {
		t.Errorf("Cards() = %v, want %v (arrival order, game-rules §31.1)", got, want)
	}
	if graveyard.Len() != 2 {
		t.Errorf("Len() = %d, want 2", graveyard.Len())
	}
}

func TestGraveyard_HasNoCapacityLimit(t *testing.T) {
	// Arrange
	graveyard := newGraveyard()
	const total = 30

	// Act
	for i := 1; i <= total; i++ {
		graveyard.add(instanceOf(t, InstanceID(fmt.Sprintf("g#%d", i)), deckCard(t, i)))
	}

	// Assert
	if graveyard.Len() != total {
		t.Errorf("Len() = %d, want %d", graveyard.Len(), total)
	}
}

func TestGraveyard_GetFindsOnlyCardsInTheGraveyard(t *testing.T) {
	// Arrange
	graveyard := newGraveyard()
	a := instanceOf(t, "g#A", deckCard(t, 1))
	graveyard.add(a)

	// Act
	found, ok := graveyard.Get("g#A")
	_, missing := graveyard.Get("g#Z")

	// Assert
	if !ok || found != a {
		t.Errorf("Get(g#A) = %v, %v, want the instance", found, ok)
	}
	if missing {
		t.Errorf("Get(g#Z) found a card that is not in the graveyard")
	}
}

func TestGraveyard_CardsReturnsACopy(t *testing.T) {
	// Arrange
	graveyard := newGraveyard()
	graveyard.add(instanceOf(t, "g#A", deckCard(t, 1)))
	graveyard.add(instanceOf(t, "g#B", deckCard(t, 2)))

	// Act
	cards := graveyard.Cards()
	cards[0] = nil
	cards = cards[:1]
	_ = cards

	// Assert
	if got, want := ids(graveyard.Cards()), []InstanceID{"g#A", "g#B"}; !slices.Equal(got, want) {
		t.Errorf("Cards() = %v, want %v (mutating the copy must not affect the graveyard)", got, want)
	}
}
