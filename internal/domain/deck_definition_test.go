package domain

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// cardID derives a distinct card identity from a sequence number, so test
// decks can be described compactly.
func cardID(n int) CardID {
	return CardID(fmt.Sprintf("CARD-%03d", n))
}

// deckCard builds a valid Monster definition for deck composition tests.
func deckCard(t *testing.T, n int) *Card {
	t.Helper()

	return mustCard(t, cardID(n), fmt.Sprintf("Card %d", n), CardTypeMonster, 1000+n, 1000, 1+n%5, SubtypeWarrior())
}

// singleEntries builds one entry per card with one copy each.
func singleEntries(t *testing.T, from, to int) []DeckEntry {
	t.Helper()

	entries := make([]DeckEntry, 0, to-from+1)
	for n := from; n <= to; n++ {
		entries = append(entries, DeckEntry{Card: deckCard(t, n), Copies: 1})
	}
	return entries
}

// validDeck builds a composition with exactly DeckSize cards: 15 cards with two
// copies each.
func validDeck(t *testing.T) DeckDefinition {
	t.Helper()

	entries := make([]DeckEntry, 0, DeckSize/DeckCopiesLimit)
	for n := 1; n <= DeckSize/DeckCopiesLimit; n++ {
		entries = append(entries, DeckEntry{Card: deckCard(t, n), Copies: DeckCopiesLimit})
	}

	deck, err := NewDeckDefinition(entries)
	if err != nil {
		t.Fatalf("validDeck is not valid: %v", err)
	}
	return deck
}

func TestDeckDefinition_AcceptsExactlyThirtyCards(t *testing.T) {
	// Arrange
	deck := validDeck(t)

	// Act
	err := deck.Validate()

	// Assert
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
	if got := deck.TotalCards(); got != DeckSize {
		t.Errorf("TotalCards() = %d, want %d", got, DeckSize)
	}
}

func TestDeckDefinition_AcceptsMixedCopyCounts(t *testing.T) {
	// Arrange
	entries := singleEntries(t, 1, 28)
	entries = append(entries, DeckEntry{Card: deckCard(t, 29), Copies: 2})

	// Act
	deck, err := NewDeckDefinition(entries)

	// Assert
	if err != nil {
		t.Fatalf("NewDeckDefinition() returned unexpected error: %v", err)
	}
	if got := deck.TotalCards(); got != DeckSize {
		t.Errorf("TotalCards() = %d, want %d", got, DeckSize)
	}
}

func TestDeckDefinition_RejectsInvalidComposition(t *testing.T) {
	tests := []struct {
		name    string
		entries func(t *testing.T) []DeckEntry
		wantErr error
	}{
		{
			name:    "fewer than thirty cards",
			entries: func(t *testing.T) []DeckEntry { return singleEntries(t, 1, 29) },
			wantErr: ErrDeckSize,
		},
		{
			name:    "more than thirty cards",
			entries: func(t *testing.T) []DeckEntry { return singleEntries(t, 1, 31) },
			wantErr: ErrDeckSize,
		},
		{
			name: "single entry with too many copies",
			entries: func(t *testing.T) []DeckEntry {
				entries := singleEntries(t, 1, 28)
				return append(entries, DeckEntry{Card: deckCard(t, 29), Copies: DeckCopiesLimit + 1})
			},
			wantErr: ErrTooManyCopies,
		},
		{
			name: "single entry without copies",
			entries: func(t *testing.T) []DeckEntry {
				entries := singleEntries(t, 1, 29)
				return append(entries, DeckEntry{Card: deckCard(t, 30), Copies: 0})
			},
			wantErr: ErrInvalidDeck,
		},
		{
			name: "single entry with negative copies",
			entries: func(t *testing.T) []DeckEntry {
				entries := singleEntries(t, 1, 29)
				return append(entries, DeckEntry{Card: deckCard(t, 30), Copies: -1})
			},
			wantErr: ErrInvalidDeck,
		},
		{
			name: "entry without a card",
			entries: func(t *testing.T) []DeckEntry {
				entries := singleEntries(t, 1, 29)
				return append(entries, DeckEntry{Copies: 1})
			},
			wantErr: ErrInvalidDeck,
		},
		{
			name: "same card in two entries",
			entries: func(t *testing.T) []DeckEntry {
				entries := singleEntries(t, 1, 28)
				return append(entries, DeckEntry{Card: deckCard(t, 1), Copies: 1}, DeckEntry{Card: deckCard(t, 2), Copies: 1})
			},
			wantErr: ErrDuplicateCard,
		},
		{
			name: "duplicate entries that would sum to the copy limit",
			entries: func(t *testing.T) []DeckEntry {
				entries := singleEntries(t, 1, 29)
				return append(entries, DeckEntry{Card: deckCard(t, 1), Copies: 1})
			},
			wantErr: ErrDuplicateCard,
		},
		{
			name:    "empty composition",
			entries: func(t *testing.T) []DeckEntry { return nil },
			wantErr: ErrDeckSize,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			entries := tc.entries(t)

			// Act
			deck, err := NewDeckDefinition(entries)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if deck.TotalCards() != 0 {
				t.Errorf("invalid deck was handed out with %d cards", deck.TotalCards())
			}
		})
	}
}

func TestDeckDefinition_ValidateDoesNotModifyTheComposition(t *testing.T) {
	// Arrange
	deck := validDeck(t)
	before := deck.Entries()

	// Act
	err := deck.Validate()

	// Assert
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
	after := deck.Entries()
	if len(after) != len(before) {
		t.Fatalf("len(Entries) = %d, want %d", len(after), len(before))
	}
	for i := range before {
		if after[i].Card != before[i].Card || after[i].Copies != before[i].Copies {
			t.Errorf("entry %d changed: %+v, want %+v", i, after[i], before[i])
		}
	}
}

func TestDeckDefinition_ErrorDescribesTheViolation(t *testing.T) {
	// Arrange
	entries := singleEntries(t, 1, 10)

	// Act
	_, err := NewDeckDefinition(entries)

	// Assert
	if !errors.Is(err, ErrDeckSize) {
		t.Fatalf("error = %v, want ErrDeckSize", err)
	}
	if got := err.Error(); !strings.Contains(got, "10") {
		t.Errorf("error message %q does not report the actual total 10", got)
	}
}

func TestNewDeckDefinition_RejectsNilAndEmptyEntries(t *testing.T) {
	tests := []struct {
		name    string
		entries []DeckEntry
	}{
		{name: "nil entries", entries: nil},
		{name: "empty entries", entries: []DeckEntry{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			_, err := NewDeckDefinition(tc.entries)

			// Assert
			if !errors.Is(err, ErrDeckSize) {
				t.Fatalf("error = %v, want ErrDeckSize", err)
			}
		})
	}
}

func TestDeckDefinition_ZeroValueIsInvalid(t *testing.T) {
	// Arrange
	var deck DeckDefinition

	// Act
	err := deck.Validate()

	// Assert
	if !errors.Is(err, ErrDeckSize) {
		t.Fatalf("error = %v, want ErrDeckSize", err)
	}
}

func TestNewDeckDefinition_CopiesTheInputEntries(t *testing.T) {
	// Arrange
	entries := validDeck(t).Entries()
	deck, err := NewDeckDefinition(entries)
	if err != nil {
		t.Fatalf("NewDeckDefinition() returned unexpected error: %v", err)
	}

	// Act
	entries[0].Copies = 99
	entries[1] = DeckEntry{}

	// Assert
	if err := deck.Validate(); err != nil {
		t.Fatalf("deck changed after input mutation: %v", err)
	}
	if got := deck.TotalCards(); got != DeckSize {
		t.Errorf("TotalCards() = %d, want %d", got, DeckSize)
	}
}

func TestDeckDefinition_EntriesReturnsACopy(t *testing.T) {
	// Arrange
	deck := validDeck(t)

	// Act
	returned := deck.Entries()
	returned[0].Copies = 99
	returned[1] = DeckEntry{}

	// Assert
	if err := deck.Validate(); err != nil {
		t.Fatalf("deck changed after mutating Entries(): %v", err)
	}
	if got := deck.Entries()[0].Copies; got != DeckCopiesLimit {
		t.Errorf("Entries()[0].Copies = %d, want %d", got, DeckCopiesLimit)
	}
}
