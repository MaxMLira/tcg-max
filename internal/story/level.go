package story

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// Level is one stage of the Story Mode progression: a fixed sequence of
// DuelsPerLevel duels the player must win to reach the next level
// (game-rules §55 and §56, glossary §12 "Story Level").
//
// The level number counts from 1 in the order of the campaign. The initial
// Life Points of the duels of the level are defined by the domain for the
// level number (game-rules §3.1), so the level never stores its own life
// points value: InitialLifePoints delegates to the domain.
//
// The duels are stored in a fixed-size array: a level with a number of duels
// different from DuelsPerLevel is not representable.
type Level struct {
	number int
	duels  [DuelsPerLevel]Duel
}

// NewLevel validates and builds a level.
//
// The number must be a level the domain supports (domain.InitialLifePoints,
// 1 to domain.MaxSupportedLevel for now), because the level would otherwise
// promise duels whose initial Life Points the domain cannot define. The
// resulting error wraps both ErrInvalidLevel and domain.ErrInvalidLevel.
//
// Every duel must occupy the slot matching its number: the duel of slot 1
// has number 1, and so on. Progression is strictly linear
// (game-rules §56), so the slot of a duel is part of its identity. A zero
// value of Duel in any slot is rejected.
func NewLevel(number int, duels [DuelsPerLevel]Duel) (Level, error) {
	if _, err := domain.InitialLifePoints(number); err != nil {
		return Level{}, fmt.Errorf("%w: %w", ErrInvalidLevel, err)
	}
	for i, duel := range duels {
		expected := i + 1
		if duel.number != expected {
			return Level{}, fmt.Errorf("%w: duel at slot %d has number %d, expected %d", ErrInvalidLevel, expected, duel.number, expected)
		}
		if duel.character.id == "" {
			return Level{}, fmt.Errorf("%w: duel %d has no character", ErrInvalidLevel, expected)
		}
	}
	return Level{number: number, duels: duels}, nil
}

// Number returns the level number, counting from 1.
func (l Level) Number() int { return l.number }

// Duels returns the four duels of the level, in the order they are played.
func (l Level) Duels() [DuelsPerLevel]Duel { return l.duels }

// Duel returns the duel with the given number (1 to DuelsPerLevel). ok is
// false for a number outside the range.
func (l Level) Duel(number int) (Duel, bool) {
	if number < 1 || number > DuelsPerLevel {
		return Duel{}, false
	}
	return l.duels[number-1], true
}

// InitialLifePoints returns the Life Points a participant starts the duels of
// the level with (game-rules §3.1). The rule lives in the domain; this
// method only applies it to the level number. For a level built with
// NewLevel the call always succeeds; the error covers the zero value of
// Level.
func (l Level) InitialLifePoints() (int, error) {
	return domain.InitialLifePoints(l.number)
}
