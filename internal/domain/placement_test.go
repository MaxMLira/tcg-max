package domain

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

// weightMonster builds a Monster definition with the given weight.
func weightMonster(t *testing.T, n, weight int) *Card {
	t.Helper()

	return mustCard(t, cardID(700+n), fmt.Sprintf("Weight %d monster %d", weight, n), CardTypeMonster, 1000, 1000, weight, SubtypeMage())
}

// padCard builds a filler Monster definition of weight 1.
func padCard(t *testing.T, n int) *Card {
	t.Helper()

	return mustCard(t, cardID(900+n), fmt.Sprintf("Pad %d", n), CardTypeMonster, 1000, 1000, 1, SubtypeAquatic())
}

// paddedDeck builds a 30-card deck: the head cards first, one copy each,
// padded with padCard up to DeckSize.
func paddedDeck(t *testing.T, head ...*Card) DeckDefinition {
	t.Helper()

	entries := make([]DeckEntry, 0, DeckSize)
	for _, c := range head {
		entries = append(entries, DeckEntry{Card: c, Copies: 1})
	}
	for n := 1; len(entries) < DeckSize; n++ {
		entries = append(entries, DeckEntry{Card: padCard(t, n), Copies: 1})
	}
	deck, err := NewDeckDefinition(entries)
	if err != nil {
		t.Fatalf("paddedDeck is not valid: %v", err)
	}
	return deck
}

// duelWithP1Deck builds a Duel where p1 uses paddedDeck(t, head...) and p2
// uses validDeck, with a noop shuffler. p1's opening hand is therefore the
// first four cards of head in order (p1#1..p1#4), p2 is never the first
// player, and p1 plays turn 1.
func duelWithP1Deck(t *testing.T, head ...*Card) *Duel {
	t.Helper()

	a := ParticipantSetup{ID: "p1", Deck: paddedDeck(t, head...), LifePoints: 4000}
	b := ParticipantSetup{ID: "p2", Deck: validDeck(t), LifePoints: 4000}
	d, err := NewDuel(a, b, noopShuffler{})
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}
	return d
}

func TestPlace_SummonsASingleCardWithoutTributes(t *testing.T) {
	// Arrange — p1's hand: p1#1 (weight 2) and three pads (weight 1).
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))

	// Act
	err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	monsters := p1.Field().Monsters()
	if len(monsters) != 1 || monsters[0].ID() != "p1#1" {
		t.Fatalf("Field = %v, want p1#1 alone", ids(monsters))
	}
	if monsters[0].Position() != PositionAttack || monsters[0].Face() != FaceUp {
		t.Errorf("placed monster = (%v, %v), want (Attack, Face-up)", monsters[0].Position(), monsters[0].Face())
	}
	if got := p1.Graveyard().Cards(); len(got) != 0 {
		t.Errorf("Graveyard = %v, want empty", ids(got))
	}
	if got := p1.Hand().Len(); got != HandLimit {
		t.Errorf("hand holds %d cards after replenishment, want %d", got, HandLimit)
	}
	assertInvariants(t, p1)
}

func TestPlace_SummonsWithThreeTributesKeepsTheGraveyardOrder(t *testing.T) {
	// Arrange — p1's hand: weights 2, 7, 10, 14.
	d := duelWithP1Deck(t, weightMonster(t, 1, 2), weightMonster(t, 2, 7), weightMonster(t, 3, 10), weightMonster(t, 4, 14))

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2", "p1#3", "p1#4"}, PositionDefense, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	monsters := p1.Field().Monsters()
	if len(monsters) != 1 || monsters[0].ID() != "p1#4" {
		t.Fatalf("Field = %v, want p1#4 (weight 14) alone", ids(monsters))
	}
	if monsters[0].Position() != PositionDefense || monsters[0].Face() != FaceUp {
		t.Errorf("placed monster = (%v, %v), want (Defense, Face-up)", monsters[0].Position(), monsters[0].Face())
	}
	if got, want := ids(p1.Graveyard().Cards()), []InstanceID{"p1#1", "p1#2", "p1#3"}; !slices.Equal(got, want) {
		t.Errorf("Graveyard = %v, want %v (selection order)", got, want)
	}
	if got := p1.Hand().Len(); got != HandLimit {
		t.Errorf("hand holds %d cards after replenishment, want %d", got, HandLimit)
	}
	assertInvariants(t, p1)
}

func TestPlace_RejectsAnEmptyOrOversizedSelection(t *testing.T) {
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	before := duelStateOf(d)

	tests := []struct {
		name string
		sel  []InstanceID
	}{
		{name: "empty", sel: nil},
		{name: "five cards", sel: []InstanceID{"p1#1", "p1#2", "p1#3", "p1#4", "p1#5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := d.Place(tt.sel, PositionAttack, FaceUp)

			// Assert
			if !errors.Is(err, ErrInvalidSelection) {
				t.Fatalf("Place() error = %v, want ErrInvalidSelection", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
				t.Errorf("Duel state changed by a rejected selection")
			}
		})
	}
}
func TestPlace_RejectsADuplicateSelection(t *testing.T) {
	// Arrange
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2", "p1#1"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("Place() error = %v, want ErrInvalidSelection", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected selection")
	}
}

func TestPlace_RejectsACardNotInTheActiveHand(t *testing.T) {
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))

	tests := []struct {
		name string
		id   InstanceID
	}{
		{name: "a card in p1's deck", id: "p1#5"},
		{name: "a card in p2's hand", id: "p2#1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			before := duelStateOf(d)

			// Act
			err := d.Place([]InstanceID{tt.id}, PositionAttack, FaceUp)

			// Assert
			if !errors.Is(err, ErrCardNotFound) {
				t.Fatalf("Place() error = %v, want ErrCardNotFound", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
				t.Errorf("Duel state changed by a rejected selection")
			}
		})
	}
}

func TestRequiredTributes_TiersByWeight(t *testing.T) {
	tests := []struct {
		weight int
		want   int
	}{
		{weight: 1, want: 0},
		{weight: 5, want: 0},
		{weight: 6, want: 1},
		{weight: 8, want: 1},
		{weight: 9, want: 2},
		{weight: 12, want: 2},
		{weight: 13, want: 3},
		{weight: 30, want: 3},
	}
	for _, tt := range tests {
		if got := requiredTributes(tt.weight); got != tt.want {
			t.Errorf("requiredTributes(%d) = %d, want %d (game-rules §16)", tt.weight, got, tt.want)
		}
	}
}

func TestPlace_TributeTiersByWeight(t *testing.T) {
	// p1's hand: p1#1 (weight 2), p1#2 (weight 7), p1#3 (weight 10),
	// p1#4 (weight 14).
	tests := []struct {
		name      string
		selection []InstanceID
		wantErr   error
	}{
		{name: "weight 2 needs zero, none given", selection: []InstanceID{"p1#1"}},
		{name: "weight 7 needs one, one given", selection: []InstanceID{"p1#1", "p1#2"}},
		{name: "weight 7 needs one, none given", selection: []InstanceID{"p1#2"}, wantErr: ErrInvalidPlacement},
		{name: "weight 10 needs two, two given", selection: []InstanceID{"p1#1", "p1#2", "p1#3"}},
		{name: "weight 10 needs two, one given", selection: []InstanceID{"p1#2", "p1#3"}, wantErr: ErrInvalidPlacement},
		{name: "weight 14 needs three, three given", selection: []InstanceID{"p1#1", "p1#2", "p1#3", "p1#4"}},
		{name: "weight 14 needs three, two given", selection: []InstanceID{"p1#2", "p1#3", "p1#4"}, wantErr: ErrInvalidPlacement},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := duelWithP1Deck(t, weightMonster(t, 1, 2), weightMonster(t, 2, 7), weightMonster(t, 3, 10), weightMonster(t, 4, 14))
			before := duelStateOf(d)

			// Act
			err := d.Place(tt.selection, PositionAttack, FaceUp)

			// Assert
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Place() error = %v, want %v", err, tt.wantErr)
				}
				if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
					t.Errorf("Duel state changed by a rejected summon")
				}
				return
			}
			if err != nil {
				t.Fatalf("Place() returned unexpected error: %v", err)
			}
			p1, _ := d.Participant("p1")
			last := tt.selection[len(tt.selection)-1]
			if _, ok := p1.Field().Get(last); !ok {
				t.Fatalf("the placed card %q is not on the Field", last)
			}
		})
	}
}

func TestPlace_FaceDownOnlyWithoutTributes(t *testing.T) {
	// Arrange — face-down defense with zero tributes is legal (game-rules §19).
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))

	// Act
	err := d.Place([]InstanceID{"p1#1"}, PositionDefense, FaceDown)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	c, ok := p1.Field().Get("p1#1")
	if !ok {
		t.Fatal("the placed card is not on the Field")
	}
	if c.Position() != PositionDefense || c.Face() != FaceDown {
		t.Errorf("placed monster = (%v, %v), want (Defense, Face-down)", c.Position(), c.Face())
	}
}

func TestPlace_RejectsFaceDownWithTributes(t *testing.T) {
	// Arrange
	d := duelWithP1Deck(t, weightMonster(t, 1, 2), weightMonster(t, 2, 7))
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2"}, PositionDefense, FaceDown)

	// Assert
	if !errors.Is(err, ErrInvalidPlacement) {
		t.Fatalf("Place() error = %v, want ErrInvalidPlacement (game-rules §19)", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected summon")
	}
}

func TestPlace_RejectsAnIncoherentPositionFace(t *testing.T) {
	// Arrange
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceDown)

	// Assert
	if !errors.Is(err, ErrInvalidPlacement) {
		t.Fatalf("Place() error = %v, want ErrInvalidPlacement", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected summon")
	}
}

func TestPlace_RejectsAMagicOrTrapAsThePlacedCard(t *testing.T) {
	tests := []struct {
		name string
		card *Card
	}{
		{name: "magic", card: magicCard(t, 1)},
		{name: "trap", card: trapCard(t, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := duelWithP1Deck(t, tt.card)
			before := duelStateOf(d)

			// Act
			err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp)

			// Assert
			if !errors.Is(err, ErrUnsupportedCardType) {
				t.Fatalf("Place() error = %v, want ErrUnsupportedCardType", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
				t.Errorf("Duel state changed by a rejected placement")
			}
		})
	}
}

func TestPlace_RejectsANonMonsterAmongTheTributes(t *testing.T) {
	// Arrange — p1#1 is a Magic, p1#2 a weight-7 Monster.
	d := duelWithP1Deck(t, magicCard(t, 1), weightMonster(t, 2, 7))
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("Place() error = %v, want ErrInvalidSelection (tributes are Monsters, game-rules §17)", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected selection")
	}
}

func TestPlace_FullMonsterAreaBlocksTheSummonEvenWithTributes(t *testing.T) {
	// Arrange — fill p1's monster area over four rounds. In the MVP tributes
	// come from the hand (mvp-story-mode §1.5), so a full area always blocks a
	// new summon (mvp-story-mode §1.2).
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	for round := 1; round <= MonsterSlots; round++ {
		p1, _ := d.Participant("p1")
		p2, _ := d.Participant("p2")
		if err := d.Place([]InstanceID{p1.Hand().Cards()[0].ID()}, PositionAttack, FaceUp); err != nil {
			t.Fatalf("round %d: p1 Place() = %v", round, err)
		}
		endFullTurn(t, d)
		if err := d.Place([]InstanceID{p2.Hand().Cards()[0].ID()}, PositionAttack, FaceUp); err != nil {
			t.Fatalf("round %d: p2 Place() = %v", round, err)
		}
		endFullTurn(t, d)
	}
	p1, _ := d.Participant("p1")
	if got := len(p1.Field().Monsters()); got != MonsterSlots {
		t.Fatalf("p1's monster area holds %d monsters, want %d", got, MonsterSlots)
	}

	// Act — the fifth summon with a full selection.
	sel := ids(p1.Hand().Cards())
	before := duelStateOf(d)
	err := d.Place(sel, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrFieldFull) {
		t.Fatalf("Place() error = %v, want ErrFieldFull", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a blocked summon")
	}
}

func TestPlace_OnlyOncePerPhase1(t *testing.T) {
	// Arrange
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	if err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp); err != nil {
		t.Fatalf("first Place() returned unexpected error: %v", err)
	}
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrPlacementUsed) {
		t.Fatalf("second Place() error = %v, want ErrPlacementUsed (game-rules §14.1)", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by the second placement")
	}
}

func TestPlace_RejectedInPhase2(t *testing.T) {
	// Arrange
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrWrongPhase) {
		t.Fatalf("Place() error = %v, want ErrWrongPhase", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected placement")
	}
}

func TestPlace_ReplenishesTheHandWithTheNextDeckCards(t *testing.T) {
	// Arrange — a full selection empties the hand.
	d := duelWithP1Deck(t, weightMonster(t, 1, 2), weightMonster(t, 2, 7), weightMonster(t, 3, 10), weightMonster(t, 4, 14))

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2", "p1#3", "p1#4"}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got, want := ids(p1.Hand().Cards()), []InstanceID{"p1#5", "p1#6", "p1#7", "p1#8"}; !slices.Equal(got, want) {
		t.Errorf("hand = %v after replenishment, want %v", got, want)
	}
	assertInvariants(t, p1)
}

// padsThenWeights builds a 30-card deck: pads first, then the given weight
// Monsters last, in order, one copy each.
func padsThenWeights(t *testing.T, weights ...int) DeckDefinition {
	t.Helper()

	entries := make([]DeckEntry, 0, DeckSize)
	for n := 1; len(entries)+len(weights) < DeckSize; n++ {
		entries = append(entries, DeckEntry{Card: padCard(t, n), Copies: 1})
	}
	for i, w := range weights {
		entries = append(entries, DeckEntry{Card: weightMonster(t, i+1, w), Copies: 1})
	}
	deck, err := NewDeckDefinition(entries)
	if err != nil {
		t.Fatalf("padsThenWeights is not valid: %v", err)
	}
	return deck
}

func TestPlace_DeckOutLosesImmediatelyAndFreezesTheDuel(t *testing.T) {
	// Arrange — p1's hand holds weights 2, 7, 10, 14 and the deck holds none.
	// The summon of the weight-14 Monster needs exactly three tributes.
	a := ParticipantSetup{ID: "p1", Deck: padsThenWeights(t, 2, 7, 10, 14), LifePoints: 4000}
	b := ParticipantSetup{ID: "p2", Deck: validDeck(t), LifePoints: 4000}
	d, err := NewDuel(a, b, noopShuffler{})
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	for _, c := range slices.Clone(p1.Hand().Cards()) {
		if err := p1.DiscardFromHand(c.ID()); err != nil {
			t.Fatalf("DiscardFromHand() = %v", err)
		}
	}
	discardUntilDeckHas(t, p1, 4)
	drawN(t, p1, 4)
	sel := ids(p1.Hand().Cards())
	if p1.Deck().Len() != 0 {
		t.Fatalf("setup: deck holds %d cards, want 0", p1.Deck().Len())
	}

	// Act — the full selection empties the hand; the deck cannot refill it.
	err = d.Place(sel, PositionDefense, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() = %v, want the action to succeed and the loss to be the outcome", err)
	}
	if !d.Finished() {
		t.Fatal("the Duel is not finished, want a deck-out loss (game-rules §8)")
	}
	res, _ := d.Result()
	if res.Winner() != "p2" || res.Loser() != "p1" || res.Reason() != EndReasonDeckOut {
		t.Errorf("result = (winner %q, loser %q, reason %v), want (p2, p1, deck out)", res.Winner(), res.Loser(), res.Reason())
	}
	if got, want := ids(p1.Field().Monsters()), []InstanceID{sel[len(sel)-1]}; !slices.Equal(got, want) {
		t.Errorf("Field = %v, want the placed card %v", got, want)
	}
	// the last three graveyard entries are the tributes, in selection order
	got := ids(p1.Graveyard().Cards())
	if !slices.Equal(got[len(got)-3:], sel[:len(sel)-1]) {
		t.Errorf("Graveyard = %v, want the tributes %v appended in order", got, sel[:len(sel)-1])
	}
	if got := d.Turn(); !got.IsFirst() {
		t.Errorf("turn after the loss = %v, want still turn 1", got)
	}
	// the Duel is frozen
	if err := d.EndPhase1(); !errors.Is(err, ErrDuelEnded) {
		t.Errorf("EndPhase1() error = %v, want ErrDuelEnded", err)
	}
	if err := d.Place(sel, PositionAttack, FaceUp); !errors.Is(err, ErrDuelEnded) {
		t.Errorf("Place() error = %v, want ErrDuelEnded", err)
	}
}
func TestPlace_AllowsTheFourthMonsterWhenTheAreaHasRoom(t *testing.T) {
	// Arrange — p1's area holds three monsters after three rounds.
	d := duelWithP1Deck(t, weightMonster(t, 1, 2))
	for round := 1; round < MonsterSlots; round++ {
		p1, _ := d.Participant("p1")
		if err := d.Place([]InstanceID{p1.Hand().Cards()[0].ID()}, PositionAttack, FaceUp); err != nil {
			t.Fatalf("round %d: p1 Place() = %v", round, err)
		}
		endFullTurn(t, d)
		p2, _ := d.Participant("p2")
		if err := d.Place([]InstanceID{p2.Hand().Cards()[0].ID()}, PositionAttack, FaceUp); err != nil {
			t.Fatalf("round %d: p2 Place() = %v", round, err)
		}
		endFullTurn(t, d)
	}
	p1, _ := d.Participant("p1")
	if got := len(p1.Field().Monsters()); got != MonsterSlots-1 {
		t.Fatalf("p1's monster area holds %d monsters, want %d", got, MonsterSlots-1)
	}

	// Act — the fourth monster fits.
	err := d.Place([]InstanceID{p1.Hand().Cards()[0].ID()}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	if got := len(p1.Field().Monsters()); got != MonsterSlots {
		t.Errorf("p1's monster area holds %d monsters, want %d", got, MonsterSlots)
	}
	assertInvariants(t, p1)
}

// finishedDuelWithP1Monster builds a Duel that p1 has just lost by deck-out,
// leaving one of p1's Monsters on the Field: the placement that emptied p1's
// hand could not be replenished (game-rules §8). It returns the finished Duel
// and the id of the Monster still on p1's Field.
func finishedDuelWithP1Monster(t *testing.T) (*Duel, InstanceID) {
	t.Helper()

	a := ParticipantSetup{ID: "p1", Deck: padsThenWeights(t, 2, 7, 10, 14), LifePoints: 4000}
	b := ParticipantSetup{ID: "p2", Deck: validDeck(t), LifePoints: 4000}
	d, err := NewDuel(a, b, noopShuffler{})
	if err != nil {
		t.Fatalf("NewDuel() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	for _, c := range slices.Clone(p1.Hand().Cards()) {
		if err := p1.DiscardFromHand(c.ID()); err != nil {
			t.Fatalf("DiscardFromHand() = %v", err)
		}
	}
	discardUntilDeckHas(t, p1, 4)
	drawN(t, p1, 4)
	sel := ids(p1.Hand().Cards())
	if err := d.Place(sel, PositionDefense, FaceUp); err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	if !d.Finished() {
		t.Fatal("expected the deck-out to finish the Duel")
	}
	return d, sel[len(sel)-1]
}
func TestPlace_AcceptsMoreTributesThanRequired(t *testing.T) {
	// A weight-7 Monster needs one tribute. Summoning it with two is legal:
	// game-rules §16 fixes the minimum, and offering more was an approved S4
	// decision. This test pins that decision.
	d := duelWithP1Deck(t, weightMonster(t, 1, 2), weightMonster(t, 2, 7))

	err := d.Place([]InstanceID{"p1#1", "p1#3", "p1#2"}, PositionAttack, FaceUp)

	if err != nil {
		t.Fatalf("Place() = %v, want success", err)
	}
	p1, _ := d.Participant("p1")
	if _, ok := p1.Field().Get("p1#2"); !ok {
		t.Fatal("the placed card p1#2 is not on the Field")
	}
	if got, want := ids(p1.Graveyard().Cards()), []InstanceID{"p1#1", "p1#3"}; !slices.Equal(got, want) {
		t.Errorf("Graveyard = %v, want %v in selection order", got, want)
	}
	assertInvariants(t, p1)
}

func TestPlace_AfterARejectedPlacementThePlacementIsStillAvailable(t *testing.T) {
	// Arrange
	d := duelWithP1Deck(t, weightMonster(t, 1, 2), weightMonster(t, 2, 7))

	// Act — summon the weight-7 Monster without its tribute.
	err := d.Place([]InstanceID{"p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrInvalidPlacement) {
		t.Fatalf("Place() error = %v, want ErrInvalidPlacement", err)
	}
	// The rejected placement did not consume the Phase 1 placement.
	if err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp); err != nil {
		t.Fatalf("second Place() = %v, want success", err)
	}
	p1, _ := d.Participant("p1")
	if _, ok := p1.Field().Get("p1#1"); !ok {
		t.Fatal("the card placed by the second action is not on the Field")
	}
}
