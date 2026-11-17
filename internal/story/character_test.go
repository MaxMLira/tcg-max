package story

import (
	"errors"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

func TestNewCharacter(t *testing.T) {
	deck := testDeck(t)
	tests := []struct {
		name     string
		id       string
		charName string
		deck     domain.DeckDefinition
		wantErr  error
	}{
		{name: "valid character", id: "char-richard", charName: "Richard", deck: deck},
		{name: "blank id rejected", id: "   ", charName: "Richard", deck: deck, wantErr: ErrInvalidCharacter},
		{name: "blank name rejected", id: "char-richard", charName: "  ", deck: deck, wantErr: ErrInvalidCharacter},
		{name: "invalid deck rejected", id: "char-richard", charName: "Richard", deck: domain.DeckDefinition{}, wantErr: ErrInvalidCharacter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCharacter(tt.id, tt.charName, tt.deck)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewCharacter() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewCharacter() unexpected error: %v", err)
			}
			if c.ID() != tt.id || c.Name() != tt.charName {
				t.Fatalf("accessors = (%q, %q), want (%q, %q)", c.ID(), c.Name(), tt.id, tt.charName)
			}
			if got := c.Deck().TotalCards(); got != domain.DeckSize {
				t.Fatalf("Deck().TotalCards() = %d, want %d", got, domain.DeckSize)
			}
		})
	}
}

func TestNewCharacterInvalidDeckWrapsDomainError(t *testing.T) {
	_, err := NewCharacter("char-x", "X", domain.DeckDefinition{})
	if !errors.Is(err, ErrInvalidCharacter) {
		t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidCharacter)
	}
	if !errors.Is(err, domain.ErrDeckSize) {
		t.Fatalf("err = %v, want wrapped %v", err, domain.ErrDeckSize)
	}
}

func TestCharacterDeckEntriesCopyIsolated(t *testing.T) {
	c := testCharacter(t, "char-copy", "Copy")
	entries := c.Deck().Entries()
	entries[0].Copies = 99
	if got := c.Deck().Entries()[0].Copies; got != 2 {
		t.Fatalf("Entries()[0].Copies = %d after mutating the copy, want 2", got)
	}
}
