package story

import (
	"fmt"
	"strings"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// Character is the named opponent the player faces in a Story Mode duel
// (glossary §12 "Character", PR-STORY-003).
//
// A Character carries the deck the Character plays with (architecture §9).
// The Character is an immutable definition: the same definition is reused
// every time the duel is played, and the runtime opponent of a Duel is
// materialized by the Duel Engine (tasks S3 to S7), not by this type.
type Character struct {
	id   string
	name string
	deck domain.DeckDefinition
}

// NewCharacter validates and builds a Character.
//
// The id and the name are required and must not be blank, since a Character
// must be identifiable and named (glossary §12: "A named opponent used by
// Story Mode"). The deck must be a valid deck composition (game-rules §4 and
// §4.1): the deck of an opponent is the input of the duel, so it is
// validated with the same rules as every deck.
//
// The result is an immutable value: the stored deck is a DeckDefinition
// value, whose entries are not reachable for mutation.
func NewCharacter(id, name string, deck domain.DeckDefinition) (Character, error) {
	if strings.TrimSpace(id) == "" {
		return Character{}, fmt.Errorf("%w: id must not be blank", ErrInvalidCharacter)
	}
	if strings.TrimSpace(name) == "" {
		return Character{}, fmt.Errorf("%w: name must not be blank", ErrInvalidCharacter)
	}
	if err := deck.Validate(); err != nil {
		return Character{}, fmt.Errorf("%w: deck is invalid: %w", ErrInvalidCharacter, err)
	}
	return Character{id: id, name: name, deck: deck}, nil
}

// ID returns the Character identity.
func (c Character) ID() string { return c.id }

// Name returns the Character's display name.
func (c Character) Name() string { return c.name }

// Deck returns the deck the Character plays with.
func (c Character) Deck() domain.DeckDefinition { return c.deck }
