package story

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// DuelsPerLevel is the exact number of duels of every Story Level
// (game-rules §55: each level contains a fixed sequence of 4 duels).
const DuelsPerLevel = 4

// Duel is the definition of one duel of the campaign: its position in the
// level, the Character the player faces and the three predefined cards the
// victory grants (game-rules §57 and §58, docs/mvp-story-mode.md §1.1).
//
// A story Duel is the scheduled slot of the progression, not the runtime
// Duel executed by the Duel Engine (tasks S3 to S7): the runtime Duel is
// built from this definition (the opponent deck, the initial Life Points of
// the level, the player's deck) by the application.
//
// The duel number is the position of the duel inside the level, counting
// from 1. The four duels of a level are played in this order, and the
// difficulty grows from duel 1 to duel 4 (game-rules §57).
type Duel struct {
	number    int
	character Character
	rewards   [3]*domain.Card
}

// NewDuel validates and builds a story duel.
//
// The number must be between 1 and DuelsPerLevel. The character must have
// been built with NewCharacter. The rewards must be the three predefined
// cards of the victory (game-rules §58); a nil reward card is rejected.
// Duplicated reward cards are allowed: the rule states three cards, not
// three distinct cards.
func NewDuel(number int, character Character, rewards [3]*domain.Card) (Duel, error) {
	if number < 1 || number > DuelsPerLevel {
		return Duel{}, fmt.Errorf("%w: number %d, must be between 1 and %d", ErrInvalidDuel, number, DuelsPerLevel)
	}
	if character.id == "" {
		return Duel{}, fmt.Errorf("%w: a duel requires a character", ErrInvalidDuel)
	}
	for i, card := range rewards {
		if card == nil {
			return Duel{}, fmt.Errorf("%w: reward %d is missing", ErrInvalidDuel, i+1)
		}
	}
	return Duel{number: number, character: character, rewards: rewards}, nil
}

// Number returns the position of the duel inside the level, counting from 1.
func (d Duel) Number() int { return d.number }

// Character returns the opponent of the duel.
func (d Duel) Character() Character { return d.character }

// Rewards returns a copy of the three predefined cards granted by the victory
// of the duel (game-rules §58). Changing the returned slice does not affect
// the duel.
func (d Duel) Rewards() []*domain.Card {
	return append([]*domain.Card(nil), d.rewards[:]...)
}
