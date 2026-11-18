package domain

import "testing"

func TestDuelResult_Accessors(t *testing.T) {
	// Arrange
	result := DuelResult{winner: "p2", loser: "p1", reason: EndReasonDeckOut}

	// Act / Assert
	if got := result.Winner(); got != "p2" {
		t.Errorf("Winner() = %q, want p2", got)
	}
	if got := result.Loser(); got != "p1" {
		t.Errorf("Loser() = %q, want p1", got)
	}
	if got := result.Reason(); got != EndReasonDeckOut {
		t.Errorf("Reason() = %v, want EndReasonDeckOut", got)
	}
}

func TestEndReason_ZeroValueIsNoReason(t *testing.T) {
	// Arrange
	var zero EndReason

	// Assert
	if zero == EndReasonDeckOut {
		t.Errorf("the zero EndReason must not be EndReasonDeckOut")
	}
	if got := (DuelResult{}).Reason(); got != 0 {
		t.Errorf("zero DuelResult Reason() = %v, want 0", got)
	}
}
