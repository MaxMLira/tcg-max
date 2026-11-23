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

func TestEndTurn_ResetsThePlacementOfTheNextTurn(t *testing.T) {
	// Arrange — p1 used its placement on turn 1.
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	if err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp); err != nil {
		t.Fatalf("Place() = %v", err)
	}
	// p1's turn ends; two turns later it is p1's turn again.
	endFullTurn(t, d) // -> p2, turn 2
	endFullTurn(t, d) // -> p1, turn 3

	// Act — p1 may place again on its new turn.
	p1, _ := d.Participant("p1")
	err := d.Place([]InstanceID{p1.Hand().Cards()[0].ID()}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() on p1's second turn = %v, want success", err)
	}
}

func TestEndTurn_ResetsTheMonstersTurnState(t *testing.T) {
	// Arrange — p1 changed its monster's position on turn 1.
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	if err := d.ChangePosition(id, PositionDefense, FaceUp); err != nil {
		t.Fatalf("ChangePosition() = %v", err)
	}
	endFullTurn(t, d) // -> p2, turn 2
	endFullTurn(t, d) // -> p1, turn 3

	// Act — the same monster may change position again on p1's new turn.
	err := d.ChangePosition(id, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("ChangePosition() on the new turn = %v, want success", err)
	}
}

func TestEndPhase1AndEndTurn_RejectedInAFinishedDuel(t *testing.T) {
	d, _ := finishedDuelWithP1Monster(t)

	if err := d.EndPhase1(); !errors.Is(err, ErrDuelEnded) {
		t.Errorf("EndPhase1() = %v, want ErrDuelEnded", err)
	}
	if err := d.EndTurn(); !errors.Is(err, ErrDuelEnded) {
		t.Errorf("EndTurn() = %v, want ErrDuelEnded", err)
	}
}
func TestEndTurn_ResetsTheTurnStateOfBothParticipants(t *testing.T) {
	// Arrange — turn 1: p1's Monster already changed its position.
	d, p1Monster := duelWithMonsterAt(t, PositionAttack, FaceUp)
	if err := d.ChangePosition(p1Monster, PositionDefense, FaceUp); err != nil {
		t.Fatalf("p1 ChangePosition() = %v", err)
	}
	// Turn 2: p2 places its own Monster and changes its position too.
	endFullTurn(t, d)
	p2, _ := d.Participant("p2")
	p2Monster := p2.Hand().Cards()[0].ID()
	if err := d.Place([]InstanceID{p2Monster}, PositionAttack, FaceUp); err != nil {
		t.Fatalf("p2 Place() = %v", err)
	}
	if err := d.ChangePosition(p2Monster, PositionDefense, FaceUp); err != nil {
		t.Fatalf("p2 ChangePosition() = %v", err)
	}

	// Act — end turn 2.
	endFullTurn(t, d) // -> p1, turn 3

	// Assert — p1's Monster may change again: its flag was reset.
	if err := d.ChangePosition(p1Monster, PositionAttack, FaceUp); err != nil {
		t.Fatalf("p1 ChangePosition() on turn 3 = %v, want success", err)
	}
	endFullTurn(t, d) // -> p2, turn 4
	if err := d.ChangePosition(p2Monster, PositionAttack, FaceUp); err != nil {
		t.Fatalf("p2 ChangePosition() on turn 4 = %v, want success", err)
	}
}
