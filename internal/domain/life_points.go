package domain

import "fmt"

// MaxSupportedLevel is the highest Story Level the domain currently supports.
//
// The game rules define levels 1 to 10 (game-rules §3.1), but the first part of
// the Story Mode MVP implements levels 1 to 5 only. Levels 6 to 10 (which also
// start with 8000 Life Points, game-rules §3.1) will be added later and are
// intentionally rejected by InitialLifePoints until then.
const MaxSupportedLevel = 5

// InitialLifePoints returns the Life Points a Participant starts a Duel with
// for the given Story Level (game-rules §3.1): 4000 for levels 1 to 4 and 8000
// for level 5.
//
// Levels below 1 and above MaxSupportedLevel are rejected with ErrInvalidLevel.
// The level is a plain int on purpose: the Level type belongs to the Story Mode
// progression, not to this rule.
func InitialLifePoints(level int) (int, error) {
	switch {
	case level >= 1 && level <= 4:
		return 4000, nil
	case level == 5:
		return 8000, nil
	default:
		return 0, fmt.Errorf("%w: level %d, supported levels are 1 to %d", ErrInvalidLevel, level, MaxSupportedLevel)
	}
}
