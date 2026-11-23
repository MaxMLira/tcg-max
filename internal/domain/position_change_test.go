package domain

import (
	"errors"
	"reflect"
	"testing"
)

// duelWithMonsterAt builds a Duel where p1 has placed a single weight-2
// Monster on the Field in the given position and face, and it is still
// Phase 1 of p1's turn. It returns the Duel and the Monster's id.
func duelWithMonsterAt(t *testing.T, position Position, face Face) (*Duel, InstanceID) {
	t.Helper()

	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	id := InstanceID("p1#1")
	if err := d.Place([]InstanceID{id}, position, face); err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	return d, id
}

func TestChangePosition_FaceUpTogglesBetweenAttackAndDefense(t *testing.T) {
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)

	// Turn 1 — Attack -> Defense.
	if err := d.ChangePosition(id, PositionDefense, FaceUp); err != nil {
		t.Fatalf("ChangePosition() = %v", err)
	}
	c, _ := d.Participants()[0].Field().Get(id)
	if c.Position() != PositionDefense || c.Face() != FaceUp {
		t.Errorf("monster = (%v, %v), want (Defense, Face-up)", c.Position(), c.Face())
	}

	// Turn 3 — the restriction resets, so the monster toggles back.
	endFullTurn(t, d)
	endFullTurn(t, d)
	if err := d.ChangePosition(id, PositionAttack, FaceUp); err != nil {
		t.Fatalf("ChangePosition() on turn 3 = %v", err)
	}
	c, _ = d.Participants()[0].Field().Get(id)
	if c.Position() != PositionAttack || c.Face() != FaceUp {
		t.Errorf("monster = (%v, %v), want (Attack, Face-up)", c.Position(), c.Face())
	}
}

func TestChangePosition_RevealsAFaceDownMonster(t *testing.T) {
	d, id := duelWithMonsterAt(t, PositionDefense, FaceDown)

	if err := d.ChangePosition(id, PositionAttack, FaceUp); err != nil {
		t.Fatalf("ChangePosition() = %v", err)
	}
	c, _ := d.Participants()[0].Field().Get(id)
	if c.Position() != PositionAttack || c.Face() != FaceUp {
		t.Errorf("monster = (%v, %v), want (Attack, Face-up)", c.Position(), c.Face())
	}
}

func TestChangePosition_RevealsFaceDownKeepingDefense(t *testing.T) {
	// Arrange — same position, different face, so it is allowed.
	d, id := duelWithMonsterAt(t, PositionDefense, FaceDown)

	// Act
	err := d.ChangePosition(id, PositionDefense, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("ChangePosition() = %v", err)
	}
	c, _ := d.Participants()[0].Field().Get(id)
	if c.Position() != PositionDefense || c.Face() != FaceUp {
		t.Errorf("monster = (%v, %v), want (Defense, Face-up)", c.Position(), c.Face())
	}
}
func TestChangePosition_RejectsFaceUpToFaceDown(t *testing.T) {
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	before := duelStateOf(d)

	err := d.ChangePosition(id, PositionDefense, FaceDown)

	if !errors.Is(err, ErrInvalidPositionChange) {
		t.Fatalf("ChangePosition() = %v, want ErrInvalidPositionChange", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected change")
	}
}

func TestChangePosition_RejectsAnIncoherentTarget(t *testing.T) {
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	before := duelStateOf(d)

	err := d.ChangePosition(id, PositionAttack, FaceDown)

	if !errors.Is(err, ErrInvalidPositionChange) {
		t.Fatalf("ChangePosition() = %v, want ErrInvalidPositionChange", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected change")
	}
}

func TestChangePosition_RejectsTheSameState(t *testing.T) {
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	before := duelStateOf(d)

	err := d.ChangePosition(id, PositionAttack, FaceUp)

	if !errors.Is(err, ErrInvalidPositionChange) {
		t.Fatalf("ChangePosition() = %v, want ErrInvalidPositionChange", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected change")
	}
}

func TestChangePosition_AtMostOncePerTurn(t *testing.T) {
	// Arrange – one change already used on this turn.
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	if err := d.ChangePosition(id, PositionDefense, FaceUp); err != nil {
		t.Fatalf("first ChangePosition() = %v", err)
	}
	before := duelStateOf(d)

	// Act
	err := d.ChangePosition(id, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrInvalidPositionChange) {
		t.Fatalf("second ChangePosition() = %v, want ErrInvalidPositionChange", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by the second change")
	}

	// The restriction resets on the next turn.
	endFullTurn(t, d)
	endFullTurn(t, d)
	if err := d.ChangePosition(id, PositionAttack, FaceUp); err != nil {
		t.Errorf("ChangePosition() on the next turn = %v, want success", err)
	}
}

func TestChangePosition_RejectsACardNotOnTheField(t *testing.T) {
	d, _ := duelWithMonsterAt(t, PositionAttack, FaceUp)
	before := duelStateOf(d)

	err := d.ChangePosition("p1#5", PositionDefense, FaceUp)

	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("ChangePosition() = %v, want ErrCardNotFound", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected change")
	}
}

func TestChangePosition_DoesNotConsumeThePlacement(t *testing.T) {
	// Arrange — the monster was placed on turn 1, which consumed that turn's
	// placement. On p1's next turn the placement is fresh.
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	endFullTurn(t, d)
	endFullTurn(t, d)

	// Act — change the position, then place another card on the same turn.
	if err := d.ChangePosition(id, PositionDefense, FaceUp); err != nil {
		t.Fatalf("ChangePosition() = %v", err)
	}
	p1, _ := d.Participant("p1")
	err := d.Place([]InstanceID{p1.Hand().Cards()[0].ID()}, PositionAttack, FaceUp)

	// Assert — the position change did not consume the placement.
	if err != nil {
		t.Fatalf("Place() after a position change = %v, want success", err)
	}
	if c, _ := p1.Field().Get(id); c.Position() != PositionDefense {
		t.Errorf("the changed monster is in %v, want Defense", c.Position())
	}
	if got := len(p1.Field().Monsters()); got != 2 {
		t.Errorf("p1's field holds %d monsters, want 2", got)
	}
}

func TestChangePosition_RejectedInPhase2(t *testing.T) {
	// Arrange
	d, id := duelWithMonsterAt(t, PositionAttack, FaceUp)
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() = %v", err)
	}
	before := duelStateOf(d)

	// Act
	err := d.ChangePosition(id, PositionDefense, FaceUp)

	// Assert
	if !errors.Is(err, ErrWrongPhase) {
		t.Fatalf("ChangePosition() = %v, want ErrWrongPhase", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected change")
	}
}

func TestChangePosition_RejectedInAFinishedDuel(t *testing.T) {
	d, id := finishedDuelWithP1Monster(t)

	err := d.ChangePosition(id, PositionAttack, FaceUp)

	if !errors.Is(err, ErrDuelEnded) {
		t.Fatalf("ChangePosition() = %v, want ErrDuelEnded", err)
	}
}
