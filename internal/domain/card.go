package domain

import (
	"fmt"
	"strings"
)

// CardID is the stable identity of a card definition.
type CardID string

// CardType is the type of a card definition (game-rules §10 "Tipos de cartas").
type CardType int

const (
	// CardTypeMonster represents a Monster card. (game-rules §10.1)
	CardTypeMonster CardType = iota + 1

	// CardTypeMagic represents a Magic card. (game-rules §10.2)
	CardTypeMagic

	// CardTypeTrap represents a Trap card. (game-rules §10.3)
	CardTypeTrap
)

// String returns the card type name used by the game rules.
func (t CardType) String() string {
	switch t {
	case CardTypeMonster:
		return "Monster"
	case CardTypeMagic:
		return "Magic"
	case CardTypeTrap:
		return "Trap"
	default:
		return fmt.Sprintf("CardType(%d)", int(t))
	}
}

// valid reports whether t is one of the card types defined by the game rules.
func (t CardType) valid() bool {
	switch t {
	case CardTypeMonster, CardTypeMagic, CardTypeTrap:
		return true
	default:
		return false
	}
}

// Card is the immutable definition of a card (architecture §5.4, glossary §2).
//
// A Card describes what the card is, never the runtime state of a specific copy
// during a Duel. Runtime changes caused by duel effects (modifiers, positions,
// attacks) belong to CardInstance and other duel state; the definition is never
// mutated (domain rule 7).
//
// The fields are private so the compiler, not a convention, enforces that rule:
// the only way to obtain a Card is NewCard, and the only way to read it is
// through its accessors.
type Card struct {
	id     CardID
	name   string
	typ    CardType
	atk    int
	def    int
	weight int
}

// NewCard validates and builds a card definition.
//
// Only constraints stated by the game rules or approved for the Story Mode MVP
// (docs/mvp-story-mode.md) are enforced:
//   - ID and Name are required and must not be blank, since a card must be
//     identifiable and named (architecture §5.4, glossary §2);
//   - Type must be Monster, Magic or Trap (game-rules §10);
//   - Monster: ATK and DEF must not be negative (game-rules §63) and Weight must
//     be at least 1 (game-rules §16);
//   - Magic: ATK, DEF and Weight are not applicable and must be zero.
//
// Trap cards are outside the Story Mode MVP; their validation is unchanged
// (the Monster constraints above apply) and is not decided here.
//
// The game rules define no upper bound for ATK, DEF or Weight, so none is
// enforced here.
func NewCard(id CardID, name string, typ CardType, atk, def, weight int) (*Card, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, fmt.Errorf("%w: id must not be blank", ErrInvalidCard)
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("%w: name must not be blank", ErrInvalidCard)
	}
	if !typ.valid() {
		return nil, fmt.Errorf("%w: unsupported type %s", ErrInvalidCard, typ)
	}

	if typ == CardTypeMagic {
		if atk != 0 || def != 0 || weight != 0 {
			return nil, fmt.Errorf("%w: magic must not define atk, def or weight, got %d, %d, %d",
				ErrInvalidCard, atk, def, weight)
		}
	} else {
		if atk < 0 {
			return nil, fmt.Errorf("%w: atk must not be negative, got %d", ErrInvalidCard, atk)
		}
		if def < 0 {
			return nil, fmt.Errorf("%w: def must not be negative, got %d", ErrInvalidCard, def)
		}
		if weight < 1 {
			return nil, fmt.Errorf("%w: weight must be at least 1, got %d", ErrInvalidCard, weight)
		}
	}

	return &Card{
		id:     id,
		name:   name,
		typ:    typ,
		atk:    atk,
		def:    def,
		weight: weight,
	}, nil
}

// ID returns the card identity.
func (c *Card) ID() CardID { return c.id }

// Name returns the card name.
func (c *Card) Name() string { return c.name }

// Type returns the card type.
func (c *Card) Type() CardType { return c.typ }

// ATK returns the base attack value of the card.
func (c *Card) ATK() int { return c.atk }

// DEF returns the base defense value of the card.
func (c *Card) DEF() int { return c.def }

// Weight returns the Weight used to determine the tributes required to Summon
// the card (game-rules §16).
func (c *Card) Weight() int { return c.weight }
