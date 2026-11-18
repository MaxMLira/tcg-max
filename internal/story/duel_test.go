package story

import (
	"errors"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

func TestNewDuel(t *testing.T) {
	character := testCharacter(t, "char-duel", "Duel")
	rewards := [3]*domain.Card{testCard(t, 10), testCard(t, 11), testCard(t, 12)}
	tests := []struct {
		name      string
		number    int
		character Character
		rewards   [3]*domain.Card
		wantErr   error
	}{
		{name: "valid duel number 1", number: 1, character: character, rewards: rewards},
		{name: "valid duel number 4", number: 4, character: character, rewards: rewards},
		{name: "number 0 rejected", number: 0, character: character, rewards: rewards, wantErr: ErrInvalidDuel},
		{name: "number 5 rejected", number: 5, character: character, rewards: rewards, wantErr: ErrInvalidDuel},
		{name: "missing character rejected", number: 1, character: Character{}, rewards: rewards, wantErr: ErrInvalidDuel},
		{
			name:      "nil reward rejected",
			number:    1,
			character: character,
			rewards:   [3]*domain.Card{rewards[0], nil, rewards[2]},
			wantErr:   ErrInvalidDuel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDuel(tt.number, tt.character, tt.rewards)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewDuel() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewDuel() unexpected error: %v", err)
			}
			if d.Number() != tt.number {
				t.Fatalf("Number() = %d, want %d", d.Number(), tt.number)
			}
			if d.Character().ID() != tt.character.ID() {
				t.Fatalf("Character().ID() = %q, want %q", d.Character().ID(), tt.character.ID())
			}
			if got := len(d.Rewards()); got != 3 {
				t.Fatalf("len(Rewards()) = %d, want 3", got)
			}
		})
	}
}

func TestDuelRewardsCopyIsolated(t *testing.T) {
	d := testDuel(t, 1, testCharacter(t, "char-r", "R"))
	got := d.Rewards()
	got[0] = nil
	if d.Rewards()[0] == nil {
		t.Fatal("mutating the returned Rewards() slice changed the duel")
	}
}
