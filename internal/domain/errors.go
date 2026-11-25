package domain

import "errors"

// Domain errors describe invalid game operations. Callers distinguish them with
// errors.Is, so every error returned by the domain wraps one of these sentinels
// together with the detail that explains the failure.
//
// Domain errors never carry transport or infrastructure semantics.
var (
	// ErrInvalidCard reports a card definition that violates its own invariants
	// (missing identity, negative combat values, unsupported card type, ...).
	ErrInvalidCard = errors.New("invalid card")

	// ErrInvalidSubtype reports a Monster subtype that violates its own
	// invariant: a subtype must not be blank.
	ErrInvalidSubtype = errors.New("invalid subtype")

	// ErrInvalidDeck reports a malformed deck composition: a nil card, an entry
	// without copies, or a composition that cannot be interpreted.
	ErrInvalidDeck = errors.New("invalid deck")

	// ErrDeckSize reports a deck whose total number of cards is not the size
	// required by the game rules.
	ErrDeckSize = errors.New("invalid deck size")

	// ErrTooManyCopies reports a deck entry that exceeds the maximum number of
	// copies allowed for the same card.
	ErrTooManyCopies = errors.New("too many copies of the same card")

	// ErrDuplicateCard reports the same card appearing in more than one deck
	// entry. Entries are never merged silently, so the composition is rejected.
	ErrDuplicateCard = errors.New("duplicate card in deck")

	// ErrInvalidInstance reports a card instance that cannot be built: a
	// missing card definition, a blank identity or a blank owner
	// (architecture §5.5, game-rules §64).
	ErrInvalidInstance = errors.New("invalid card instance")

	// ErrInvalidParticipant reports a participant, or a participant-owned
	// zone, that cannot be built: a blank identity, an identity with
	// surrounding whitespace or a non-positive life points value (a Duel ends
	// when Life Points reach 0 or less, mvp-story-mode §1.2).
	ErrInvalidParticipant = errors.New("invalid participant")

	// ErrInvalidLevel reports a Story Level outside the supported range. The
	// game rules define levels 1 to 10 (game-rules §3.1); only the levels up to
	// MaxSupportedLevel are supported for now.
	ErrInvalidLevel = errors.New("invalid level")

	// ErrHandFull reports adding a card to a Hand that already holds the
	// maximum number of cards (game-rules §5.2).
	ErrHandFull = errors.New("hand is full")

	// ErrCardNotFound reports an operation on a card instance that is not in
	// the zone the operation reads from.
	ErrCardNotFound = errors.New("card not found in zone")

	// ErrFieldFull reports placing a card in a Field area (monster or
	// Magic/Trap) that has no free slot (game-rules §9.1 and §9.2).
	ErrFieldFull = errors.New("field area is full")

	// ErrDeckEmpty reports drawing from a Deck that has no cards (game-rules
	// §8). The consequence of the failed draw (losing the Duel) belongs to the
	// Duel, not to the Deck.
	ErrDeckEmpty = errors.New("deck is empty")

	// ErrInvalidDuel reports a Duel that cannot be built: the two Participants
	// have the same id, so their InstanceIDs would collide (ADR-011, architecture §40).
	// Errors of the Participants themselves keep their own sentinels.
	ErrInvalidDuel = errors.New("invalid duel")

	// ErrDuelEnded reports an operation on a Duel that already has a result:
	// a finished Duel accepts no further state change. Game-rules §8 defines
	// only the loss by an exhausted Deck (and §2 that the Duel ends at once);
	// the freeze itself is the S3 design decision (mvp-story-mode §1.2, "duel
	// ends immediately").
	ErrDuelEnded = errors.New("duel has ended")

	// ErrWrongPhase reports an action performed in a Phase where it is not
	// permitted (game-rules §11, §14 and §23, Story Mode MVP): the placement
	// and the position actions and EndPhase1 are Phase 1 actions, and EndTurn
	// is a Phase 2 action.
	ErrWrongPhase = errors.New("action not allowed in the current phase")

	// ErrPlacementUsed reports a second placement in the same Phase 1: only
	// one placement is allowed per Phase 1 (game-rules §14.1).
	ErrPlacementUsed = errors.New("a card was already placed in this phase 1")

	// ErrInvalidSelection reports a malformed selection for the placement
	// action: a size outside 1 to 4, a card selected more than once, or a
	// non-Monster among the selected cards (mvp-story-mode §1.5 and A4,
	// game-rules §17).
	ErrInvalidSelection = errors.New("invalid card selection")

	// ErrInvalidPlacement reports a card that cannot be placed: a Monster
	// without enough tributes for its Weight, or a position and face that are
	// not a legal summon (game-rules §16, §18.2 and §19).
	ErrInvalidPlacement = errors.New("card cannot be placed")

	// ErrUnsupportedCardType reports a placement of a card type without a
	// placement rule in the current scope: Magic placement arrives with the
	// Equipment task, and Trap placement is rejected (mvp-story-mode A2).
	ErrUnsupportedCardType = errors.New("card type has no placement rule")

	// ErrInvalidPositionChange reports a change of position of a Monster that
	// is not permitted: an incoherent target, the current position and face
	// again, a face-up Monster turned face-down, a second change in the same
	// turn, or a card that is not a Monster on the Field (game-rules §18.2,
	// §20 and §21.1, mvp-story-mode S4).
	ErrInvalidPositionChange = errors.New("invalid position change")

	// ErrInvalidEquipmentActivation reports an activation of a Magic
	// Equipment that is not permitted (game-rules §35.2,
	// docs/mvp-story-mode.md §1.4): the Equipment is not on the Field, the
	// card is not an Equipment, the Equipment is already linked to a Monster
	// (the MVP defines no re-link), the target is not a face-up Monster on the
	// Field of the active Participant, or the target's subtype does not match
	// the accepted subtype of the Equipment.
	ErrInvalidEquipmentActivation = errors.New("invalid equipment activation")
)
