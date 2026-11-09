package domain

// Graveyard is the permanent zone of the cards a Participant has lost or used
// (game-rules §31, glossary "Zones"). It is public: both Participants can read
// it. Each Participant has its own Graveyard, which holds only the cards that
// Participant owns.
//
// Cards keep their arrival order (game-rules §31.1) and are never removed from
// it: the Graveyard has no capacity limit and no remove operation. The mutator
// is unexported (only the owning Participant moves cards) and readers return
// copies. A Graveyard is created by its Participant; the zero value is not part
// of any Participant and must not be used.
type Graveyard struct {
	cards []*CardInstance
}

// newGraveyard builds an empty Graveyard.
func newGraveyard() *Graveyard {
	return &Graveyard{}
}

// Len returns the number of cards in the Graveyard.
func (g *Graveyard) Len() int { return len(g.cards) }

// Cards returns a copy of the cards in the Graveyard, in arrival order.
// Changing the returned slice does not affect the Graveyard.
func (g *Graveyard) Cards() []*CardInstance {
	return append([]*CardInstance(nil), g.cards...)
}

// Get returns the instance with the given id when it is in the Graveyard.
func (g *Graveyard) Get(id InstanceID) (*CardInstance, bool) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, true
		}
	}
	return nil, false
}

// add appends card after the cards already in the Graveyard. The Graveyard is
// unbounded, so adding never fails.
func (g *Graveyard) add(card *CardInstance) {
	g.cards = append(g.cards, card)
}
