package domain

import "fmt"

// HandLimit is the maximum number of cards a Hand can hold (game-rules §5.2).
const HandLimit = 4

// Hand is the zone of the cards a Participant holds and can play
// (architecture §13, glossary "Zones"). It holds at most HandLimit cards
// (game-rules §5.2) and keeps them in the order they arrived.
//
// The mutators are unexported: only the owning Participant moves cards between
// zones. Readers return copies, so callers cannot change the zone through them.
// A Hand is created by its Participant; the zero value is not part of any
// Participant and must not be used.
type Hand struct {
	cards []*CardInstance
}

// newHand builds an empty Hand.
func newHand() *Hand {
	return &Hand{cards: make([]*CardInstance, 0, HandLimit)}
}

// Len returns the number of cards in the Hand.
func (h *Hand) Len() int { return len(h.cards) }

// Full reports whether the Hand holds HandLimit cards.
func (h *Hand) Full() bool { return len(h.cards) >= HandLimit }

// Cards returns a copy of the cards in the Hand, in insertion order. Changing
// the returned slice does not affect the Hand.
func (h *Hand) Cards() []*CardInstance {
	return append([]*CardInstance(nil), h.cards...)
}

// Get returns the instance with the given id when it is in the Hand.
func (h *Hand) Get(id InstanceID) (*CardInstance, bool) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, true
		}
	}
	return nil, false
}

// canAdd reports why card cannot be added, without changing the Hand. A full
// Hand rejects with ErrHandFull (game-rules §5.2).
func (h *Hand) canAdd(card *CardInstance) error {
	if h.Full() {
		return fmt.Errorf("%w: hand holds %d cards, limit is %d", ErrHandFull, len(h.cards), HandLimit)
	}
	return nil
}

// add appends card to the Hand. When the Hand is full it returns ErrHandFull
// and leaves the Hand unchanged.
func (h *Hand) add(card *CardInstance) error {
	if err := h.canAdd(card); err != nil {
		return err
	}
	h.cards = append(h.cards, card)
	return nil
}

// remove takes the instance with the given id out of the Hand, keeping the
// order of the remaining cards. It returns ErrCardNotFound, leaving the Hand
// unchanged, when the instance is not in the Hand.
func (h *Hand) remove(id InstanceID) (*CardInstance, error) {
	for i, card := range h.cards {
		if card.ID() == id {
			h.cards = append(h.cards[:i], h.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, fmt.Errorf("%w: %q is not in the hand", ErrCardNotFound, id)
}
