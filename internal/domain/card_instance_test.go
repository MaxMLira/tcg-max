package domain

import (
	"errors"
	"strings"
	"testing"
)

// instanceOf builds a card instance owned by "p1" or fails the test.
func instanceOf(t *testing.T, id InstanceID, card *Card) *CardInstance {
	t.Helper()

	inst, err := newCardInstance(id, card, "p1")
	if err != nil {
		t.Fatalf("newCardInstance() returned unexpected error: %v", err)
	}
	return inst
}

func TestNewCardInstance_StartsInDeckWithIdentityAndOwner(t *testing.T) {
	// Arrange
	card := deckCard(t, 1)

	// Act
	inst, err := newCardInstance("p1#1", card, "p1")

	// Assert
	if err != nil {
		t.Fatalf("newCardInstance() returned unexpected error: %v", err)
	}
	if inst.ID() != "p1#1" {
		t.Errorf("ID() = %q, want %q", inst.ID(), "p1#1")
	}
	if inst.Owner() != "p1" {
		t.Errorf("Owner() = %q, want %q", inst.Owner(), "p1")
	}
	if inst.Card() != card {
		t.Errorf("Card() returned a different definition pointer")
	}
	if inst.Location() != LocationDeck {
		t.Errorf("Location() = %s, want %s", inst.Location(), LocationDeck)
	}
}

func TestNewCardInstance_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		id    InstanceID
		card  func(t *testing.T) *Card
		owner ParticipantID
	}{
		{name: "nil card", id: "p1#1", card: func(*testing.T) *Card { return nil }, owner: "p1"},
		{name: "blank id", id: "  ", card: func(t *testing.T) *Card { return deckCard(t, 1) }, owner: "p1"},
		{name: "empty id", id: "", card: func(t *testing.T) *Card { return deckCard(t, 1) }, owner: "p1"},
		{name: "blank owner", id: "p1#1", card: func(t *testing.T) *Card { return deckCard(t, 1) }, owner: " "},
		{name: "empty owner", id: "p1#1", card: func(t *testing.T) *Card { return deckCard(t, 1) }, owner: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			card := tt.card(t)

			// Act
			inst, err := newCardInstance(tt.id, card, tt.owner)

			// Assert
			if !errors.Is(err, ErrInvalidInstance) {
				t.Fatalf("newCardInstance() error = %v, want ErrInvalidInstance", err)
			}
			if inst != nil {
				t.Errorf("newCardInstance() = %v, want nil instance on error", inst)
			}
		})
	}
}

func TestCardInstance_TwoCopiesOfACardAreDistinctInstances(t *testing.T) {
	// Arrange
	card := deckCard(t, 1)

	// Act
	first := instanceOf(t, "p1#1", card)
	second := instanceOf(t, "p1#2", card)

	// Assert
	if first == second || first.ID() == second.ID() {
		t.Errorf("two copies must be distinct instances, got %v and %v", first.ID(), second.ID())
	}
	if first.Card() != second.Card() {
		t.Errorf("two copies must share the same definition")
	}
}

func TestCardInstance_SetLocationChangesOnlyTheLocation(t *testing.T) {
	// Arrange
	card := deckCard(t, 1)
	inst := instanceOf(t, "p1#1", card)

	// Act
	inst.setLocation(LocationHand)

	// Assert
	if inst.Location() != LocationHand {
		t.Errorf("Location() = %s, want %s", inst.Location(), LocationHand)
	}
	if inst.ID() != "p1#1" || inst.Owner() != "p1" || inst.Card() != card {
		t.Errorf("setLocation changed identity, owner or definition")
	}
}

func TestLocation_StringNamesEachZone(t *testing.T) {
	tests := []struct {
		location Location
		want     string
	}{
		{LocationDeck, "Deck"},
		{LocationHand, "Hand"},
		{LocationField, "Field"},
		{LocationGraveyard, "Graveyard"},
		{Location(0), "Location(0)"},
		{Location(99), "Location(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			// Act
			got := tt.location.String()

			// Assert
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateParticipantID(t *testing.T) {
	tests := []struct {
		name    string
		id      ParticipantID
		wantErr bool
	}{
		{name: "empty", id: "", wantErr: true},
		{name: "spaces", id: "   ", wantErr: true},
		{name: "tab", id: "\t", wantErr: true},
		{name: "leading space", id: " p1", wantErr: true},
		{name: "trailing space", id: "p1 ", wantErr: true},
		{name: "surrounded by spaces", id: " p1 ", wantErr: true},
		{name: "surrounded by tab and newline", id: "\tp1\n", wantErr: true},
		{name: "short id", id: "p1"},
		{name: "hyphenated id", id: "player-1"},
		{name: "internal space is accepted", id: "player one"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (the table entry)

			// Act
			err := validateParticipantID(tt.id)

			// Assert
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidParticipant) {
					t.Fatalf("validateParticipantID(%q) error = %v, want %v", tt.id, err, ErrInvalidParticipant)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateParticipantID(%q) returned unexpected error: %v", tt.id, err)
			}
		})
	}
}

func TestValidateParticipantID_MessageDistinguishesBlankFromPadded(t *testing.T) {
	// Arrange
	blank, padded := ParticipantID("  "), ParticipantID(" p1")

	// Act
	blankErr := validateParticipantID(blank)
	paddedErr := validateParticipantID(padded)

	// Assert
	if blankErr == nil || !strings.Contains(blankErr.Error(), "must not be blank") {
		t.Errorf("blank id error = %v, want it to mention \"must not be blank\"", blankErr)
	}
	if paddedErr == nil || !strings.Contains(paddedErr.Error(), "surrounding whitespace") || !strings.Contains(paddedErr.Error(), `" p1"`) {
		t.Errorf("padded id error = %v, want it to mention the surrounding whitespace and the quoted id", paddedErr)
	}
}
