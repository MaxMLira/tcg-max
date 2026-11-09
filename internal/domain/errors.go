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
)
