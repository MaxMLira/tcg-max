package domain

import "testing"

func TestTurn_Accessors(t *testing.T) {
	tests := []struct {
		name      string
		turn      Turn
		wantNum   int
		wantID    ParticipantID
		wantFirst bool
	}{
		{name: "first turn", turn: Turn{number: FirstTurnNumber, participant: "p1"}, wantNum: 1, wantID: "p1", wantFirst: true},
		{name: "second turn", turn: Turn{number: 2, participant: "p2"}, wantNum: 2, wantID: "p2", wantFirst: false},
		{name: "later turn", turn: Turn{number: 10, participant: "p1"}, wantNum: 10, wantID: "p1", wantFirst: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			num, id, first := tt.turn.Number(), tt.turn.Participant(), tt.turn.IsFirst()

			// Assert
			if num != tt.wantNum {
				t.Errorf("Number() = %d, want %d", num, tt.wantNum)
			}
			if id != tt.wantID {
				t.Errorf("Participant() = %q, want %q", id, tt.wantID)
			}
			if first != tt.wantFirst {
				t.Errorf("IsFirst() = %v, want %v", first, tt.wantFirst)
			}
		})
	}
}

func TestFirstTurnNumber_IsOne(t *testing.T) {
	// Assert
	if FirstTurnNumber != 1 {
		t.Errorf("FirstTurnNumber = %d, want 1", FirstTurnNumber)
	}
	if !(Turn{number: FirstTurnNumber}).IsFirst() {
		t.Errorf("a turn numbered FirstTurnNumber must be the first turn")
	}
	if (Turn{}).IsFirst() {
		t.Errorf("the zero Turn must not be the first turn")
	}
}
