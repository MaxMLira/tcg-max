package domain

import "testing"

func TestPhaseString(t *testing.T) {
	tests := []struct {
		phase Phase
		want  string
	}{
		{Phase1, "Phase 1"},
		{Phase2, "Phase 2"},
		{Phase(0), "Phase(0)"},
	}
	for _, tt := range tests {
		if got := tt.phase.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestNewDuel_StartsInPhase1(t *testing.T) {
	// Act
	d := newTestDuel(t, noopShuffler{})

	// Assert
	if got := d.Turn().Phase(); got != Phase1 {
		t.Errorf("Turn().Phase() = %v, want %v", got, Phase1)
	}
}
