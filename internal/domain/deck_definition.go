package domain

import "fmt"

const (
	// DeckSize is the exact number of cards a deck must contain
	// (game-rules §4).
	DeckSize = 30

	// DeckCopiesLimit is the maximum number of copies of the same card in a
	// deck (game-rules §4.1).
	DeckCopiesLimit = 2
)

// DeckEntry is a card definition together with the number of copies used in the
// deck.
type DeckEntry struct {
	Card   *Card
	Copies int
}

// DeckDefinition is the composition of a deck: which cards and how many copies
// of each. It is the persistent description of a deck used as duel input, not
// the runtime deck of a Duel.
//
// The runtime deck (the ordered cards being drawn during a Duel) is DeckZone.
// Keeping the two names distinct avoids mixing deck composition with duel
// state.
//
// A DeckDefinition does not create card instances: expanding the composition
// into copies requires an owner and an identity, which belong to the duel.
//
// The entries are private so a validated deck cannot be changed afterwards.
// Duel initialization must only receive decks built with NewDeckDefinition; the
// zero value is reported invalid by Validate.
type DeckDefinition struct {
	entries []DeckEntry
}

// NewDeckDefinition copies the entries and validates the composition against
// the deck rules (see Validate). It returns the zero DeckDefinition and the
// validation error when the composition is invalid, so an invalid deck is never
// handed out.
func NewDeckDefinition(entries []DeckEntry) (DeckDefinition, error) {
	deck := DeckDefinition{entries: append([]DeckEntry(nil), entries...)}
	if err := deck.Validate(); err != nil {
		return DeckDefinition{}, err
	}
	return deck, nil
}

// Entries returns a copy of the deck entries. Changing the returned slice does
// not affect the deck.
func (d DeckDefinition) Entries() []DeckEntry {
	return append([]DeckEntry(nil), d.entries...)
}

// TotalCards returns the number of card copies described by the composition,
// counting every entry. It does not validate the composition.
func (d DeckDefinition) TotalCards() int {
	total := 0
	for _, entry := range d.entries {
		total += entry.Copies
	}
	return total
}

// Validate checks the composition against the deck rules.
//
// Rules enforced (game-rules §4 and §4.1):
//   - every entry references an existing card;
//   - every entry has at least one copy and at most DeckCopiesLimit copies;
//   - the same card does not appear in more than one entry;
//   - the total number of cards is exactly DeckSize.
//
// Duplicate entries are rejected instead of being merged, because merging would
// hide an invalid composition where the same card is counted twice.
//
// Returned errors wrap one of ErrInvalidDeck, ErrTooManyCopies, ErrDuplicateCard
// or ErrDeckSize with the detail that explains the violation, so callers can
// classify the failure with errors.Is.
func (d DeckDefinition) Validate() error {
	seen := make(map[CardID]struct{}, len(d.entries))

	for i, entry := range d.entries {
		if entry.Card == nil {
			return fmt.Errorf("%w: entry %d has no card", ErrInvalidDeck, i)
		}
		if entry.Copies < 1 {
			return fmt.Errorf("%w: entry %d has %d copies", ErrInvalidDeck, i, entry.Copies)
		}
		if entry.Copies > DeckCopiesLimit {
			return fmt.Errorf("%w: card %q has %d copies, limit is %d",
				ErrTooManyCopies, entry.Card.ID(), entry.Copies, DeckCopiesLimit)
		}

		id := entry.Card.ID()
		if _, duplicated := seen[id]; duplicated {
			return fmt.Errorf("%w: card %q appears in more than one entry", ErrDuplicateCard, id)
		}
		seen[id] = struct{}{}
	}

	if total := d.TotalCards(); total != DeckSize {
		return fmt.Errorf("%w: deck has %d cards, expected %d", ErrDeckSize, total, DeckSize)
	}

	return nil
}
