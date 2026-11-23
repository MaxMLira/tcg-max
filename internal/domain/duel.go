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

// requiredTributes returns the number of tributes the Weight of a Monster
// requires for its summon (game-rules §16): 0 for 1 to 5, 1 for 6 to 8, 2 for
// 9 to 12 and 3 for 13 and more.
func requiredTributes(weight int) int {
	switch {
	case weight <= 5:
		return 0
	case weight <= 8:
		return 1
	case weight <= 12:
		return 2
	default:
		return 3
	}
}

// activeParticipant returns the Participant who plays the current Turn. The
// Turn always names one of the Duel's Participants, so it cannot fail.
func (d *Duel) activeParticipant() *Participant {
	if p := d.participants[0]; p.ID() == d.turn.participant {
		return p
	}
	return d.participants[1]
}

// Place is the selection, discard and placement action (mvp-story-mode §1.5,
// game-rules §13.3 and §14). The active Participant selects 1 to 4 distinct
// cards from its Hand; every selected card except the last goes to the
// Graveyard, in selection order, and the last card is placed on the Field in
// the given position and face.
//
// The action is atomic: a rejected action leaves the Hand, the Field, the
// Graveyard and the turn state unchanged. It is only allowed in Phase 1, at
// most once per Phase 1 (game-rules §14.1).
//
// In this scope the selected cards must be Monsters: the non-last cards are
// the tributes, which are Monsters (game-rules §17), and the last card is a
// Monster whose summon requires, by Weight, 0 to 3 tributes
// (game-rules §16). A face-down summon is only allowed without tributes
// (game-rules §19). Placing a Magic or a Trap is rejected with
// ErrUnsupportedCardType: Magic placement arrives with the Equipment task and
// Trap placement is rejected (mvp-story-mode A2).
//
// A Monster whose summon needs tributes receives its tributes, its position
// and its face in this same action (mvp-story-mode A4). When the action
// succeeds, the Hand is replenished automatically (game-rules §5.2 and §7,
// Story Mode MVP); a Deck that cannot complete the Hand loses the Duel at
// once (game-rules §8).
func (d *Duel) Place(selection []InstanceID, position Position, face Face) error {
	if d.result != nil {
		return fmt.Errorf("%w: cannot act in a finished duel", ErrDuelEnded)
	}
	if d.turn.phase != Phase1 {
		return fmt.Errorf("%w: placement is a phase 1 action and the turn is in %s", ErrWrongPhase, d.turn.phase)
	}
	if d.placementUsed {
		return fmt.Errorf("%w: only one placement is allowed in phase 1 (game-rules §14.1)", ErrPlacementUsed)
	}
	if !validPositionFace(position, face) {
		return fmt.Errorf("%w: %s %s is not a coherent position and face", ErrInvalidPlacement, face, position)
	}

	p := d.activeParticipant()

	if len(selection) < 1 || len(selection) > 4 {
		return fmt.Errorf("%w: the selection holds %d cards and must hold 1 to 4", ErrInvalidSelection, len(selection))
	}

	seen := make(map[InstanceID]struct{}, len(selection))
	cards := make([]*CardInstance, len(selection))
	for i, id := range selection {
		if _, dup := seen[id]; dup {
			return fmt.Errorf("%w: %q is selected more than once", ErrInvalidSelection, id)
		}
		seen[id] = struct{}{}
		card, ok := p.Hand().Get(id)
		if !ok {
			return fmt.Errorf("%w: %q is not in the hand of %q", ErrCardNotFound, id, p.ID())
		}
		if card.Card().Type() != CardTypeMonster {
			if i == len(selection)-1 {
				return fmt.Errorf("%w: a %s cannot be placed", ErrUnsupportedCardType, card.Card().Type())
			}
			return fmt.Errorf("%w: only monsters can be tributed", ErrInvalidSelection)
		}
		cards[i] = card
	}

	placed := cards[len(selection)-1]
	tributes := len(selection) - 1
	need := requiredTributes(placed.Card().Weight())
	if tributes < need {
		return fmt.Errorf("%w: a monster of weight %d requires %d tributes and the selection carries %d", ErrInvalidPlacement, placed.Card().Weight(), need, tributes)
	}
	if face == FaceDown && tributes > 0 {
		return fmt.Errorf("%w: a monster summoned with tributes must enter face-up (game-rules §19)", ErrInvalidPlacement)
	}

	if err := p.applyPlacement(selection, position, face); err != nil {
		return err
	}

	d.placementUsed = true
	return d.replenishHand(p.ID())
}

// ChangePosition changes the position of a Monster of the active Participant
// on its Field (game-rules §21.1). It is a phase 1 action and it does not
// count as a placement (game-rules §14.2).
//
// A face-up Monster switches between Attack and Defense; a face-down Monster
// is revealed by its controller and enters face-up, in Attack or Defense
// (game-rules §20). The target must be a coherent pair (a face-down Monster
// is always in Defense, game-rules §18.2) and different from the current
// position and face, and a Monster changes its position at most once per turn
// (mvp-story-mode S4). The restriction that a Monster that attacked in the
// turn cannot change its position (game-rules §21.1) is enforced by the
// combat state of the combat task.
func (d *Duel) ChangePosition(id InstanceID, position Position, face Face) error {
	if d.result != nil {
		return fmt.Errorf("%w: cannot act in a finished duel", ErrDuelEnded)
	}
	if d.turn.phase != Phase1 {
		return fmt.Errorf("%w: changing position is a phase 1 action and the turn is in %s", ErrWrongPhase, d.turn.phase)
	}
	if !validPositionFace(position, face) {
		return fmt.Errorf("%w: %s %s is not a coherent position and face", ErrInvalidPositionChange, face, position)
	}

	p := d.activeParticipant()
	card, ok := p.Field().Get(id)
	if !ok {
		return fmt.Errorf("%w: %q is not on the field of %q", ErrCardNotFound, id, p.ID())
	}
	if card.Card().Type() != CardTypeMonster {
		return fmt.Errorf("%w: %q is a %s, not a monster", ErrInvalidPositionChange, id, card.Card().Type())
	}
	if card.positionChanged() {
		return fmt.Errorf("%w: %q already changed its position in this turn", ErrInvalidPositionChange, id)
	}
	if card.Position() == position && card.Face() == face {
		return fmt.Errorf("%w: %q is already in %s %s", ErrInvalidPositionChange, id, face, position)
	}
	if card.Face() == FaceUp && face == FaceDown {
		return fmt.Errorf("%w: %q is a face-up monster and cannot be turned face-down", ErrInvalidPositionChange, id)
	}

	return p.changePosition(id, position, face)
}

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
