package domain

import (
	"errors"
	"reflect"
	"testing"
)

// endFullTurn ends the current player's turn: Phase 1 -> Phase 2 -> next
// player's Phase 1.
func endFullTurn(t *testing.T, d *Duel) {
	t.Helper()

	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() = %v", err)
	}
	if err := d.EndTurn(); err != nil {
		t.Fatalf("EndTurn() = %v", err)
	}
}

func TestNewDuel_StartsInPhase1OfTheFirstTurn(t *testing.T) {
	d := newTestDuel(t, noopShuffler{})

	turn := d.Turn()
	if turn.Number() != FirstTurnNumber || !turn.IsFirst() {
		t.Errorf("turn = %v, want number %d and IsFirst", turn.Number(), FirstTurnNumber)
	}
	if turn.Phase() != Phase1 {
		t.Errorf("phase = %v, want Phase 1", turn.Phase())
	}
	if turn.Participant() == "" {
		t.Error("the first turn has no participant")
	}
}

func TestEndPhase1_MovesToPhase2(t *testing.T) {
	d := newTestDuel(t, noopShuffler{})

	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() = %v", err)
	}
	if got := d.Turn().Phase(); got != Phase2 {
		t.Errorf("phase = %v, want Phase 2", got)
	}

	// A second EndPhase1 is rejected.
	before := duelStateOf(d)
	if err := d.EndPhase1(); !errors.Is(err, ErrWrongPhase) {
		t.Errorf("second EndPhase1() = %v, want ErrWrongPhase", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by the rejected EndPhase1")
	}
}

func TestEndTurn_AlternatesPlayersAndResetsToPhase1(t *testing.T) {
	d := newTestDuel(t, noopShuffler{})
	first := d.Turn().Participant()

	endFullTurn(t, d)

	turn := d.Turn()
	if turn.Number() != FirstTurnNumber+1 {
		t.Errorf("turn number = %d, want %d", turn.Number(), FirstTurnNumber+1)
	}
	if turn.Participant() == first {
		t.Error("the turn did not alternate to the other participant")
	}
	if turn.Phase() != Phase1 {
		t.Errorf("phase = %v, want Phase 1", turn.Phase())
	}
	if turn.IsFirst() {
		t.Error("the second turn reports IsFirst")
	}
}

func TestEndTurn_RejectedInPhase1(t *testing.T) {
	d := newTestDuel(t, noopShuffler{})
	before := duelStateOf(d)

	err := d.EndTurn()

	if !errors.Is(err, ErrWrongPhase) {
		t.Fatalf("EndTurn() in phase 1 = %v, want ErrWrongPhase", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by the rejected EndTurn")
	}
}
