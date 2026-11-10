package domain

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

// mixedDeck builds a valid composition with Monsters, Magic and Trap cards, in
// that order: 8 monsters x2, 5 magic x2 and 2 traps x2 (30 cards).
func mixedDeck(t *testing.T) DeckDefinition {
	t.Helper()

	var entries []DeckEntry
	for n := 1; n <= 8; n++ {
		entries = append(entries, DeckEntry{Card: deckCard(t, n), Copies: 2})
	}
	for n := 1; n <= 5; n++ {
		entries = append(entries, DeckEntry{Card: magicCard(t, n), Copies: 2})
	}
	for n := 1; n <= 2; n++ {
		entries = append(entries, DeckEntry{Card: trapCard(t, n), Copies: 2})
	}
	deck, err := NewDeckDefinition(entries)
	if err != nil {
		t.Fatalf("mixedDeck is not valid: %v", err)
	}
	return deck
}

// newTestParticipant builds a Participant with the unshuffled validDeck (30
// Monsters) so the draw order is the expansion order.
func newTestParticipant(t *testing.T, id ParticipantID) *Participant {
	t.Helper()

	p, err := NewParticipant(id, validDeck(t), 4000, noopShuffler{})
	if err != nil {
		t.Fatalf("NewParticipant() returned unexpected error: %v", err)
	}
	return p
}

// newMixedParticipant is like newTestParticipant with mixedDeck.
func newMixedParticipant(t *testing.T, id ParticipantID) *Participant {
	t.Helper()

	p, err := NewParticipant(id, mixedDeck(t), 4000, noopShuffler{})
	if err != nil {
		t.Fatalf("NewParticipant() returned unexpected error: %v", err)
	}
	return p
}

// drawN draws n cards and returns them in draw order.
func drawN(t *testing.T, p *Participant, n int) []*CardInstance {
	t.Helper()

	cards := make([]*CardInstance, 0, n)
	for range n {
		card, err := p.Draw()
		if err != nil {
			t.Fatalf("Draw() returned unexpected error: %v", err)
		}
		cards = append(cards, card)
	}
	return cards
}

// state is a comparable picture of every zone and every location.
type state struct {
	Deck, Hand, Monsters, MagicTraps, Graveyard []InstanceID
	Locations                                   map[InstanceID]Location
}

func stateOf(p *Participant) state {
	s := state{
		Deck:       ids(p.Deck().Cards()),
		Hand:       ids(p.Hand().Cards()),
		Monsters:   ids(p.Field().Monsters()),
		MagicTraps: ids(p.Field().MagicTraps()),
		Graveyard:  ids(p.Graveyard().Cards()),
		Locations:  make(map[InstanceID]Location),
	}
	for _, zone := range [][]*CardInstance{
		p.Deck().Cards(), p.Hand().Cards(), p.Field().Monsters(), p.Field().MagicTraps(), p.Graveyard().Cards(),
	} {
		for _, c := range zone {
			s.Locations[c.ID()] = c.Location()
		}
	}
	return s
}

// assertInvariants checks the Participant invariants: the zone sizes add up to
// the deck size, every instance is in exactly one zone, its Location equals
// that zone, its owner is the Participant and the zone limits hold.
func assertInvariants(t *testing.T, p *Participant) {
	t.Helper()

	zones := []struct {
		name  string
		want  Location
		cards []*CardInstance
	}{
		{"deck", LocationDeck, p.Deck().Cards()},
		{"hand", LocationHand, p.Hand().Cards()},
		{"monsters", LocationField, p.Field().Monsters()},
		{"magic/trap", LocationField, p.Field().MagicTraps()},
		{"graveyard", LocationGraveyard, p.Graveyard().Cards()},
	}
	seen := make(map[InstanceID]string)
	total := 0
	for _, z := range zones {
		for _, c := range z.cards {
			total++
			if prev, dup := seen[c.ID()]; dup {
				t.Errorf("instance %q is in %s and in %s", c.ID(), prev, z.name)
			}
			seen[c.ID()] = z.name
			if c.Location() != z.want {
				t.Errorf("instance %q is in the %s but Location() = %s, want %s", c.ID(), z.name, c.Location(), z.want)
			}
			if c.Owner() != p.ID() {
				t.Errorf("instance %q owner = %q, want %q", c.ID(), c.Owner(), p.ID())
			}
		}
	}
	if total != DeckSize {
		t.Errorf("sum of zone sizes = %d, want %d", total, DeckSize)
	}
	if p.Hand().Len() > HandLimit {
		t.Errorf("hand holds %d cards, limit is %d", p.Hand().Len(), HandLimit)
	}
	if n := len(p.Field().Monsters()); n > MonsterSlots {
		t.Errorf("monster area holds %d cards, limit is %d", n, MonsterSlots)
	}
	if n := len(p.Field().MagicTraps()); n > MagicTrapSlots {
		t.Errorf("magic/trap area holds %d cards, limit is %d", n, MagicTrapSlots)
	}
}

// assertUnchanged fails when the Participant state differs from before.
func assertUnchanged(t *testing.T, p *Participant, before state) {
	t.Helper()

	if after := stateOf(p); !reflect.DeepEqual(before, after) {
		t.Errorf("state changed by a rejected transition:\nbefore %+v\nafter  %+v", before, after)
	}
	assertInvariants(t, p)
}

func TestNewParticipant_StartsWithAShuffledDeckAndEmptyOtherZones(t *testing.T) {
	// Arrange
	fake := &fakeShuffler{}

	// Act
	p, err := NewParticipant("p1", validDeck(t), 4000, fake)

	// Assert
	if err != nil {
		t.Fatalf("NewParticipant() returned unexpected error: %v", err)
	}
	if p.ID() != "p1" || p.LifePoints() != 4000 {
		t.Errorf("ID, LifePoints = %q, %d, want p1, 4000", p.ID(), p.LifePoints())
	}
	if p.Deck().Len() != DeckSize {
		t.Errorf("Deck().Len() = %d, want %d", p.Deck().Len(), DeckSize)
	}
	if p.Hand().Len() != 0 || len(p.Field().Monsters()) != 0 || len(p.Field().MagicTraps()) != 0 || p.Graveyard().Len() != 0 {
		t.Errorf("hand, field and graveyard must start empty")
	}
	if fake.calls != 1 {
		t.Errorf("Shuffle calls = %d, want 1", fake.calls)
	}
	assertInvariants(t, p)
}

func TestNewParticipant_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		id         ParticipantID
		def        func(t *testing.T) DeckDefinition
		lifePoints int
		shuffler   Shuffler
		wantErr    error
	}{
		{name: "blank id", id: " ", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "empty id", id: "", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "tab id", id: "\t", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "spaces id", id: "   ", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "leading space id", id: " p1", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "trailing space id", id: "p1 ", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "surrounded id", id: " p1 ", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "tab and newline id", id: "\tp1\n", def: validDeck, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "zero life points", id: "p1", def: validDeck, lifePoints: 0, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "negative life points", id: "p1", def: validDeck, lifePoints: -1, shuffler: noopShuffler{}, wantErr: ErrInvalidParticipant},
		{name: "nil shuffler", id: "p1", def: validDeck, lifePoints: 4000, shuffler: nil, wantErr: ErrInvalidDeck},
		{name: "zero deck definition", id: "p1", def: func(*testing.T) DeckDefinition { return DeckDefinition{} }, lifePoints: 4000, shuffler: noopShuffler{}, wantErr: ErrDeckSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			def := tt.def(t)

			// Act
			p, err := NewParticipant(tt.id, def, tt.lifePoints, tt.shuffler)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewParticipant() error = %v, want %v", err, tt.wantErr)
			}
			if p != nil {
				t.Errorf("NewParticipant() = %v, want nil on error", p)
			}
		})
	}
}

func TestNewParticipant_AcceptsIdsWithoutSurroundingWhitespace(t *testing.T) {
	tests := []struct {
		name string
		id   ParticipantID
	}{
		{name: "short id", id: "p1"},
		{name: "hyphenated id", id: "player-1"},
		{name: "internal space is accepted", id: "player one"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			def := validDeck(t)

			// Act
			p, err := NewParticipant(tt.id, def, 4000, noopShuffler{})

			// Assert
			if err != nil {
				t.Fatalf("NewParticipant(%q) returned unexpected error: %v", tt.id, err)
			}
			if p.ID() != tt.id {
				t.Errorf("ID() = %q, want %q (stored as given)", p.ID(), tt.id)
			}
			for _, inst := range p.Deck().Cards() {
				if inst.Owner() != tt.id {
					t.Errorf("instance %q owner = %q, want %q", inst.ID(), inst.Owner(), tt.id)
				}
			}
		})
	}
}

func TestNewParticipant_AcceptsTheLifePointsOfEachSupportedLevel(t *testing.T) {
	for level := 1; level <= MaxSupportedLevel; level++ {
		t.Run(fmt.Sprintf("level %d", level), func(t *testing.T) {
			// Arrange
			lp, err := InitialLifePoints(level)
			if err != nil {
				t.Fatalf("InitialLifePoints() returned unexpected error: %v", err)
			}

			// Act
			p, err := NewParticipant("p1", validDeck(t), lp, noopShuffler{})

			// Assert
			if err != nil {
				t.Fatalf("NewParticipant() returned unexpected error: %v", err)
			}
			if p.LifePoints() != lp {
				t.Errorf("LifePoints() = %d, want %d", p.LifePoints(), lp)
			}
		})
	}
}

func TestParticipant_DrawMovesTheTopCardFromDeckToHand(t *testing.T) {
	// Arrange
	p := newTestParticipant(t, "p1")
	top := p.Deck().Cards()[0]

	// Act
	drawn, err := p.Draw()

	// Assert
	if err != nil {
		t.Fatalf("Draw() returned unexpected error: %v", err)
	}
	if drawn != top {
		t.Errorf("Draw() = %q, want the top card %q", drawn.ID(), top.ID())
	}
	if drawn.Location() != LocationHand {
		t.Errorf("Location() = %s, want %s", drawn.Location(), LocationHand)
	}
	if got, ok := p.Hand().Get(drawn.ID()); !ok || got != drawn {
		t.Errorf("the drawn card is not in the hand")
	}
	if p.Deck().Len() != DeckSize-1 {
		t.Errorf("Deck().Len() = %d, want %d", p.Deck().Len(), DeckSize-1)
	}
	assertInvariants(t, p)
}

func TestParticipant_DrawFromEmptyDeckFailsAndStaysUnchanged(t *testing.T) {
	// Arrange
	p := newTestParticipant(t, "p1")
	for p.Deck().Len() > 0 {
		for _, card := range drawN(t, p, min(HandLimit, p.Deck().Len())) {
			if err := p.DiscardFromHand(card.ID()); err != nil {
				t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
			}
		}
	}
	before := stateOf(p)

	// Act
	drawn, err := p.Draw()

	// Assert
	if !errors.Is(err, ErrDeckEmpty) {
		t.Fatalf("Draw() error = %v, want ErrDeckEmpty", err)
	}
	if drawn != nil {
		t.Errorf("Draw() = %v, want nil on error", drawn)
	}
	assertUnchanged(t, p, before)
	if p.Graveyard().Len() != DeckSize {
		t.Errorf("Graveyard().Len() = %d, want %d", p.Graveyard().Len(), DeckSize)
	}
}

func TestParticipant_DrawWithFullHandFailsAndStaysUnchanged(t *testing.T) {
	// Arrange
	p := newTestParticipant(t, "p1")
	drawN(t, p, HandLimit)
	before := stateOf(p)
	top := p.Deck().Cards()[0]

	// Act
	drawn, err := p.Draw()

	// Assert
	if !errors.Is(err, ErrHandFull) {
		t.Fatalf("Draw() error = %v, want ErrHandFull", err)
	}
	if drawn != nil {
		t.Errorf("Draw() = %v, want nil on error", drawn)
	}
	assertUnchanged(t, p, before)
	if top.Location() != LocationDeck || p.Deck().Cards()[0] != top {
		t.Errorf("the top card must stay on top of the deck, Location() = %s", top.Location())
	}
}

func TestParticipant_CardKeepsIdentityAlongDeckHandFieldGraveyard(t *testing.T) {
	// Arrange
	p := newTestParticipant(t, "p1")
	top := p.Deck().Cards()[0]
	id, def := top.ID(), top.Card()

	steps := []struct {
		name string
		act  func() error
		want Location
	}{
		{"Draw", func() error { _, err := p.Draw(); return err }, LocationHand},
		{"PlaceFromHand", func() error { return p.PlaceFromHand(id) }, LocationField},
		{"SendFieldCardToGraveyard", func() error { return p.SendFieldCardToGraveyard(id) }, LocationGraveyard},
	}

	for _, step := range steps {
		// Act
		if err := step.act(); err != nil {
			t.Fatalf("%s returned unexpected error: %v", step.name, err)
		}

		// Assert
		if top.Location() != step.want {
			t.Errorf("after %s Location() = %s, want %s", step.name, top.Location(), step.want)
		}
		if top.ID() != id || top.Card() != def {
			t.Errorf("after %s the instance changed identity or definition", step.name)
		}
		assertInvariants(t, p)
	}
	if got, ok := p.Graveyard().Get(id); !ok || got != top {
		t.Errorf("the same *CardInstance must end in the graveyard")
	}
}

func TestParticipant_DiscardFromHandMovesTheCardToTheGraveyard(t *testing.T) {
	// Arrange
	p := newTestParticipant(t, "p1")
	cards := drawN(t, p, 3)

	// Act
	err := p.DiscardFromHand(cards[1].ID())

	// Assert
	if err != nil {
		t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
	}
	if cards[1].Location() != LocationGraveyard {
		t.Errorf("Location() = %s, want %s", cards[1].Location(), LocationGraveyard)
	}
	if got, want := ids(p.Hand().Cards()), []InstanceID{cards[0].ID(), cards[2].ID()}; !slices.Equal(got, want) {
		t.Errorf("Hand = %v, want %v", got, want)
	}
	assertInvariants(t, p)
}

// fetch draws until it finds a card of the wanted type, discarding the other
// cards drawn on the way, and leaves the found card in the hand.
func fetch(t *testing.T, p *Participant, typ CardType) *CardInstance {
	t.Helper()

	for p.Deck().Len() > 0 {
		card, err := p.Draw()
		if err != nil {
			t.Fatalf("Draw() returned unexpected error: %v", err)
		}
		if card.Card().Type() == typ {
			return card
		}
		if err := p.DiscardFromHand(card.ID()); err != nil {
			t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
		}
	}
	t.Fatalf("no %s left in the deck", typ)
	return nil
}

func TestParticipant_PlaceFromHandPutsEachTypeInItsArea(t *testing.T) {
	tests := []struct {
		typ          CardType
		wantMonsters int
		wantMagicTrp int
	}{
		{CardTypeMonster, 1, 0},
		{CardTypeMagic, 0, 1},
		{CardTypeTrap, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.typ.String(), func(t *testing.T) {
			// Arrange
			p := newMixedParticipant(t, "p1")
			card := fetch(t, p, tt.typ)

			// Act
			err := p.PlaceFromHand(card.ID())

			// Assert
			if err != nil {
				t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
			}
			if card.Location() != LocationField {
				t.Errorf("Location() = %s, want %s", card.Location(), LocationField)
			}
			if got := len(p.Field().Monsters()); got != tt.wantMonsters {
				t.Errorf("len(Monsters()) = %d, want %d", got, tt.wantMonsters)
			}
			if got := len(p.Field().MagicTraps()); got != tt.wantMagicTrp {
				t.Errorf("len(MagicTraps()) = %d, want %d", got, tt.wantMagicTrp)
			}
			if _, inHand := p.Hand().Get(card.ID()); inHand {
				t.Errorf("the placed card is still in the hand")
			}
			assertInvariants(t, p)
		})
	}
}

func TestParticipant_PlaceFromHandIntoFullAreaFailsAndStaysUnchanged(t *testing.T) {
	tests := []struct {
		typ CardType
		cap int
	}{
		{CardTypeMonster, MonsterSlots},
		{CardTypeMagic, MagicTrapSlots},
	}

	for _, tt := range tests {
		t.Run(tt.typ.String(), func(t *testing.T) {
			// Arrange
			p := newMixedParticipant(t, "p1")
			for range tt.cap {
				if err := p.PlaceFromHand(fetch(t, p, tt.typ).ID()); err != nil {
					t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
				}
			}
			extra := fetch(t, p, tt.typ)
			before := stateOf(p)

			// Act
			err := p.PlaceFromHand(extra.ID())

			// Assert
			if !errors.Is(err, ErrFieldFull) {
				t.Fatalf("PlaceFromHand() error = %v, want ErrFieldFull", err)
			}
			assertUnchanged(t, p, before)
			if extra.Location() != LocationHand {
				t.Errorf("Location() = %s, want %s", extra.Location(), LocationHand)
			}
			if _, inHand := p.Hand().Get(extra.ID()); !inHand {
				t.Errorf("the rejected card must stay in the hand")
			}
		})
	}
}

func TestParticipant_SendFieldCardToGraveyardKeepsArrivalOrder(t *testing.T) {
	// Arrange
	p := newTestParticipant(t, "p1")
	cards := drawN(t, p, 2)
	for _, c := range cards {
		if err := p.PlaceFromHand(c.ID()); err != nil {
			t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
		}
	}

	// Act: B arrives before A.
	errB := p.SendFieldCardToGraveyard(cards[1].ID())
	errA := p.SendFieldCardToGraveyard(cards[0].ID())

	// Assert
	if errB != nil || errA != nil {
		t.Fatalf("SendFieldCardToGraveyard() errors = %v, %v, want nil", errB, errA)
	}
	if got, want := ids(p.Graveyard().Cards()), []InstanceID{cards[1].ID(), cards[0].ID()}; !slices.Equal(got, want) {
		t.Errorf("Graveyard = %v, want arrival order %v", got, want)
	}
	if len(p.Field().Monsters()) != 0 {
		t.Errorf("the monster area must be empty")
	}
	assertInvariants(t, p)
}

func TestParticipant_TransitionsRejectAbsentOrForeignIDsAndStayUnchanged(t *testing.T) {
	tests := []struct {
		name string
		act  func(p *Participant, id InstanceID) error
	}{
		{"DiscardFromHand", (*Participant).DiscardFromHand},
		{"PlaceFromHand", (*Participant).PlaceFromHand},
		{"SendFieldCardToGraveyard", (*Participant).SendFieldCardToGraveyard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: p1 has one card in the hand and one on the field; the
			// foreign id is a card of p2 and the deck card is in neither zone.
			p1 := newTestParticipant(t, "p1")
			p2 := newTestParticipant(t, "p2")
			cards := drawN(t, p1, 2)
			if err := p1.PlaceFromHand(cards[0].ID()); err != nil {
				t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
			}
			foreign := drawN(t, p2, 1)[0]
			before1, before2 := stateOf(p1), stateOf(p2)

			absent := []InstanceID{"p1#999", foreign.ID(), p1.Deck().Cards()[0].ID()}
			if tt.name == "SendFieldCardToGraveyard" {
				absent = append(absent, cards[1].ID()) // in the hand, not on the field
			} else {
				absent = append(absent, cards[0].ID()) // on the field, not in the hand
			}

			for _, id := range absent {
				// Act
				err := tt.act(p1, id)

				// Assert
				if !errors.Is(err, ErrCardNotFound) {
					t.Errorf("%s(%q) error = %v, want ErrCardNotFound", tt.name, id, err)
				}
			}
			assertUnchanged(t, p1, before1)
			assertUnchanged(t, p2, before2)
		})
	}
}

func TestParticipant_ConservesTheDeckAcrossALongScriptedSequence(t *testing.T) {
	// Arrange: a deterministic script that keeps the hand at HandLimit, places
	// what fits, discards the rest and empties the field now and then, until
	// the deck runs out.
	p := newMixedParticipant(t, "p1")
	assertInvariants(t, p)
	placed, discarded, sent := 0, 0, 0

	// Act
	for step := 0; p.Deck().Len() > 0 || p.Hand().Len() > 0; step++ {
		for !p.Hand().Full() && p.Deck().Len() > 0 {
			if _, err := p.Draw(); err != nil {
				t.Fatalf("step %d: Draw() returned unexpected error: %v", step, err)
			}
			assertInvariants(t, p)
		}
		for _, card := range p.Hand().Cards() {
			if step%2 == 0 && p.Field().HasRoomFor(card.Card()) {
				if err := p.PlaceFromHand(card.ID()); err != nil {
					t.Fatalf("step %d: PlaceFromHand() returned unexpected error: %v", step, err)
				}
				placed++
			} else if err := p.DiscardFromHand(card.ID()); err != nil {
				t.Fatalf("step %d: DiscardFromHand() returned unexpected error: %v", step, err)
			} else {
				discarded++
			}
			assertInvariants(t, p)
		}
		if step%3 == 0 {
			for _, card := range append(p.Field().Monsters(), p.Field().MagicTraps()...) {
				if err := p.SendFieldCardToGraveyard(card.ID()); err != nil {
					t.Fatalf("step %d: SendFieldCardToGraveyard() returned unexpected error: %v", step, err)
				}
				sent++
				assertInvariants(t, p)
			}
		}
	}

	// Assert: every card ended in the graveyard or on the field, and the script
	// really used the three moves.
	inField := len(p.Field().Monsters()) + len(p.Field().MagicTraps())
	if got := p.Graveyard().Len() + inField; got != DeckSize {
		t.Errorf("graveyard + field = %d, want %d once the deck and hand are empty", got, DeckSize)
	}
	if placed == 0 || discarded == 0 || sent == 0 {
		t.Errorf("script did not exercise every move: placed=%d discarded=%d sent=%d", placed, discarded, sent)
	}
	assertInvariants(t, p)
}

func TestParticipant_TwoParticipantsAreIndependent(t *testing.T) {
	// Arrange
	p1 := newTestParticipant(t, "p1")
	p2 := newTestParticipant(t, "p2")
	before2 := stateOf(p2)

	// Act
	cards := drawN(t, p1, 2)
	if err := p1.PlaceFromHand(cards[0].ID()); err != nil {
		t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
	}
	if err := p1.DiscardFromHand(cards[1].ID()); err != nil {
		t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
	}

	// Assert
	assertUnchanged(t, p2, before2)
	if p2.Hand().Len() != 0 || p2.Graveyard().Len() != 0 || len(p2.Field().Monsters()) != 0 {
		t.Errorf("p2 zones changed by p1 transitions")
	}
	for _, c := range p2.Deck().Cards() {
		if c.Owner() != p2.ID() {
			t.Errorf("p2 instance %q owner = %q, want %q", c.ID(), c.Owner(), p2.ID())
		}
	}
	for _, c := range p1.Deck().Cards() {
		if c.Owner() != p1.ID() {
			t.Errorf("p1 instance %q owner = %q, want %q", c.ID(), c.Owner(), p1.ID())
		}
	}
	assertInvariants(t, p1)
}

func TestParticipant_CardDefinitionsAreUnchangedByOperations(t *testing.T) {
	// Arrange
	p := newMixedParticipant(t, "p1")
	type snapshot struct {
		id           CardID
		name         string
		typ          CardType
		atk, def, wt int
	}
	snap := func(c *Card) snapshot {
		return snapshot{c.ID(), c.Name(), c.Type(), c.ATK(), c.DEF(), c.Weight()}
	}
	before := make(map[*Card]snapshot)
	for _, c := range p.Deck().Cards() {
		before[c.Card()] = snap(c.Card())
	}

	// Act
	for p.Deck().Len() > 0 {
		card := fetch(t, p, p.Deck().Cards()[0].Card().Type())
		if p.Field().HasRoomFor(card.Card()) {
			if err := p.PlaceFromHand(card.ID()); err != nil {
				t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
			}
			if err := p.SendFieldCardToGraveyard(card.ID()); err != nil {
				t.Fatalf("SendFieldCardToGraveyard() returned unexpected error: %v", err)
			}
		} else if err := p.DiscardFromHand(card.ID()); err != nil {
			t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
		}
	}

	// Assert
	if len(before) == 0 {
		t.Fatalf("test setup: no card definitions recorded")
	}
	for _, c := range p.Graveyard().Cards() {
		if got, want := snap(c.Card()), before[c.Card()]; got != want {
			t.Errorf("definition of %q changed: %+v, want %+v", c.ID(), got, want)
		}
	}
}

func TestParticipant_ZoneAccessorsDoNotLetCallersChangeTheState(t *testing.T) {
	// Arrange
	p := newMixedParticipant(t, "p1")
	cards := drawN(t, p, 3)
	if err := p.PlaceFromHand(cards[0].ID()); err != nil {
		t.Fatalf("PlaceFromHand() returned unexpected error: %v", err)
	}
	if err := p.DiscardFromHand(cards[1].ID()); err != nil {
		t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
	}
	before := stateOf(p)

	// Act: mutate every slice returned by the readers.
	for _, s := range [][]*CardInstance{
		p.Deck().Cards(), p.Hand().Cards(), p.Field().Monsters(), p.Field().MagicTraps(), p.Graveyard().Cards(),
	} {
		for i := range s {
			s[i] = nil
		}
		slices.Reverse(s)
	}

	// Assert
	assertUnchanged(t, p, before)
}
