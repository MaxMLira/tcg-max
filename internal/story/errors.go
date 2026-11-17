package story

import "errors"

// Story errors describe invalid Story Mode operations. Callers distinguish
// them with errors.Is, so every error returned by the package wraps one of
// these sentinels together with the detail that explains the failure.
//
// Story errors never carry transport or infrastructure semantics.
var (
	// ErrInvalidCharacter reports a Character that violates its own
	// invariants: a blank identity, a blank name or an invalid deck.
	ErrInvalidCharacter = errors.New("invalid story character")

	// ErrInvalidDuel reports a duel definition that violates its own
	// invariants: a duel number outside 1 to DuelsPerLevel, a missing
	// character or a missing reward card.
	ErrInvalidDuel = errors.New("invalid story duel")

	// ErrInvalidLevel reports a level that violates its own invariants: a
	// level number the domain cannot support, or a duel whose slot does not
	// match its number. Errors built from the domain rule wrap
	// domain.ErrInvalidLevel as well, so both sentinels match errors.Is.
	ErrInvalidLevel = errors.New("invalid story level")

	// ErrInvalidCampaign reports a campaign that violates its own
	// invariants: an empty campaign, level numbers that are not sequential
	// from 1, or more levels than the domain supports.
	ErrInvalidCampaign = errors.New("invalid story campaign")

	// ErrDuelNotFound reports an operation on a duel that does not exist in
	// the campaign.
	ErrDuelNotFound = errors.New("duel not found in campaign")

	// ErrNotCurrentDuel reports recording the result of a duel that is not
	// the current one: only the first not-yet-won duel can be played
	// (game-rules §56, docs/mvp-story-mode.md §1.1).
	ErrNotCurrentDuel = errors.New("duel is not the current duel")

	// ErrCampaignComplete reports recording the result of a duel when every
	// duel of the campaign has already been won.
	ErrCampaignComplete = errors.New("campaign is complete")
)
