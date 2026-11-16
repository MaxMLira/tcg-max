package domain

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// recordingShuffler wraps another Shuffler and records the size of every
// Shuffle call, in order.
type recordingShuffler struct {
	inner Shuffler
	sizes []int
}

func (r *recordingShuffler) Shuffle(n int, swap func(i, j int)) {
	r.sizes = append(r.sizes, n)
	r.inner.Shuffle(n, swap)
}

// seededShuffler returns a reproducible Shuffler. The generator belongs to the
// test, never to the domain.
func seededShuffler(seed1, seed2 uint64) *rand.Rand {
	return rand.New(rand.NewPCG(seed1, seed2))
}

func setupOf(t *testing.T, id ParticipantID) ParticipantSetup {
	t.Helper()

	return ParticipantSetup{ID: id, Deck: validDeck(t), LifePoints: 4000}
}

// newTestDuel builds a Duel between "p1" and "p2" with validDeck for both.
func newTestDuel(t *testing.T, s Shuffler) *Duel {
	t.Helper()

	d, err := NewDuel(setupOf(t, "p1"), setupOf(t, "p2"), s)
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}
	return d
}

// duelState is a comparable picture of the whole Duel.
type duelState struct {
	P1, P2   state
	Turn     Turn
	Result   DuelResult
	Finished bool
}

func duelStateOf(d *Duel) duelState {
	res, _ := d.Result()
	ps := d.Participants()
	return duelState{P1: stateOf(ps[0]), P2: stateOf(ps[1]), Turn: d.Turn(), Result: res, Finished: d.Finished()}
}

func TestNewDuel_RejectsInvalidSetups(t *testing.T) {
	badDeck := DeckDefinition{entries: []DeckEntry{
		{Card: deckCard(t, 1), Copies: DeckCopiesLimit + 1},
		{Card: deckCard(t, 2), Copies: 1},
	}}

	tests := []struct {
		name    string
		a, b    ParticipantSetup
		wantErr error
	}{
		{name: "equal ids", a: setupOf(t, "p1"), b: setupOf(t, "p1"), wantErr: ErrInvalidDuel},
		{name: "blank first id", a: setupOf(t, " "), b: setupOf(t, "p2"), wantErr: ErrInvalidParticipant},
		{name: "empty second id", a: setupOf(t, "p1"), b: setupOf(t, ""), wantErr: ErrInvalidParticipant},
		{name: "padded first id", a: setupOf(t, " p1"), b: setupOf(t, "p2"), wantErr: ErrInvalidParticipant},
		{name: "padded second id", a: setupOf(t, "p1"), b: setupOf(t, "p2 "), wantErr: ErrInvalidParticipant},
		{name: "too many copies in first deck", a: ParticipantSetup{ID: "p1", Deck: badDeck, LifePoints: 4000}, b: setupOf(t, "p2"), wantErr: ErrTooManyCopies},
		{name: "too many copies in second deck", a: setupOf(t, "p1"), b: ParticipantSetup{ID: "p2", Deck: badDeck, LifePoints: 4000}, wantErr: ErrTooManyCopies},
		{name: "zero deck of first", a: ParticipantSetup{ID: "p1", LifePoints: 4000}, b: setupOf(t, "p2"), wantErr: ErrDeckSize},
		{name: "zero life points", a: ParticipantSetup{ID: "p1", Deck: validDeck(t)}, b: setupOf(t, "p2"), wantErr: ErrInvalidParticipant},
		{name: "negative life points", a: setupOf(t, "p1"), b: ParticipantSetup{ID: "p2", Deck: validDeck(t), LifePoints: -1}, wantErr: ErrInvalidParticipant},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			d, err := NewDuel(tt.a, tt.b, noopShuffler{})

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewDuel() error = %v, want %v", err, tt.wantErr)
			}
			if d != nil {
				t.Errorf("NewDuel() = %v, want nil Duel on error", d)
			}
		})
	}
}

func TestNewDuel_RejectsANilShuffler(t *testing.T) {
	// Act
	d, err := NewDuel(setupOf(t, "p1"), setupOf(t, "p2"), nil)

	// Assert
	if !errors.Is(err, ErrInvalidDeck) {
		t.Fatalf("NewDuel() error = %v, want ErrInvalidDeck", err)
	}
	if d != nil {
		t.Errorf("NewDuel() = %v, want nil Duel on error", d)
	}
}

func TestNewDuel_EqualIdsAreRejectedBeforeShuffling(t *testing.T) {
	// Arrange
	rec := &recordingShuffler{inner: noopShuffler{}}

	// Act
	d, err := NewDuel(setupOf(t, "same"), setupOf(t, "same"), rec)

	// Assert
	if !errors.Is(err, ErrInvalidDuel) {
		t.Fatalf("NewDuel() error = %v, want ErrInvalidDuel", err)
	}
	if !strings.Contains(err.Error(), `"same"`) {
		t.Errorf("error %q should mention the id", err)
	}
	if d != nil {
		t.Errorf("NewDuel() = %v, want nil Duel", d)
	}
	if len(rec.sizes) != 0 {
		t.Errorf("Shuffle called %d times (%v), want 0", len(rec.sizes), rec.sizes)
	}
}

func TestNewDuel_FailureIdentifiesTheFailingParticipant(t *testing.T) {
	// Arrange
	bad := ParticipantSetup{ID: "p2", Deck: DeckDefinition{}, LifePoints: 4000}

	// Act
	d, err := NewDuel(setupOf(t, "p1"), bad, noopShuffler{})

	// Assert
	if !errors.Is(err, ErrDeckSize) || d != nil {
		t.Fatalf("NewDuel() = %v, %v, want nil Duel and ErrDeckSize", d, err)
	}
	if !strings.Contains(err.Error(), `"p2"`) {
		t.Errorf("error %q should identify the failing participant", err)
	}
}

func TestNewDuel_IdsDifferingOnlyByCaseAreDifferentParticipants(t *testing.T) {
	// Act
	d, err := NewDuel(setupOf(t, "p1"), setupOf(t, "P1"), noopShuffler{})

	// Assert
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}
	if _, ok := d.Participant("P1"); !ok {
		t.Errorf("Participant(P1) not found")
	}
}

func TestNewDuel_CallsTheShufflerForTwoDecksAndTheFirstPlayer(t *testing.T) {
	// Arrange
	rec := &recordingShuffler{inner: noopShuffler{}}

	// Act
	newTestDuel(t, rec)

	// Assert
	if want := []int{DeckSize, DeckSize, 2}; !slices.Equal(rec.sizes, want) {
		t.Errorf("Shuffle sizes = %v, want %v", rec.sizes, want)
	}
}

func TestNewDuel_DealsExactlyFourCardsToEachParticipant(t *testing.T) {
	// Arrange
	d := newTestDuel(t, &fakeShuffler{})

	// Act
	ps := d.Participants()

	// Assert
	seen := make(map[InstanceID]ParticipantID)
	for _, p := range ps {
		if p.Hand().Len() != HandLimit {
			t.Errorf("%q hand = %d cards, want %d (the first player gets no extra card)", p.ID(), p.Hand().Len(), HandLimit)
		}
		if p.Deck().Len() != DeckSize-HandLimit {
			t.Errorf("%q deck = %d cards, want %d", p.ID(), p.Deck().Len(), DeckSize-HandLimit)
		}
		if p.LifePoints() != 4000 {
			t.Errorf("%q life points = %d, want 4000", p.ID(), p.LifePoints())
		}
		for _, c := range p.Hand().Cards() {
			if c.Location() != LocationHand {
				t.Errorf("%q location = %s, want %s", c.ID(), c.Location(), LocationHand)
			}
		}
		for _, c := range slices.Concat(p.Hand().Cards(), p.Deck().Cards()) {
			if c.Owner() != p.ID() {
				t.Errorf("%q owner = %q, want %q", c.ID(), c.Owner(), p.ID())
			}
			if !strings.HasPrefix(string(c.ID()), string(p.ID())+"#") {
				t.Errorf("%q does not start with %q#", c.ID(), p.ID())
			}
			if other, dup := seen[c.ID()]; dup {
				t.Errorf("instance id %q repeated (also held by %q)", c.ID(), other)
			}
			seen[c.ID()] = p.ID()
		}
		if len(p.Field().Monsters())+len(p.Field().MagicTraps())+p.Graveyard().Len() != 0 {
			t.Errorf("%q field and graveyard must start empty", p.ID())
		}
		assertInvariants(t, p)
	}
	if len(seen) != 2*DeckSize {
		t.Errorf("distinct instances = %d, want %d", len(seen), 2*DeckSize)
	}
}

func TestNewDuel_OpeningHandIsTheTopOfTheShuffledDeck(t *testing.T) {
	// Arrange: the reversing fake puts the last expansion ids on top.
	d := newTestDuel(t, &fakeShuffler{})
	p1, _ := d.Participant("p1")

	// Act
	got := ids(p1.Hand().Cards())

	// Assert
	want := []InstanceID{"p1#30", "p1#29", "p1#28", "p1#27"}
	if !slices.Equal(got, want) {
		t.Errorf("opening hand = %v, want %v", got, want)
	}
}

func TestNewDuel_ParticipantIdsContainingHashStillGiveUniqueInstanceIds(t *testing.T) {
	// Arrange
	d, err := NewDuel(setupOf(t, "a#1"), setupOf(t, "a"), noopShuffler{})
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}

	// Act
	seen := make(map[InstanceID]bool)
	for _, p := range d.Participants() {
		for _, c := range slices.Concat(p.Hand().Cards(), p.Deck().Cards()) {
			// Assert
			if seen[c.ID()] {
				t.Errorf("instance id %q repeated", c.ID())
			}
			seen[c.ID()] = true
			if !strings.HasPrefix(string(c.ID()), string(p.ID())+"#") {
				t.Errorf("instance id %q does not split as %q#<n>", c.ID(), p.ID())
			}
			if c.Owner() != p.ID() {
				t.Errorf("%q owner = %q, want %q", c.ID(), c.Owner(), p.ID())
			}
		}
	}
	if len(seen) != 2*DeckSize {
		t.Errorf("distinct instances = %d, want %d", len(seen), 2*DeckSize)
	}
}

func TestNewDuel_ParticipantsAreIndependentEvenWithEqualDeckDefinitions(t *testing.T) {
	// Arrange
	def := validDeck(t)
	d, err := NewDuel(
		ParticipantSetup{ID: "p1", Deck: def, LifePoints: 4000},
		ParticipantSetup{ID: "p2", Deck: def, LifePoints: 4000},
		noopShuffler{},
	)
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	p2, _ := d.Participant("p2")
	before2 := stateOf(p2)

	// Act
	if err := p1.DiscardFromHand(p1.Hand().Cards()[0].ID()); err != nil {
		t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
	}
	if _, err := p1.Draw(); err != nil {
		t.Fatalf("Draw() returned unexpected error: %v", err)
	}

	// Assert
	assertUnchanged(t, p2, before2)
	if p1.Deck() == p2.Deck() || p1.Hand() == p2.Hand() || p1.Graveyard() == p2.Graveyard() || p1.Field() == p2.Field() {
		t.Errorf("participants share a zone")
	}
	for _, c1 := range p1.Hand().Cards() {
		for _, c2 := range p2.Hand().Cards() {
			if c1 == c2 || c1.ID() == c2.ID() {
				t.Errorf("the hands share the card %q", c1.ID())
			}
		}
	}
}

func TestNewDuel_StartsAtTurnOneWithTheChosenFirstPlayer(t *testing.T) {
	tests := []struct {
		name      string
		shuffler  Shuffler
		wantFirst ParticipantID
	}{
		{name: "noop shuffler keeps the first setup first", shuffler: noopShuffler{}, wantFirst: "p1"},
		{name: "reversing shuffler selects the second setup", shuffler: &fakeShuffler{}, wantFirst: "p2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			d := newTestDuel(t, tt.shuffler)

			// Assert
			turn := d.Turn()
			if turn.Number() != FirstTurnNumber || turn.Participant() != tt.wantFirst || !turn.IsFirst() {
				t.Errorf("Turn() = {%d, %q, first=%v}, want {%d, %q, first=true}",
					turn.Number(), turn.Participant(), turn.IsFirst(), FirstTurnNumber, tt.wantFirst)
			}
			if _, ok := d.Result(); ok {
				t.Errorf("Result() reports a result for a new Duel")
			}
			if d.Finished() {
				t.Errorf("Finished() = true for a new Duel")
			}
			if first, ok := d.Participant(turn.Participant()); !ok || first.Hand().Len() != HandLimit {
				t.Errorf("the first player must hold exactly %d cards", HandLimit)
			}
		})
	}
}

func TestDuel_ParticipantLookup(t *testing.T) {
	// Arrange
	d := newTestDuel(t, noopShuffler{})

	// Act
	p1, ok1 := d.Participant("p1")
	p2, ok2 := d.Participant("p2")
	none, okNone := d.Participant("p3")

	// Assert
	if !ok1 || p1 != d.Participants()[0] {
		t.Errorf("Participant(p1) = %v, %v, want the first participant", p1, ok1)
	}
	if !ok2 || p2 != d.Participants()[1] {
		t.Errorf("Participant(p2) = %v, %v, want the second participant", p2, ok2)
	}
	if okNone || none != nil {
		t.Errorf("Participant(p3) = %v, %v, want nil, false", none, okNone)
	}
}

func TestDuel_ReplenishHandFillsTheHandOfOneParticipant(t *testing.T) {
	tests := []struct {
		name      string
		discard   int
		wantDrawn int
	}{
		{name: "three cards draws one", discard: 1, wantDrawn: 1},
		{name: "one card draws three", discard: 3, wantDrawn: 3},
		{name: "full hand is a no-op", discard: 0, wantDrawn: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := newTestDuel(t, noopShuffler{})
			p1, _ := d.Participant("p1")
			p2, _ := d.Participant("p2")
			for _, c := range p1.Hand().Cards()[:tt.discard] {
				if err := p1.DiscardFromHand(c.ID()); err != nil {
					t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
				}
			}
			deckBefore := p1.Deck().Len()
			before2 := stateOf(p2)

			// Act
			err := d.replenishHand("p1")

			// Assert
			if err != nil {
				t.Fatalf("replenishHand() returned unexpected error: %v", err)
			}
			if p1.Hand().Len() != HandLimit {
				t.Errorf("Hand().Len() = %d, want %d", p1.Hand().Len(), HandLimit)
			}
			if got := p1.Deck().Len(); got != deckBefore-tt.wantDrawn {
				t.Errorf("Deck().Len() = %d, want %d", got, deckBefore-tt.wantDrawn)
			}
			assertUnchanged(t, p2, before2)
			if d.Finished() {
				t.Errorf("Finished() = true after a successful replenishment")
			}

			// A repeated call must not draw again.
			afterFirst := duelStateOf(d)
			if err := d.replenishHand("p1"); err != nil {
				t.Fatalf("second replenishHand() returned unexpected error: %v", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, afterFirst) {
				t.Errorf("a repeated replenishHand changed the state")
			}
		})
	}
}

func TestDuel_ReplenishHandRejectsAnUnknownParticipant(t *testing.T) {
	// Arrange
	d := newTestDuel(t, noopShuffler{})
	before := duelStateOf(d)

	// Act
	err := d.replenishHand("ghost")

	// Assert
	if !errors.Is(err, ErrInvalidParticipant) {
		t.Fatalf("replenishHand() error = %v, want ErrInvalidParticipant", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("state changed by a rejected replenishment")
	}
}

func TestDuel_ReplenishHandWithAShortDeckEndsTheDuelByDeckOut(t *testing.T) {
	tests := []struct {
		name       string
		loser      ParticipantID
		wantWinner ParticipantID
	}{
		{name: "first setup loses", loser: "p1", wantWinner: "p2"},
		{name: "second setup loses", loser: "p2", wantWinner: "p1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 2 cards left in the deck and an empty hand.
			d := newTestDuel(t, noopShuffler{})
			p, _ := d.Participant(tt.loser)
			for _, c := range p.Hand().Cards() {
				if err := p.DiscardFromHand(c.ID()); err != nil {
					t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
				}
			}
			discardUntilDeckHas(t, p, 2)
			beforeLoser := stateOf(p)

			// Act
			err := d.replenishHand(tt.loser)

			// Assert
			if err != nil {
				t.Fatalf("replenishHand() returned unexpected error: %v (a loss is an outcome, not a rejection)", err)
			}
			res, ok := d.Result()
			if !ok || !d.Finished() {
				t.Fatalf("Result() ok = %v, Finished() = %v, want both true", ok, d.Finished())
			}
			if res.Loser() != tt.loser || res.Winner() != tt.wantWinner || res.Reason() != EndReasonDeckOut {
				t.Errorf("Result = {winner %q, loser %q, reason %v}, want {%q, %q, EndReasonDeckOut}",
					res.Winner(), res.Loser(), res.Reason(), tt.wantWinner, tt.loser)
			}
			assertUnchanged(t, p, beforeLoser)

			// A finished Duel is frozen.
			frozen := duelStateOf(d)
			for _, id := range []ParticipantID{"p1", "p2"} {
				if err := d.replenishHand(id); !errors.Is(err, ErrDuelEnded) {
					t.Errorf("replenishHand(%q) after the end: error = %v, want ErrDuelEnded", id, err)
				}
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, frozen) {
				t.Errorf("state changed after the Duel ended")
			}
		})
	}
}

func TestNewDuel_IsDeterministicForTheSameSeed(t *testing.T) {
	// Arrange
	build := func() *Duel { return newTestDuel(t, seededShuffler(1, 2)) }

	// Act
	first, second := build(), build()

	// Assert
	if first.Turn() != second.Turn() {
		t.Errorf("first player differs: %v vs %v", first.Turn(), second.Turn())
	}
	cardOf := func(d *Duel) map[InstanceID]CardID {
		m := make(map[InstanceID]CardID)
		for _, p := range d.Participants() {
			for _, c := range slices.Concat(p.Hand().Cards(), p.Deck().Cards()) {
				m[c.ID()] = c.Card().ID()
			}
		}
		return m
	}
	if !reflect.DeepEqual(cardOf(first), cardOf(second)) {
		t.Errorf("InstanceID to CardID mapping differs for the same seed")
	}
	for i := range 2 {
		a, b := first.Participants()[i], second.Participants()[i]
		if got, want := ids(a.Deck().Cards()), ids(b.Deck().Cards()); !slices.Equal(got, want) {
			t.Errorf("deck order of %q differs:\n%v\n%v", a.ID(), got, want)
		}
		if got, want := ids(a.Hand().Cards()), ids(b.Hand().Cards()); !slices.Equal(got, want) {
			t.Errorf("hand of %q differs: %v vs %v", a.ID(), got, want)
		}
	}
}

// prepareHandAndDeck leaves p with exactly inHand cards in the Hand and
// deckLeft in the Deck, using only the exported transitions (the other cards
// go to the Graveyard).
func prepareHandAndDeck(t *testing.T, p *Participant, inHand, deckLeft int) {
	t.Helper()

	for _, c := range slices.Clone(p.Hand().Cards()) {
		if err := p.DiscardFromHand(c.ID()); err != nil {
			t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
		}
	}
	discardUntilDeckHas(t, p, deckLeft+inHand)
	drawN(t, p, inHand)
	if p.Hand().Len() != inHand || p.Deck().Len() != deckLeft {
		t.Fatalf("setup: hand, deck = %d, %d, want %d, %d", p.Hand().Len(), p.Deck().Len(), inHand, deckLeft)
	}
}

func TestDuel_ReplenishHandWithExactlyTheMissingCardsDoesNotLose(t *testing.T) {
	tests := []struct {
		name     string
		inHand   int
		deckLeft int
	}{
		{name: "empty hand and exactly four cards", inHand: 0, deckLeft: 4},
		{name: "three cards and exactly one card", inHand: 3, deckLeft: 1},
		{name: "two cards and exactly two cards", inHand: 2, deckLeft: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := newTestDuel(t, noopShuffler{})
			p, _ := d.Participant("p1")
			other, _ := d.Participant("p2")
			prepareHandAndDeck(t, p, tt.inHand, tt.deckLeft)
			beforeOther := stateOf(other)

			// Act
			err := d.replenishHand("p1")

			// Assert
			if err != nil {
				t.Fatalf("replenishHand() returned unexpected error: %v", err)
			}
			if p.Hand().Len() != HandLimit || p.Deck().Len() != 0 {
				t.Errorf("hand, deck = %d, %d, want %d, 0", p.Hand().Len(), p.Deck().Len(), HandLimit)
			}
			if res, ok := d.Result(); ok || res != (DuelResult{}) || d.Finished() {
				t.Errorf("Result() = %+v, %v, Finished() = %v, want no result: drawing the last card is not a loss", res, ok, d.Finished())
			}
			assertUnchanged(t, other, beforeOther)
		})
	}
}

func TestDuel_ReplenishHandOneCardShortEndsTheDuelAndLeavesTheZonesAlone(t *testing.T) {
	tests := []struct {
		name       string
		loser      ParticipantID
		wantWinner ParticipantID
		inHand     int
		deckLeft   int
	}{
		{name: "empty hand and three cards", loser: "p1", wantWinner: "p2", inHand: 0, deckLeft: 3},
		{name: "empty hand and three cards, second setup", loser: "p2", wantWinner: "p1", inHand: 0, deckLeft: 3},
		{name: "three cards and an empty deck", loser: "p1", wantWinner: "p2", inHand: 3, deckLeft: 0},
		{name: "two cards and one card", loser: "p1", wantWinner: "p2", inHand: 2, deckLeft: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := newTestDuel(t, noopShuffler{})
			p, _ := d.Participant(tt.loser)
			other, _ := d.Participant(tt.wantWinner)
			prepareHandAndDeck(t, p, tt.inHand, tt.deckLeft)
			beforeLoser, beforeOther := stateOf(p), stateOf(other)

			// Act
			err := d.replenishHand(tt.loser)

			// Assert
			if err != nil {
				t.Fatalf("replenishHand() returned unexpected error: %v (a loss is an outcome, not a rejection)", err)
			}
			res, ok := d.Result()
			if !ok || !d.Finished() {
				t.Fatalf("Result() ok = %v, Finished() = %v, want both true", ok, d.Finished())
			}
			if res.Loser() != tt.loser || res.Winner() != tt.wantWinner || res.Reason() != EndReasonDeckOut {
				t.Errorf("Result = {winner %q, loser %q, reason %v}, want {%q, %q, EndReasonDeckOut}",
					res.Winner(), res.Loser(), res.Reason(), tt.wantWinner, tt.loser)
			}
			assertUnchanged(t, p, beforeLoser)
			assertUnchanged(t, other, beforeOther)

			// The Duel is frozen afterwards.
			frozen := duelStateOf(d)
			for _, id := range []ParticipantID{"p1", "p2"} {
				if err := d.replenishHand(id); !errors.Is(err, ErrDuelEnded) {
					t.Errorf("replenishHand(%q) after the end: error = %v, want ErrDuelEnded", id, err)
				}
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, frozen) {
				t.Errorf("state changed after the Duel ended")
			}
		})
	}
}

// permutingShuffler applies a different permutation on every call, in order:
// reverse reverses the n elements and identity leaves them alone. It records
// the size of each call.
type permutingShuffler struct {
	reverse []bool
	sizes   []int
}

func (s *permutingShuffler) Shuffle(n int, swap func(i, j int)) {
	call := len(s.sizes)
	s.sizes = append(s.sizes, n)
	if call >= len(s.reverse) || !s.reverse[call] {
		return
	}
	for i := 0; i < n/2; i++ {
		swap(i, n-1-i)
	}
}

func TestNewDuel_BuildsTheFirstSetupBeforeTheSecondAndFlipsTheCoinLast(t *testing.T) {
	// Arrange: call 1 (deck of a) reverses, call 2 (deck of b) keeps the
	// expansion order, call 3 (coin flip) reverses.
	s := &permutingShuffler{reverse: []bool{true, false, true}}

	// Act
	d := newTestDuel(t, s)

	// Assert
	if want := []int{DeckSize, DeckSize, 2}; !slices.Equal(s.sizes, want) {
		t.Fatalf("Shuffle sizes = %v, want %v", s.sizes, want)
	}

	// Top of the Deck first: the hand holds the first four cards drawn.
	wholeDeck := func(id ParticipantID) []InstanceID {
		p, _ := d.Participant(id)
		return slices.Concat(ids(p.Hand().Cards()), ids(p.Deck().Cards()))
	}
	var reversed, original []InstanceID
	for n := DeckSize; n >= 1; n-- {
		reversed = append(reversed, InstanceID(fmt.Sprintf("p1#%d", n)))
	}
	for n := 1; n <= DeckSize; n++ {
		original = append(original, InstanceID(fmt.Sprintf("p2#%d", n)))
	}
	if got := wholeDeck("p1"); !slices.Equal(got, reversed) {
		t.Errorf("deck order of p1 = %v, want the reversed expansion order %v (first shuffle call)", got, reversed)
	}
	if got := wholeDeck("p2"); !slices.Equal(got, original) {
		t.Errorf("deck order of p2 = %v, want the expansion order %v (second shuffle call)", got, original)
	}
	if got := d.Turn().Participant(); got != "p2" {
		t.Errorf("first player = %q, want %q (the third call reverses the two candidates)", got, "p2")
	}
}
