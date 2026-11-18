package story

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// Progress is the state of one player inside a campaign: which duels the
// player has already won.
//
// Progression is strictly linear (game-rules §56): the player may only play
// the first not-yet-won duel. Winning it moves the player to the next duel;
// losing it changes nothing, so the same duel can be played again (retry,
// docs/mvp-story-mode.md §1.1) without losing the previous victories.
//
// The state is therefore one fact: the number of consecutive duels won from
// the start of the campaign. Everything else (current level, current duel,
// won duels, completion) is derived from it, so no duplicated flags exist
// (domain rule 4).
//
// Progress is an immutable value: the result transitions return a new
// Progress and leave the receiver unchanged (domain rule 6: the result of a
// duel is calculated before the state applies it).
type Progress struct {
	// current is the number of duels already won, counting from the first
	// duel of the campaign. It is also the 0-based index of the duel the
	// player must play next.
	current int
}

// NewProgress returns the progress of a player that has not won any duel:
// the current duel is the first one of the campaign.
func NewProgress() Progress { return Progress{} }

// Count returns the number of duels the player has already won.
func (p Progress) Count() int { return p.current }

// IsComplete reports whether the player has won every duel of the campaign.
// An empty campaign is never complete.
func (p Progress) IsComplete(c Campaign) bool {
	if c.DuelCount() == 0 {
		return false
	}
	return p.current >= c.DuelCount()
}

// Current returns the level and the duel the player must play next: the
// first not-yet-won duel of the campaign (game-rules §56). ok is false when
// the campaign is empty, when the campaign is complete, or when the progress
// does not fit the campaign.
func (p Progress) Current(c Campaign) (Level, Duel, bool) {
	if p.IsComplete(c) {
		return Level{}, Duel{}, false
	}
	level, ok := c.Level(p.current/DuelsPerLevel + 1)
	if !ok {
		return Level{}, Duel{}, false
	}
	duel, ok := level.Duel(p.current%DuelsPerLevel + 1)
	if !ok {
		return Level{}, Duel{}, false
	}
	return level, duel, true
}

// HasWon reports whether the duel has already been won by the player. A duel
// that does not exist in the campaign has never been won.
func (p Progress) HasWon(c Campaign, level, duel int) bool {
	d, ok := c.Duel(level, duel)
	if !ok {
		return false
	}
	position := (level-1)*DuelsPerLevel + (d.Number() - 1)
	return p.current > position
}

// RecordVictory applies the victory of a duel to the progress and returns
// the new progress (game-rules §56: the victory of a duel advances the
// player to the next duel; the victory of the four duels of a level is
// required to reach the next level).
//
// The duel must be the current one: a victory of a duel that does not exist
// in the campaign is reported with ErrDuelNotFound, a victory of a duel that
// is not the current one (an already won duel or a future duel) with
// ErrNotCurrentDuel, and a victory after the campaign is complete with
// ErrCampaignComplete. In every error case the receiver is returned
// unchanged, so a rejected result never changes the state, and a won duel
// can never be played or rewarded again.
//
// The victory grants the three predefined cards of the duel (game-rules
// §58): they are exposed by LastRewards on the returned progress. The
// progress does not store cards: granting them to the player belongs to the
// Collection, which is outside the Story Mode MVP.
func (p Progress) RecordVictory(c Campaign, level, duel int) (Progress, error) {
	if p.IsComplete(c) {
		return p, fmt.Errorf("%w: all %d duels are already won", ErrCampaignComplete, c.DuelCount())
	}
	if _, ok := c.Duel(level, duel); !ok {
		return p, fmt.Errorf("%w: level %d, duel %d", ErrDuelNotFound, level, duel)
	}
	currentLevel, currentDuel, ok := p.Current(c)
	if !ok {
		return p, fmt.Errorf("%w: level %d, duel %d", ErrDuelNotFound, level, duel)
	}
	if level != currentLevel.Number() || duel != currentDuel.Number() {
		return p, fmt.Errorf("%w: the current duel is level %d, duel %d; got level %d, duel %d",
			ErrNotCurrentDuel, currentLevel.Number(), currentDuel.Number(), level, duel)
	}
	return Progress{current: p.current + 1}, nil
}

// RecordLoss applies the loss of a duel to the progress and returns the
// progress (game-rules §56: losing a duel does not advance the progression).
//
// The state is returned unchanged: the current duel remains the current duel
// and the player may play it again (retry, docs/mvp-story-mode.md §1.1)
// without losing the previous victories. The duel must be the current one,
// with the same error semantics of RecordVictory.
func (p Progress) RecordLoss(c Campaign, level, duel int) (Progress, error) {
	if p.IsComplete(c) {
		return p, fmt.Errorf("%w: all %d duels are already won", ErrCampaignComplete, c.DuelCount())
	}
	if _, ok := c.Duel(level, duel); !ok {
		return p, fmt.Errorf("%w: level %d, duel %d", ErrDuelNotFound, level, duel)
	}
	currentLevel, currentDuel, ok := p.Current(c)
	if !ok {
		return p, fmt.Errorf("%w: level %d, duel %d", ErrDuelNotFound, level, duel)
	}
	if level != currentLevel.Number() || duel != currentDuel.Number() {
		return p, fmt.Errorf("%w: the current duel is level %d, duel %d; got level %d, duel %d",
			ErrNotCurrentDuel, currentLevel.Number(), currentDuel.Number(), level, duel)
	}
	return p, nil
}

// LastRewards returns the three predefined cards granted by the duel that
// the player has most recently won (game-rules §58). The returned cards are
// the input for the player's Collection, which is outside the Story Mode
// MVP: the domain defines and exposes the reward, the application grants it.
//
// ok is false when the player has not won any duel yet or when the campaign
// does not contain the duel of the progress.
func (p Progress) LastRewards(c Campaign) ([]*domain.Card, bool) {
	if p.current == 0 {
		return nil, false
	}
	d, ok := c.Duel((p.current-1)/DuelsPerLevel+1, (p.current-1)%DuelsPerLevel+1)
	if !ok {
		return nil, false
	}
	return d.Rewards(), true
}
