package domain

import (
	"fmt"
	"strings"
)

// ParticipantID is the identity of a Participant of a Duel.
//
// It is stored exactly as given. NewParticipant and NewDeckZone reject a blank
// id and an id with leading or trailing whitespace (see validateParticipantID);
// whitespace inside the id is accepted.
//
// Uniqueness is NOT checked here: nothing at this level knows the other
// Participant. The Duel (task S3) must reject two Participants with the same
// id before starting, because the InstanceIDs are "<participant id>#<n>" and
// would collide.
type ParticipantID string

// validateParticipantID checks that id can identify a Participant: it must not
// be blank and must not have leading or trailing whitespace. Violations wrap
// ErrInvalidParticipant.
//
// Surrounding whitespace is rejected instead of trimmed because a
// ParticipantID is an identity, used as a lookup key and as the prefix of
// every InstanceID: trimming would silently change the string the caller
// passed. This matches the other identities of the domain (CardID, InstanceID),
// which are stored as given and only rejected when blank.
func validateParticipantID(id ParticipantID) error {
	if strings.TrimSpace(string(id)) == "" {
		return fmt.Errorf("%w: id must not be blank", ErrInvalidParticipant)
	}
	if string(id) != strings.TrimSpace(string(id)) {
		return fmt.Errorf("%w: id %q has surrounding whitespace", ErrInvalidParticipant, string(id))
	}
	return nil
}

// InstanceID is the identity of one CardInstance. Two copies of the same card
// definition are two different instances with different InstanceIDs
// (game-rules §64, architecture §5.5).
type InstanceID string

// Location is the zone that currently holds a CardInstance (architecture §13,
// glossary "Zones"). The zero value is not a valid location.
type Location int

const (
	// LocationDeck means the instance is in its owner's Deck.
	LocationDeck Location = iota + 1

	// LocationHand means the instance is in its owner's Hand.
	LocationHand

	// LocationField means the instance is on its owner's Field.
	LocationField

	// LocationGraveyard means the instance is in its owner's Graveyard
	// (game-rules §31).
	LocationGraveyard
)

// String returns the zone name used by the game rules.
func (l Location) String() string {
	switch l {
	case LocationDeck:
		return "Deck"
	case LocationHand:
		return "Hand"
	case LocationField:
		return "Field"
	case LocationGraveyard:
		return "Graveyard"
	default:
		return fmt.Sprintf("Location(%d)", int(l))
	}
}

// CardInstance is a concrete copy of a card during a Duel (architecture §5.5,
// game-rules §64).
//
// The Card is the shared, immutable definition and is never mutated by the
// Duel (domain rule 7); everything that changes during a Duel belongs to the
// instance. The instance has one owner for its whole life and is, at any
// moment, in exactly one zone of that owner. The location is stored here
// because architecture §5.5 lists it as part of the instance state; it has a
// single writer, the transfer between the owner's zones (see Participant),
// which is the only code allowed to change it.
//
// A Monster on the Field also carries its position and face (game-rules
// §18 and §21), meaningful only while it is in the Field; the placement and
// the position-change transitions are the only writers. Whether the Monster
// changed its position in the current turn is turn-specific state that the
// Duel resets when a new turn begins (architecture §12, game-rules §21.1).
//
// A Magic Equipment in the Field also carries its link to the equipped
// Monster (game-rules §35.2.4), meaningful only while it is linked; the
// activation and the transition that removes the Monster from the Field are
// the only writers. The link is not turn-specific: the bonus it grants
// persists while the Equipment stays equipped (game-rules §38), so the Duel
// never resets it. The Monster holds no reference to its Equipment: the
// +800 ATK bonus is derived from the links of the Field, and the base Card
// definition is never mutated (domain rule 7).
//
// Instances are created only by the deck expansion of DeckZone, so the zero
// value is not usable and the constructor is not exported.
type CardInstance struct {
	id       InstanceID
	card     *Card
	owner    ParticipantID
	location Location
	position Position
	face     Face
	changed  bool // position changed during the current turn (game-rules §21.1)

	// equipmentLink is the id of the Monster a Magic Equipment is linked to
	// (game-rules §35.2.4); the zero value means "not linked". Only an
	// Equipment in the Field can hold a link.
	equipmentLink InstanceID
}

// newCardInstance validates and builds an instance that starts in the Deck.
// A nil card, a blank id and a blank owner are rejected with ErrInvalidInstance.
func newCardInstance(id InstanceID, card *Card, owner ParticipantID) (*CardInstance, error) {
	if card == nil {
		return nil, fmt.Errorf("%w: card must not be nil", ErrInvalidInstance)
	}
	if strings.TrimSpace(string(id)) == "" {
		return nil, fmt.Errorf("%w: id must not be blank", ErrInvalidInstance)
	}
	if strings.TrimSpace(string(owner)) == "" {
		return nil, fmt.Errorf("%w: owner must not be blank", ErrInvalidInstance)
	}
	return &CardInstance{id: id, card: card, owner: owner, location: LocationDeck}, nil
}

// ID returns the instance identity.
func (c *CardInstance) ID() InstanceID { return c.id }

// Card returns the shared card definition. It is never mutated by a Duel.
func (c *CardInstance) Card() *Card { return c.card }

// Owner returns the Participant that owns the instance.
func (c *CardInstance) Owner() ParticipantID { return c.owner }

// Location returns the zone that currently holds the instance.
func (c *CardInstance) Location() Location { return c.location }

// Position returns the Monster's position on the Field (game-rules §21). It
// is meaningful only for a Monster in LocationField; for any other card or
// zone it returns the zero Position.
func (c *CardInstance) Position() Position { return c.position }

// Face returns the Monster's face on the Field (game-rules §18). It is
// meaningful only for a Monster in LocationField; for any other card or zone
// it returns the zero Face.
func (c *CardInstance) Face() Face { return c.face }

// positionChanged reports whether the Monster changed its position during the
// current turn (game-rules §21.1, mvp-story-mode S4). The Duel resets it when
// a new turn begins.
func (c *CardInstance) positionChanged() bool { return c.changed }

// setFieldState is the single writer of the position and the face. Only the
// placement and the position-change transitions of the owning Participant may
// call it. The pair must be coherent (see validPositionFace): a face-down
// Monster is always in Defense (game-rules §18.2). The callers validate the
// pair before calling; an incoherent pair here is a programming error, so it
// panics instead of silently storing an invalid state. The pair is meaningful
// only while the instance is a Monster on the Field.
func (c *CardInstance) setFieldState(position Position, face Face) {
	if !validPositionFace(position, face) {
		panic(fmt.Sprintf("setFieldState of %q with incoherent pair %s %s", c.id, face, position))
	}
	c.position = position
	c.face = face
}

// markPositionChanged records a position change of the current turn
// (game-rules §21.1).
func (c *CardInstance) markPositionChanged() { c.changed = true }

// resetTurnState clears the turn-specific state of the instance
// (architecture §12). The Duel calls it, through the owning Participant, when
// a new turn begins.
func (c *CardInstance) resetTurnState() { c.changed = false }

// setLocation is the single writer of the location. Only the transfer between
// zones may call it, so the location always equals the containing zone.
func (c *CardInstance) setLocation(l Location) { c.location = l }
