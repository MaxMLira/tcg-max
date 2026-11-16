package domain

// EndReason is why a Duel ended. The zero value means no reason.
type EndReason int

const (
	// EndReasonDeckOut means the loser could not complete the Hand because the
	// Deck had too few cards (game-rules §7 and §8). The reason for losing
	// all Life Points is added with the combat task (S6).
	EndReasonDeckOut EndReason = iota + 1
)

// DuelResult is the outcome of a finished Duel: who won, who lost and why
// (mvp-story-mode §1.2). It is an immutable value created by the Duel.
type DuelResult struct {
	winner ParticipantID
	loser  ParticipantID
	reason EndReason
}

// Winner returns the id of the Participant who won the Duel.
func (r DuelResult) Winner() ParticipantID { return r.winner }

// Loser returns the id of the Participant who lost the Duel.
func (r DuelResult) Loser() ParticipantID { return r.loser }

// Reason returns why the Duel ended.
func (r DuelResult) Reason() EndReason { return r.reason }
