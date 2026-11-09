package domain

import "fmt"

// DeckZone is the runtime deck of a Participant: the ordered cards still to be
// drawn (architecture §13, glossary "Deck Zone"). It is distinct from
// DeckDefinition, which only describes the composition.
//
// The top of the deck is index 0. Cards are drawn from the top. The mutator is
// unexported (only the owning Participant moves cards) and readers return
// copies. A DeckZone is created by NewDeckZone; the zero value is not usable.
type DeckZone struct {
	cards []*CardInstance
}

// NewDeckZone expands a deck composition into card instances owned by owner and
// shuffles them with s.
//
// The definition is validated with DeckDefinition.Validate and its sentinel
// errors are propagated (so the zero DeckDefinition is rejected). A blank owner
// or one with surrounding whitespace is rejected with ErrInvalidParticipant and a nil shuffler with ErrInvalidDeck.
//
// The entries are expanded in order, each one Copies times, and the instances
// get the ids "<owner>#<n>" with n = 1..DeckSize in that expansion order. Ids
// are assigned before the shuffle, so only the order of the deck is random. Two
// copies of the same card are two distinct instances (game-rules §64,
// architecture §5.5). Shuffler.Shuffle is called exactly once, with the number of cards.
func NewDeckZone(owner ParticipantID, def DeckDefinition, s Shuffler) (*DeckZone, error) {
	if err := def.Validate(); err != nil {
		return nil, err
	}
	if err := validateParticipantID(owner); err != nil {
		return nil, fmt.Errorf("deck owner: %w", err)
	}
	if s == nil {
		return nil, fmt.Errorf("%w: shuffler must not be nil", ErrInvalidDeck)
	}

	cards := make([]*CardInstance, 0, def.TotalCards())
	for _, entry := range def.Entries() {
		for range entry.Copies {
			id := InstanceID(fmt.Sprintf("%s#%d", owner, len(cards)+1))
			inst, err := newCardInstance(id, entry.Card, owner)
			if err != nil {
				return nil, err
			}
			cards = append(cards, inst)
		}
	}

	s.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return &DeckZone{cards: cards}, nil
}

// Len returns the number of cards left in the Deck.
func (d *DeckZone) Len() int { return len(d.cards) }

// Cards returns a copy of the cards in the Deck, top first. Changing the
// returned slice does not affect the Deck.
func (d *DeckZone) Cards() []*CardInstance {
	return append([]*CardInstance(nil), d.cards...)
}

// draw takes the top card out of the Deck. When the Deck is empty it returns
// ErrDeckEmpty and leaves the Deck unchanged.
func (d *DeckZone) draw() (*CardInstance, error) {
	if len(d.cards) == 0 {
		return nil, fmt.Errorf("%w: no card to draw", ErrDeckEmpty)
	}
	top := d.cards[0]
	d.cards = d.cards[1:]
	return top, nil
}
