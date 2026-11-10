package domain

import "fmt"

// Participant is one side of a Duel: a player or an opponent. The engine treats
// both the same way (domain rule 10).
//
// A Participant owns its Life Points and its four zones: Deck, Hand, Field and
// Graveyard (architecture §13). Every CardInstance it holds is owned by it, and
// each instance is in exactly one of the four zones at any time, with
// CardInstance.Location equal to that zone. No rule removes a card from the
// Duel, so the four zone sizes always add up to the size of the deck it
// started with (DeckSize); this is a consequence of the transitions below, not
// a runtime check.
//
// The zones are readable through the accessors, but their mutators are
// unexported: the exported transitions below are the only way to move a card.
// Each is atomic: it checks the source and the capacity of the target before
// changing anything, so a rejected transition leaves all state unchanged.
// Allowed transitions are Deck to Hand, Hand to Field, Hand to Graveyard and
// Field to Graveyard; there is no other (Graveyard to Field is Revival,
// outside the Story Mode MVP).
//
// A Participant is created by NewParticipant; the zero value is not usable. It
// holds no lock: a Duel is a sequential state machine and is serialized by its
// caller.
type Participant struct {
	id         ParticipantID
	lifePoints int
	deck       *DeckZone
	hand       *Hand
	field      *Field
	graveyard  *Graveyard
}

// NewParticipant builds a Participant with the given Life Points, a shuffled
// Deck built from def and empty Hand, Field and Graveyard.
//
// A blank id, an id with leading or trailing whitespace, or Life Points that
// are not positive are rejected with ErrInvalidParticipant. Life Points must
// be positive because a Duel ends immediately when a player's Life Points
// reach 0 or less (mvp-story-mode §1.2), so a Participant cannot start with
// 0 or less. The id uniqueness between the two Participants of a Duel is NOT
// checked here: the Duel (task S3) must reject two Participants with the same
// id, because the InstanceIDs ("<participant id>#<n>") would collide.
//
// It receives the Life Points value, not the Story Level: use
// InitialLifePoints to obtain it. The deck errors of NewDeckZone are
// propagated.
func NewParticipant(id ParticipantID, def DeckDefinition, lifePoints int, s Shuffler) (*Participant, error) {
	if err := validateParticipantID(id); err != nil {
		return nil, err
	}
	if lifePoints <= 0 {
		return nil, fmt.Errorf("%w: life points must be positive, got %d", ErrInvalidParticipant, lifePoints)
	}
	deck, err := NewDeckZone(id, def, s)
	if err != nil {
		return nil, err
	}
	return &Participant{
		id:         id,
		lifePoints: lifePoints,
		deck:       deck,
		hand:       newHand(),
		field:      newField(),
		graveyard:  newGraveyard(),
	}, nil
}

// ID returns the Participant identity.
func (p *Participant) ID() ParticipantID { return p.id }

// LifePoints returns the current Life Points.
func (p *Participant) LifePoints() int { return p.lifePoints }

// Deck returns the Participant's Deck for reading.
func (p *Participant) Deck() *DeckZone { return p.deck }

// Hand returns the Participant's Hand for reading.
func (p *Participant) Hand() *Hand { return p.hand }

// Field returns the Participant's Field for reading.
func (p *Participant) Field() *Field { return p.field }

// Graveyard returns the Participant's Graveyard for reading. It is public: the
// opponent can read it too (game-rules §31).
func (p *Participant) Graveyard() *Graveyard { return p.graveyard }

// Draw moves the top card of the Deck to the Hand and returns it.
//
// It returns ErrDeckEmpty when the Deck has no card and ErrHandFull when the
// Hand already holds HandLimit cards (game-rules §5.2); in both cases nothing
// changes. What a failed draw means for the Duel is not decided here.
func (p *Participant) Draw() (*CardInstance, error) {
	if p.deck.Len() == 0 {
		return nil, fmt.Errorf("%w: participant %q has no card to draw", ErrDeckEmpty, p.id)
	}
	card := p.deck.cards[0]

	err := p.transfer(card, LocationHand,
		p.hand.canAdd,
		func() error { _, err := p.deck.draw(); return err },
		p.hand.add,
	)
	if err != nil {
		return nil, err
	}
	return card, nil
}

// DiscardFromHand moves the card with the given id from the Hand to the
// Graveyard (game-rules §31). It returns ErrCardNotFound, leaving all state
// unchanged, when the card is not in this Participant's Hand.
func (p *Participant) DiscardFromHand(id InstanceID) error {
	card, ok := p.hand.Get(id)
	if !ok {
		return fmt.Errorf("%w: %q is not in the hand of %q", ErrCardNotFound, id, p.id)
	}

	return p.transfer(card, LocationGraveyard,
		nil, // the Graveyard has no capacity limit
		func() error { _, err := p.hand.remove(id); return err },
		p.addToGraveyard,
	)
}

// PlaceFromHand moves the card with the given id from the Hand to the Field.
// The Field area is chosen by the card type (see Field).
//
// It returns ErrCardNotFound when the card is not in the Hand and ErrFieldFull
// when its area has no free slot; in both cases all state is unchanged.
func (p *Participant) PlaceFromHand(id InstanceID) error {
	card, ok := p.hand.Get(id)
	if !ok {
		return fmt.Errorf("%w: %q is not in the hand of %q", ErrCardNotFound, id, p.id)
	}

	return p.transfer(card, LocationField,
		p.field.canAdd,
		func() error { _, err := p.hand.remove(id); return err },
		p.field.add,
	)
}

// SendFieldCardToGraveyard moves the card with the given id from the Field to
// the Graveyard (game-rules §31). It returns ErrCardNotFound, leaving all state
// unchanged, when the card is not on this Participant's Field.
func (p *Participant) SendFieldCardToGraveyard(id InstanceID) error {
	card, ok := p.field.Get(id)
	if !ok {
		return fmt.Errorf("%w: %q is not on the field of %q", ErrCardNotFound, id, p.id)
	}

	return p.transfer(card, LocationGraveyard,
		nil, // the Graveyard has no capacity limit
		func() error { _, err := p.field.remove(id); return err },
		p.addToGraveyard,
	)
}

// addToGraveyard adapts the infallible Graveyard.add to the put step of
// transfer.
func (p *Participant) addToGraveyard(card *CardInstance) error {
	p.graveyard.add(card)
	return nil
}

// transfer is the single place where a card changes zone. The caller has
// already found card in the source zone. It runs in three steps so that a
// rejected transition changes nothing:
//
//  1. canPut checks, without mutating, that the target zone can take the card;
//     nil means the target has no capacity limit (the Graveyard);
//  2. take removes the card from the source and put adds it to the target;
//     after step 1 neither can fail;
//  3. the card's Location is set to the target zone. This is the only call of
//     the single private writer, so Location always equals the containing zone.
//
// A failure in step 2 would mean the check of step 1 disagrees with the
// mutators, which is a programming error; it is returned instead of ignored.
func (p *Participant) transfer(
	card *CardInstance,
	to Location,
	canPut func(*CardInstance) error,
	take func() error,
	put func(*CardInstance) error,
) error {
	if canPut != nil {
		if err := canPut(card); err != nil {
			return err
		}
	}
	if err := take(); err != nil {
		return fmt.Errorf("transfer of %q to %s failed after validation: %w", card.ID(), to, err)
	}
	if err := put(card); err != nil {
		return fmt.Errorf("transfer of %q to %s failed after validation: %w", card.ID(), to, err)
	}

	card.setLocation(to)
	return nil
}
