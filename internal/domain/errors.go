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
)
