package domain

// FirstTurnNumber is the number of the first Turn of a Duel.
const FirstTurnNumber = 1

// Turn identifies one turn of a Duel: its number, counted from
// FirstTurnNumber, and the Participant who plays it (game-rules §11,
// architecture §12).
//
// It is a small immutable value: the Duel creates it and nothing changes it
// afterwards. Phases and the advance to the next turn belong to the turn
// structure (task S4) and are not part of it yet. The zero value is not a
// turn of any Duel.
type Turn struct {
	number      int
	participant ParticipantID
}

// Number returns the turn number, starting at FirstTurnNumber.
func (t Turn) Number() int { return t.number }

// Participant returns the id of the Participant who plays the turn.
func (t Turn) Participant() ParticipantID { return t.participant }

// IsFirst reports whether this is the first turn of the Duel. In the Story
// Mode MVP the Participant who plays it does not attack (game-rules §6); that
// restriction is applied by the combat action, not by Turn.
func (t Turn) IsFirst() bool { return t.number == FirstTurnNumber }
