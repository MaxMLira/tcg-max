package domain

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// fakeShuffler is a deterministic Shuffler: it records how it was called and
// applies a fixed permutation (reversing the order) through swap.
type fakeShuffler struct {
	calls int
	gotN  int
}

func (f *fakeShuffler) Shuffle(n int, swap func(i, j int)) {
	f.calls++
	f.gotN = n
	for i := 0; i < n/2; i++ {
		swap(i, n-1-i)
	}
}

// noopShuffler leaves the expansion order untouched.
type noopShuffler struct{}

func (noopShuffler) Shuffle(int, func(i, j int)) {}

// newTestDeck builds a DeckZone from validDeck with the given shuffler.
func newTestDeck(t *testing.T, owner ParticipantID, s Shuffler) *DeckZone {
	t.Helper()

	deck, err := NewDeckZone(owner, validDeck(t), s)
	if err != nil {
		t.Fatalf("NewDeckZone() returned unexpected error: %v", err)
	}
	return deck
}

func TestNewDeckZone_BuildsThirtyDistinctInstancesOwnedByTheParticipant(t *testing.T) {
	// Arrange
	def := validDeck(t)

	// Act
	deck, err := NewDeckZone("p1", def, noopShuffler{})

	// Assert
	if err != nil {
		t.Fatalf("NewDeckZone() returned unexpected error: %v", err)
	}
	if deck.Len() != DeckSize {
		t.Fatalf("Len() = %d, want %d", deck.Len(), DeckSize)
	}
	seen := make(map[InstanceID]bool, DeckSize)
	perCard := make(map[CardID]int)
	for _, inst := range deck.Cards() {
		if seen[inst.ID()] {
			t.Errorf("duplicated instance id %q", inst.ID())
		}
		seen[inst.ID()] = true
		if inst.Owner() != "p1" {
			t.Errorf("instance %q owner = %q, want p1", inst.ID(), inst.Owner())
		}
		if inst.Location() != LocationDeck {
			t.Errorf("instance %q location = %s, want %s", inst.ID(), inst.Location(), LocationDeck)
		}
		perCard[inst.Card().ID()]++
	}
	for n := 1; n <= DeckSize; n++ {
		if id := InstanceID(fmt.Sprintf("p1#%d", n)); !seen[id] {
			t.Errorf("missing instance id %q", id)
		}
	}
	for _, entry := range def.Entries() {
		if got := perCard[entry.Card.ID()]; got != entry.Copies {
			t.Errorf("instances of %q = %d, want %d copies", entry.Card.ID(), got, entry.Copies)
		}
	}
}

func TestNewDeckZone_AssignsIdsInExpansionOrderBeforeShuffling(t *testing.T) {
	// Arrange
	def := validDeck(t)

	// Act
	deck, err := NewDeckZone("p1", def, noopShuffler{})

	// Assert
	if err != nil {
		t.Fatalf("NewDeckZone() returned unexpected error: %v", err)
	}
	cards := deck.Cards()
	// validDeck has 15 entries with two copies each: p1#1 and p1#2 are CARD-001.
	for i, want := range []CardID{cardID(1), cardID(1), cardID(2), cardID(2)} {
		if cards[i].Card().ID() != want {
			t.Errorf("card at top+%d = %q, want %q", i, cards[i].Card().ID(), want)
		}
		if wantID := InstanceID(fmt.Sprintf("p1#%d", i+1)); cards[i].ID() != wantID {
			t.Errorf("card at top+%d id = %q, want %q", i, cards[i].ID(), wantID)
		}
	}
}

func TestNewDeckZone_TwoCopiesOfACardAreDistinctInstances(t *testing.T) {
	// Arrange
	deck := newTestDeck(t, "p1", noopShuffler{})
	cards := deck.Cards()

	// Act
	first, second := cards[0], cards[1]

	// Assert
	if first.Card().ID() != second.Card().ID() {
		t.Fatalf("test setup: expected two copies of one card, got %q and %q", first.Card().ID(), second.Card().ID())
	}
	if first == second || first.ID() == second.ID() {
		t.Errorf("two copies must be distinct instances, got %q and %q", first.ID(), second.ID())
	}
	if first.Card() != second.Card() {
		t.Errorf("two copies must share the same definition")
	}
}

func TestNewDeckZone_CallsTheShufflerOnceWithTheDeckSizeAndAppliesItsPermutation(t *testing.T) {
	// Arrange
	fake := &fakeShuffler{}
	unshuffled := newTestDeck(t, "p1", noopShuffler{})

	// Act
	deck := newTestDeck(t, "p1", fake)

	// Assert
	if fake.calls != 1 {
		t.Errorf("Shuffle calls = %d, want exactly 1", fake.calls)
	}
	if fake.gotN != DeckSize {
		t.Errorf("Shuffle n = %d, want %d", fake.gotN, DeckSize)
	}
	want := ids(unshuffled.Cards())
	slices.Reverse(want)
	if got := ids(deck.Cards()); !slices.Equal(got, want) {
		t.Errorf("Cards() ids = %v, want the reversed expansion %v", got, want)
	}
}

func TestNewDeckZone_SeededGeneratorIsReproducibleAndKeepsTheSameCards(t *testing.T) {
	// Arrange
	seeded := func() Shuffler { return rand.New(rand.NewPCG(1, 2)) }

	// Act
	first := newTestDeck(t, "p1", seeded())
	second := newTestDeck(t, "p1", seeded())
	unshuffled := newTestDeck(t, "p1", noopShuffler{})

	// Assert
	firstIDs, secondIDs := ids(first.Cards()), ids(second.Cards())
	if !slices.Equal(firstIDs, secondIDs) {
		t.Errorf("same seed gave different orders:\n%v\n%v", firstIDs, secondIDs)
	}
	if slices.Equal(firstIDs, ids(unshuffled.Cards())) {
		t.Errorf("the seeded shuffle left the order unchanged")
	}
	sorted := slices.Clone(firstIDs)
	slices.Sort(sorted)
	wantSorted := ids(unshuffled.Cards())
	slices.Sort(wantSorted)
	if !slices.Equal(sorted, wantSorted) {
		t.Errorf("shuffled deck is not a permutation of the expanded deck")
	}
}

func TestNewDeckZone_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		owner    ParticipantID
		def      func(t *testing.T) DeckDefinition
		shuffler Shuffler
		wantErr  error
	}{
		{name: "zero deck definition", owner: "p1", def: func(*testing.T) DeckDefinition { return DeckDefinition{} }, shuffler: noopShuffler{}, wantErr: ErrDeckSize},
		{name: "blank owner", owner: "  ", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "empty owner", owner: "", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "tab owner", owner: "\t", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "spaces owner", owner: "   ", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "leading space owner", owner: " p1", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "trailing space owner", owner: "p1 ", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "surrounded owner", owner: " p1 ", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "tab and newline owner", owner: "\tp1\n", def: validDeck, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "nil shuffler", owner: "p1", def: validDeck, shuffler: nil, wantErr: ErrInvalidDeck},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			def := tt.def(t)

			// Act
			deck, err := NewDeckZone(tt.owner, def, tt.shuffler)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewDeckZone() error = %v, want %v", err, tt.wantErr)
			}
			if deck != nil {
				t.Errorf("NewDeckZone() = %v, want nil deck on error", deck)
			}
		})
	}
}

func TestNewDeckZone_AcceptsOwnersWithoutSurroundingWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		owner ParticipantID
	}{
		{name: "short id", owner: "p1"},
		{name: "hyphenated id", owner: "player-1"},
		{name: "internal space is accepted", owner: "player one"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			def := validDeck(t)

			// Act
			deck, err := NewDeckZone(tt.owner, def, noopShuffler{})

			// Assert
			if err != nil {
				t.Fatalf("NewDeckZone(%q) returned unexpected error: %v", tt.owner, err)
			}
			for i, inst := range deck.Cards() {
				if inst.Owner() != tt.owner {
					t.Errorf("instance %q owner = %q, want %q", inst.ID(), inst.Owner(), tt.owner)
				}
				wantID := InstanceID(fmt.Sprintf("%s#%d", tt.owner, i+1))
				if inst.ID() != wantID {
					t.Errorf("instance %d id = %q, want %q", i, inst.ID(), wantID)
				}
			}
		})
	}
}

func TestNewDeckZone_PropagatesDeckDefinitionErrors(t *testing.T) {
	// Arrange
	invalid := DeckDefinition{entries: []DeckEntry{
		{Card: deckCard(t, 1), Copies: DeckCopiesLimit + 1},
		{Card: deckCard(t, 2), Copies: 1},
	}}

	// Act
	deck, err := NewDeckZone("p1", invalid, noopShuffler{})

	// Assert
	if !errors.Is(err, ErrTooManyCopies) {
		t.Fatalf("NewDeckZone() error = %v, want ErrTooManyCopies", err)
	}
	if deck != nil {
		t.Errorf("NewDeckZone() = %v, want nil deck on error", deck)
	}
}

func TestDeckZone_DrawTakesTheTopCard(t *testing.T) {
	// Arrange
	deck := newTestDeck(t, "p1", noopShuffler{})
	top := deck.Cards()[0]

	// Act
	drawn, err := deck.draw()

	// Assert
	if err != nil {
		t.Fatalf("draw() returned unexpected error: %v", err)
	}
	if drawn != top {
		t.Errorf("draw() = %q, want the top card %q", drawn.ID(), top.ID())
	}
	if deck.Len() != DeckSize-1 {
		t.Errorf("Len() = %d, want %d", deck.Len(), DeckSize-1)
	}
	if deck.Cards()[0] == top {
		t.Errorf("the drawn card is still on top of the deck")
	}
}

func TestDeckZone_DrawFromEmptyDeckFailsAndStaysUnchanged(t *testing.T) {
	// Arrange
	deck := newTestDeck(t, "p1", noopShuffler{})
	for range DeckSize {
		if _, err := deck.draw(); err != nil {
			t.Fatalf("draw() returned unexpected error: %v", err)
		}
	}

	// Act
	drawn, err := deck.draw()

	// Assert
	if !errors.Is(err, ErrDeckEmpty) {
		t.Fatalf("draw() error = %v, want ErrDeckEmpty", err)
	}
	if drawn != nil {
		t.Errorf("draw() = %v, want nil on error", drawn)
	}
	if deck.Len() != 0 {
		t.Errorf("Len() = %d, want 0 (unchanged)", deck.Len())
	}
}

func TestDeckZone_CardsReturnsACopy(t *testing.T) {
	// Arrange
	deck := newTestDeck(t, "p1", noopShuffler{})
	before := ids(deck.Cards())

	// Act
	cards := deck.Cards()
	cards[0] = nil
	cards[1], cards[2] = cards[2], cards[1]

	// Assert
	if got := ids(deck.Cards()); !slices.Equal(got, before) {
		t.Errorf("Cards() = %v, want %v (mutating the copy must not affect the deck)", got, before)
	}
}
