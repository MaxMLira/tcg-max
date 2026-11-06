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
)
