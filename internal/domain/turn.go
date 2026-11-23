package domain

// FirstTurnNumber is the number of the first Turn of a Duel.
const FirstTurnNumber = 1

// Turn identifies one turn of a Duel: its number, counted from
// FirstTurnNumber, the Participant who plays it and its current Phase
// (game-rules §11, glossary "Turn", architecture §12).
//
// The Duel is the single writer of the Turn: it creates it and updates it
// when the phase or the turn changes (architecture §12). The zero value is
// not a turn of any Duel.
type Turn struct {
	number      int
	participant ParticipantID
	phase       Phase
}

// Number returns the turn number, starting at FirstTurnNumber.
func (t Turn) Number() int { return t.number }

// Participant returns the id of the Participant who plays the turn.
func (t Turn) Participant() ParticipantID { return t.participant }

// Phase returns the current Phase of the turn.
func (t Turn) Phase() Phase { return t.phase }

// IsFirst reports whether this is the first turn of the Duel. In the Story
// Mode MVP the Participant who plays it does not attack (game-rules §6); that
// restriction is applied by the combat action, not by Turn.
func (t Turn) IsFirst() bool { return t.number == FirstTurnNumber }
