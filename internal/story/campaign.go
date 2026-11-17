package story

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// Campaign is the complete Story Mode progression: the ordered levels the
// player faces from the first duel to the last (game-rules §55,
// architecture §9 "story configuration").
//
// The campaign is an immutable definition shared by every player. The state
// of a player inside the campaign is a separate value: Progress.
type Campaign struct {
	levels []Level
}

// NewCampaign validates and builds a campaign from its levels, in the order
// they are played.
//
// Invariants:
//   - the campaign has at least one level;
//   - the level numbers are sequential from 1, in the given order: the first
//     level is level 1, the second is level 2, and so on. Progression is
//     strictly linear (game-rules §56), so the order of the levels is part
//     of the definition;
//   - the campaign does not contain more levels than the domain supports
//     (domain.MaxSupportedLevel), because every level must be able to define
//     the initial Life Points of its duels (game-rules §3.1).
func NewCampaign(levels ...Level) (Campaign, error) {
	if len(levels) == 0 {
		return Campaign{}, fmt.Errorf("%w: at least one level is required", ErrInvalidCampaign)
	}
	if len(levels) > domain.MaxSupportedLevel {
		return Campaign{}, fmt.Errorf("%w: %d levels, the domain supports at most %d", ErrInvalidCampaign, len(levels), domain.MaxSupportedLevel)
	}
	for i, level := range levels {
		if level.number != i+1 {
			return Campaign{}, fmt.Errorf("%w: level at position %d has number %d, expected %d", ErrInvalidCampaign, i+1, level.number, i+1)
		}
	}
	return Campaign{levels: append([]Level(nil), levels...)}, nil
}

// LevelCount returns the number of levels of the campaign.
func (c Campaign) LevelCount() int { return len(c.levels) }

// DuelCount returns the total number of duels of the campaign: LevelCount
// times DuelsPerLevel.
func (c Campaign) DuelCount() int { return len(c.levels) * DuelsPerLevel }

// Levels returns a copy of the levels, in the order they are played.
// Changing the returned slice does not affect the campaign.
func (c Campaign) Levels() []Level {
	return append([]Level(nil), c.levels...)
}

// Level returns the level with the given number. ok is false for a number
// outside the campaign.
func (c Campaign) Level(number int) (Level, bool) {
	if number < 1 || number > len(c.levels) {
		return Level{}, false
	}
	return c.levels[number-1], true
}

// Duel returns the duel with the given level and duel numbers. ok is false
// when either number is outside the campaign.
func (c Campaign) Duel(level, duel int) (Duel, bool) {
	l, ok := c.Level(level)
	if !ok {
		return Duel{}, false
	}
	return l.Duel(duel)
}
