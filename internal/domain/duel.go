package domain

import (
	"errors"
	"fmt"
)

// ParticipantSetup is the input for one side of a Duel: the Participant id,
// its deck composition and its starting Life Points. Use InitialLifePoints to
// obtain the Life Points of a Story Level.
type ParticipantSetup struct {
	ID         ParticipantID
	Deck       DeckDefinition
	LifePoints int
}

// Duel is the state of a match between two Participants (architecture §7,
// ADR-011). It owns both Participants, the current Turn, the placement state
// of the turn and, once it has ended, the DuelResult.
//
// The Duel engine does not know whether a Participant is a human or an AI
// (domain rule 10). Whether the Duel is finished is derived from the result
// (domain rule 4): there is no separate status flag.
//
// The actions of the Duel act for the Participant who plays the current Turn.
// The Hand is refilled automatically; there is no manual draw (game-rules §7,
// Story Mode MVP). Because of that the replenishment is the unexported
// replenishHand, to be called by the Duel after a resolved action. The turn
// runs Phase 1, Phase 2 and End Turn (game-rules §11, Story Mode MVP).
//
// A Duel is created by NewDuel; the zero value is not usable. It holds no lock
// and no Shuffler: it is a sequential state machine serialized by its caller.
type Duel struct {
	participants  [2]*Participant
	turn          Turn
	placementUsed bool
	result        *DuelResult
}

// NewDuel prepares a Duel between the Participants a and b: it builds both
// Participants with shuffled decks, deals the opening hands and chooses the
// first player at random (game-rules §5.1, §6 and mvp-story-mode A5).
//
// It is all-or-nothing: everything is built in local values and the Duel is
// returned only when every step succeeded, so no half-built Duel exists. The
// steps run in this order:
//
//  1. both ids are validated (ErrInvalidParticipant);
//  2. equal ids, compared exactly (ids are stored as given), are rejected with
//     ErrInvalidDuel, because the InstanceIDs "<participant id>#<n>" would
//     collide. This happens before the Shuffler is used;
//  3. the Participants are built, a first and then b, each shuffling its Deck
//     once; the deck, Life Points and nil-Shuffler errors are propagated;
//  4. each Participant draws its opening hand, a first and then b: exactly 4
//     cards, with no extra card or advantage for the first player
//     (Story Mode MVP, game-rules §5.1 and §6);
//  5. the first player is chosen with one more Shuffle over the two ids.
//
// The Duel uses one Shuffler, called with the sizes 30, 30 and 2 in that
// order. If NewDuel fails the Shuffler may already have advanced, so the
// caller must discard it. A seeded generator (math/rand/v2, rand.New with a
// PCG source) is created by the composition root or by the tests, never by the
// domain.
//
// The first turn is Turn number FirstTurnNumber of the chosen player, in
// Phase 1. The rule that the first player does not attack in that turn is
// applied by the combat action of the combat task (game-rules §6,
// Turn.IsFirst).
func NewDuel(a, b ParticipantSetup, s Shuffler) (*Duel, error) {
	if err := validateParticipantID(a.ID); err != nil {
		return nil, fmt.Errorf("first participant: %w", err)
	}
	if err := validateParticipantID(b.ID); err != nil {
		return nil, fmt.Errorf("second participant: %w", err)
	}
	if a.ID == b.ID {
		return nil, fmt.Errorf("%w: both participants have the id %q", ErrInvalidDuel, a.ID)
	}

	pa, err := NewParticipant(a.ID, a.Deck, a.LifePoints, s)
	if err != nil {
		return nil, fmt.Errorf("first participant %q: %w", a.ID, err)
	}
	pb, err := NewParticipant(b.ID, b.Deck, b.LifePoints, s)
	if err != nil {
		return nil, fmt.Errorf("second participant %q: %w", b.ID, err)
	}

	// A validated deck has DeckSize cards, so this cannot fail; if it does it
	// is a programming error and no Duel is returned.
	if err := pa.replenishHand(); err != nil {
		return nil, fmt.Errorf("opening hand of %q: %w", a.ID, err)
	}
	if err := pb.replenishHand(); err != nil {
		return nil, fmt.Errorf("opening hand of %q: %w", b.ID, err)
	}

	candidates := [2]ParticipantID{a.ID, b.ID}
	s.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return &Duel{
		participants: [2]*Participant{pa, pb},
		turn:         Turn{number: FirstTurnNumber, participant: candidates[0], phase: Phase1},
	}, nil
}

// Participants returns the two Participants, in the order of the setups given
// to NewDuel. The array is a copy; the Participants are shared and their
// zones can only be changed through their own transitions.
func (d *Duel) Participants() [2]*Participant { return d.participants }

// Participant returns the Participant with the given id and whether it takes
// part in the Duel.
func (d *Duel) Participant(id ParticipantID) (*Participant, bool) {
	for _, p := range d.participants {
		if p.ID() == id {
			return p, true
		}
	}
	return nil, false
}

// Turn returns the current Turn.
func (d *Duel) Turn() Turn { return d.turn }

// Result returns the result of the Duel and true when it has ended; otherwise
// it returns the zero DuelResult and false.
func (d *Duel) Result() (DuelResult, bool) {
	if d.result == nil {
		return DuelResult{}, false
	}
	return *d.result, true
}

// Finished reports whether the Duel has ended.
func (d *Duel) Finished() bool { return d.result != nil }

// EndPhase1 ends Phase 1 of the current Turn and starts Phase 2
// (game-rules §11, Story Mode MVP: the turn runs Phase 1, Phase 2 and End
// Turn). It is only allowed in Phase 1.
func (d *Duel) EndPhase1() error {
	if d.result != nil {
		return fmt.Errorf("%w: cannot act in a finished duel", ErrDuelEnded)
	}
	if d.turn.phase != Phase1 {
		return fmt.Errorf("%w: the turn is already in %s", ErrWrongPhase, d.turn.phase)
	}
	d.turn.phase = Phase2
	return nil
}

// EndTurn ends the current Turn and starts the next one, which the other
// Participant plays (game-rules §11, glossary "Turn"): the turn number
// increments, the phase returns to Phase 1, the placement of the phase is
// fresh and the turn-specific state of the Monsters is reset
// (architecture §12). It is only allowed in Phase 2: the MVP turn runs Phase
// 1, Phase 2 and End Turn.
func (d *Duel) EndTurn() error {
	if d.result != nil {
		return fmt.Errorf("%w: cannot act in a finished duel", ErrDuelEnded)
	}
	if d.turn.phase != Phase2 {
		return fmt.Errorf("%w: ending the turn is a phase 2 action and the turn is in %s", ErrWrongPhase, d.turn.phase)
	}

	var next ParticipantID
	if d.participants[0].ID() == d.turn.participant {
		next = d.participants[1].ID()
	} else {
		next = d.participants[0].ID()
	}

	d.turn = Turn{number: d.turn.number + 1, participant: next, phase: Phase1}
	d.placementUsed = false
	for _, p := range d.participants {
		p.resetTurnState()
	}
	return nil
}

// replenishHand refills the Hand of the Participant with the given id up to
// HandLimit cards (game-rules §5.2 and §7, Story Mode MVP). It is the only way
// the Duel draws, and it is unexported on purpose: an exported draw would be
// the manual draw action that the MVP does not have. The actions of later
// tasks call it after an action is resolved.
//
// It returns ErrDuelEnded, wrapped, when the Duel has a result: a finished
// Duel is frozen. It returns ErrInvalidParticipant for an id that is not in the
// Duel. A Hand that is already full is left as it is.
//
// When the Deck cannot supply the missing cards the Participant loses
// immediately (game-rules §7 and §8): the other Participant wins with
// EndReasonDeckOut and nil is returned, because a loss is an outcome of the
// game, not a rejected action. The Participant's zones are unchanged in that
// case. Any other error is propagated.
func (d *Duel) replenishHand(id ParticipantID) error {
	if d.result != nil {
		return fmt.Errorf("%w: cannot replenish the hand of %q", ErrDuelEnded, id)
	}
	p, ok := d.Participant(id)
	if !ok {
		return fmt.Errorf("%w: %q does not take part in the duel", ErrInvalidParticipant, id)
	}

	err := p.replenishHand()
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrDeckEmpty) {
		return err
	}

	winner := d.participants[0].ID()
	if winner == id {
		winner = d.participants[1].ID()
	}
	d.result = &DuelResult{winner: winner, loser: id, reason: EndReasonDeckOut}
	return nil
}
