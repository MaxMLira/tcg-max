package domain

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

// equipmentCard builds a Magic Equipment definition that accepts the given
// Monster subtype (game-rules §35.2, one of the four MVP Equipment).
func equipmentCard(t *testing.T, n int, accepted MonsterSubtype) *Card {
	t.Helper()

	card, err := NewEquipmentCard(cardID(300+n), fmt.Sprintf("%s Equipment %d", accepted, n), accepted)
	if err != nil {
		t.Fatalf("NewEquipmentCard(%q) failed: %v", accepted, err)
	}
	return card
}

// subtypeMonster builds a Monster definition of the given subtype and weight 2,
// so its summon needs no tribute (game-rules §16).
func subtypeMonster(t *testing.T, n int, subtype MonsterSubtype) *Card {
	t.Helper()

	return mustCard(t, cardID(400+n), fmt.Sprintf("%s monster %d", subtype, n), CardTypeMonster, 1000, 800, 2, subtype)
}

// placeInPhase1 performs a single-card placement of the active Participant and
// fails the test when it is rejected.
func placeInPhase1(t *testing.T, d *Duel, id InstanceID) {
	t.Helper()

	if err := d.Place([]InstanceID{id}, PositionAttack, FaceUp); err != nil {
		t.Fatalf("Place(%q) returned unexpected error: %v", id, err)
	}
}

// duelWithEquipment builds a Duel where p1's deck starts with a Dragon
// Monster and then one Equipment per accepted subtype, so p1's opening hand
// is p1#1 (the Monster) followed by the Equipment. The Monster is summoned
// face-up in Attack on turn 1 and every Equipment is placed in the Magic/Trap
// area; because only one placement is allowed per Phase 1
// (game-rules §14.1), each one is placed on a different p1 turn. The Duel is
// left in p1's Phase 2, where the activation is allowed.
//
// It returns the Duel, the Monster id and the Equipment ids in placement
// order.
func duelWithEquipment(t *testing.T, accepted ...MonsterSubtype) (*Duel, InstanceID, []InstanceID) {
	t.Helper()

	head := []*Card{subtypeMonster(t, 1, SubtypeDragon())}
	for i, a := range accepted {
		head = append(head, equipmentCard(t, i+1, a))
	}
	d := duelWithP1Deck(t, head...)

	monster := InstanceID("p1#1")
	placeInPhase1(t, d, monster)

	equipment := make([]InstanceID, 0, len(accepted))
	for i := range accepted {
		endFullTurn(t, d) // -> p2
		endFullTurn(t, d) // -> p1, fresh Phase 1
		id := InstanceID(fmt.Sprintf("p1#%d", i+2))
		placeInPhase1(t, d, id)
		equipment = append(equipment, id)
	}
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	return d, monster, equipment
}

func TestPlace_PlacesAnEquipmentInAFreeMagicTrapZone(t *testing.T) {
	// Arrange — p1's hand: p1#1 is a Dragon Equipment.
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()))

	// Act — the position and face are ignored for an Equipment.
	err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got, want := ids(p1.Field().MagicTraps()), []InstanceID{"p1#1"}; !slices.Equal(got, want) {
		t.Fatalf("magic/trap area = %v, want %v", got, want)
	}
	if got := ids(p1.Field().Monsters()); len(got) != 0 {
		t.Errorf("monster area = %v, want empty: an Equipment never lands in the monster area", got)
	}
	eq, _ := p1.Field().Get("p1#1")
	if eq.Location() != LocationField {
		t.Errorf("Location() = %s, want %s", eq.Location(), LocationField)
	}
	if link := eq.EquipmentLink(); link != "" {
		t.Errorf("EquipmentLink() = %q, want empty: placement does not activate (mvp-story-mode §1.4)", link)
	}
	if got := p1.Field().EquipmentBonusFor("p1#1"); got != 0 {
		t.Errorf("EquipmentBonusFor() = %d, want 0 before activation", got)
	}
	if got := p1.Hand().Len(); got != HandLimit {
		t.Errorf("hand holds %d cards after replenishment, want %d", got, HandLimit)
	}
	assertInvariants(t, p1)
}

func TestPlace_EquipmentIgnoresThePositionAndFace(t *testing.T) {
	// The Magic/Trap area carries no position or face state (game-rules §9.2),
	// so a pair that would be rejected for a Monster is accepted here.
	tests := []struct {
		name     string
		position Position
		face     Face
	}{
		{name: "attack face-down is incoherent for a monster", position: PositionAttack, face: FaceDown},
		{name: "zero position and face", position: Position(0), face: Face(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()))

			// Act
			err := d.Place([]InstanceID{"p1#1"}, tt.position, tt.face)

			// Assert
			if err != nil {
				t.Fatalf("Place() returned unexpected error: %v", err)
			}
			p1, _ := d.Participant("p1")
			eq, ok := p1.Field().Get("p1#1")
			if !ok {
				t.Fatal("the placed Equipment is not on the Field")
			}
			if eq.Position() != Position(0) || eq.Face() != Face(0) {
				t.Errorf("Equipment = (%v, %v), want the zero position and face", eq.Position(), eq.Face())
			}
			assertInvariants(t, p1)
		})
	}
}

func TestPlace_RejectsAFifthEquipmentWhenTheFourZonesAreOccupied(t *testing.T) {
	// Arrange — place four Equipment, one per p1 turn (game-rules §14.1).
	d := duelWithP1Deck(t,
		equipmentCard(t, 1, SubtypeDragon()),
		equipmentCard(t, 2, SubtypeWarrior()),
		equipmentCard(t, 3, SubtypeMage()),
		equipmentCard(t, 4, SubtypeAquatic()),
		equipmentCard(t, 5, SubtypeDragon()),
	)
	for n := 1; n <= MagicTrapSlots; n++ {
		placeInPhase1(t, d, InstanceID(fmt.Sprintf("p1#%d", n)))
		endFullTurn(t, d) // -> p2
		endFullTurn(t, d) // -> p1
	}
	p1, _ := d.Participant("p1")
	if got := len(p1.Field().MagicTraps()); got != MagicTrapSlots {
		t.Fatalf("magic/trap area holds %d cards, want %d", got, MagicTrapSlots)
	}
	before := duelStateOf(d)

	// Act — the fifth Magic/Trap placement. Replacing an existing card is not
	// allowed (mvp-story-mode §1.4).
	err := d.Place([]InstanceID{"p1#5"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrFieldFull) {
		t.Fatalf("Place() error = %v, want ErrFieldFull (game-rules §9.2)", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a blocked placement")
	}
	assertInvariants(t, p1)
}

func TestPlace_AnEquipmentConsumesThePhase1Placement(t *testing.T) {
	// Arrange — placement and activation are distinct operations, and the
	// placement of an Equipment is the single placement of the Phase 1
	// (game-rules §14.1, mvp-story-mode §1.4).
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()), subtypeMonster(t, 1, SubtypeDragon()))
	placeInPhase1(t, d, "p1#1")
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrPlacementUsed) {
		t.Fatalf("second Place() error = %v, want ErrPlacementUsed", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by the second placement")
	}
}

func TestPlace_DiscardsTheSelectedMonstersBeforeAnEquipment(t *testing.T) {
	// Arrange — every selected card except the last goes to the Graveyard,
	// in selection order (game-rules §13.3); the last one is the Equipment.
	d := duelWithP1Deck(t,
		subtypeMonster(t, 1, SubtypeDragon()),
		subtypeMonster(t, 2, SubtypeWarrior()),
		equipmentCard(t, 1, SubtypeDragon()),
	)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2", "p1#3"}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got, want := ids(p1.Field().MagicTraps()), []InstanceID{"p1#3"}; !slices.Equal(got, want) {
		t.Errorf("magic/trap area = %v, want %v", got, want)
	}
	if got, want := ids(p1.Graveyard().Cards()), []InstanceID{"p1#1", "p1#2"}; !slices.Equal(got, want) {
		t.Errorf("Graveyard = %v, want %v in selection order", got, want)
	}
	assertInvariants(t, p1)
}

func TestPlace_StillRejectsAPlainMagicAndATrap(t *testing.T) {
	// A Magic that is not an Equipment declares no accepted subtype and has
	// no placement rule in the MVP (mvp-story-mode §1.4 and A2).
	tests := []struct {
		name string
		card *Card
	}{
		{name: "plain magic", card: magicCard(t, 1)},
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

func TestPlace_RejectsAnEquipmentAmongTheTributes(t *testing.T) {
	// Arrange — only Monsters may be selected before the last card
	// (game-rules §17, Story Mode MVP).
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()), weightMonster(t, 1, 7))
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("Place() error = %v, want ErrInvalidSelection", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected selection")
	}
}
func TestActivateEquipment_RejectedOutsidePhase2(t *testing.T) {
	// Arrange — the Equipment is placed and the turn is back in Phase 1 of a
	// later p1 turn. Activation is a Phase 2 action (game-rules §35.2).
	d, monster, equipment := duelWithEquipment(t, SubtypeDragon())
	if err := d.EndTurn(); err != nil { // the helper leaves p1 in Phase 2
		t.Fatalf("EndTurn() returned unexpected error: %v", err)
	}
	endFullTurn(t, d) // -> p1, Phase 1
	if got := d.Turn().Phase(); got != Phase1 {
		t.Fatalf("setup: phase = %v, want Phase 1", got)
	}
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment(equipment[0], monster)

	// Assert
	if !errors.Is(err, ErrWrongPhase) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrWrongPhase", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
	p1, _ := d.Participant("p1")
	if got := p1.Field().EquipmentBonusFor(monster); got != 0 {
		t.Errorf("EquipmentBonusFor() = %d, want 0 after a rejected activation", got)
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_RejectedFromTheHand(t *testing.T) {
	// Arrange — p1#2 is an Equipment still in the hand; p1#1 is a face-up
	// Dragon on the Field. An Equipment cannot be activated from the hand
	// (game-rules §35.2, mvp-story-mode §1.4).
	d := duelWithP1Deck(t, subtypeMonster(t, 1, SubtypeDragon()), equipmentCard(t, 1, SubtypeDragon()))
	placeInPhase1(t, d, "p1#1")
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if _, ok := p1.Hand().Get("p1#2"); !ok {
		t.Fatal("setup: p1#2 should still be in the hand")
	}
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment("p1#2", "p1#1")

	// Assert
	if !errors.Is(err, ErrInvalidEquipmentActivation) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
	if got := p1.Field().EquipmentBonusFor("p1#1"); got != 0 {
		t.Errorf("EquipmentBonusFor() = %d, want 0", got)
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_RejectedWhenNotPlacedOnTheField(t *testing.T) {
	// An Equipment must be placed before it can be activated
	// (game-rules §35.2.4): in the deck and in the Graveyard it is not.
	tests := []struct {
		name string
		// prepare returns the id of an Equipment that is not on the Field.
		prepare func(t *testing.T, d *Duel, p *Participant) InstanceID
	}{
		{
			name: "in the deck",
			prepare: func(t *testing.T, d *Duel, p *Participant) InstanceID {
				t.Helper()
				// The top of the deck is the first card that was never drawn.
				cards := p.Deck().Cards()
				if len(cards) == 0 {
					t.Fatal("setup: the deck is empty")
				}
				return cards[0].ID()
			},
		},
		{
			name: "in the graveyard",
			prepare: func(t *testing.T, d *Duel, p *Participant) InstanceID {
				t.Helper()
				// Discard the Equipment of the hand through a selection whose
				// last card is the Monster placed by the caller's setup.
				eq := InstanceID("p1#2")
				if err := p.DiscardFromHand(eq); err != nil {
					t.Fatalf("DiscardFromHand() returned unexpected error: %v", err)
				}
				return eq
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := duelWithP1Deck(t, subtypeMonster(t, 1, SubtypeDragon()), equipmentCard(t, 1, SubtypeDragon()))
			placeInPhase1(t, d, "p1#1")
			if err := d.EndPhase1(); err != nil {
				t.Fatalf("EndPhase1() returned unexpected error: %v", err)
			}
			p1, _ := d.Participant("p1")
			id := tt.prepare(t, d, p1)
			before := duelStateOf(d)

			// Act
			err := d.ActivateEquipment(id, "p1#1")

			// Assert
			if !errors.Is(err, ErrInvalidEquipmentActivation) {
				t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
				t.Errorf("Duel state changed by a rejected activation")
			}
			assertInvariants(t, p1)
		})
	}
}

func TestActivateEquipment_RejectsAnUnknownCard(t *testing.T) {
	// Arrange
	d, monster, _ := duelWithEquipment(t, SubtypeDragon())
	before := duelStateOf(d)

	// Act — an id of the other Participant is not a card of p1.
	err := d.ActivateEquipment("p2#1", monster)

	// Assert
	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrCardNotFound", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
}

func TestActivateEquipment_RejectedWithNoCompatibleFaceUpMonster(t *testing.T) {
	// Arrange — p1#1 is a Dragon Equipment and p1#2 a Dragon Monster that is
	// never summoned, so p1's monster area stays empty. Activation requires a
	// compatible face-up Monster on the Field (game-rules §35.2.4).
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()), subtypeMonster(t, 1, SubtypeDragon()))
	placeInPhase1(t, d, "p1#1")
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got := len(p1.Field().Monsters()); got != 0 {
		t.Fatalf("setup: monster area holds %d monsters, want 0", got)
	}
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment("p1#1", "p1#2")

	// Assert
	if !errors.Is(err, ErrInvalidEquipmentActivation) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_RejectsAFaceDownMonster(t *testing.T) {
	// Arrange — the only Dragon on the Field is face-down, so no compatible
	// face-up Monster exists (game-rules §35.2.4).
	d := duelWithP1Deck(t, subtypeMonster(t, 1, SubtypeDragon()), equipmentCard(t, 1, SubtypeDragon()))
	if err := d.Place([]InstanceID{"p1#1"}, PositionDefense, FaceDown); err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	endFullTurn(t, d) // -> p2
	endFullTurn(t, d) // -> p1
	placeInPhase1(t, d, "p1#2")
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	monster, _ := p1.Field().Get("p1#1")
	if monster.Face() != FaceDown {
		t.Fatalf("setup: the monster is %v, want face-down", monster.Face())
	}
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment("p1#2", "p1#1")

	// Assert
	if !errors.Is(err, ErrInvalidEquipmentActivation) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
	if got := p1.Field().EquipmentBonusFor("p1#1"); got != 0 {
		t.Errorf("EquipmentBonusFor() = %d, want 0", got)
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_RejectsAnIncompatibleSubtype(t *testing.T) {
	// Each Equipment accepts exactly one Monster subtype (game-rules §35.2
	// and §10.1.1), so every other subtype is rejected.
	for _, accepted := range []MonsterSubtype{SubtypeWarrior(), SubtypeMage(), SubtypeAquatic()} {
		t.Run(accepted.String(), func(t *testing.T) {
			// Arrange — a face-up Dragon on the Field and an Equipment that
			// accepts another subtype.
			d := duelWithP1Deck(t, subtypeMonster(t, 1, SubtypeDragon()), equipmentCard(t, 1, accepted))
			placeInPhase1(t, d, "p1#1")
			endFullTurn(t, d) // -> p2
			endFullTurn(t, d) // -> p1
			placeInPhase1(t, d, "p1#2")
			if err := d.EndPhase1(); err != nil {
				t.Fatalf("EndPhase1() returned unexpected error: %v", err)
			}
			before := duelStateOf(d)

			// Act
			err := d.ActivateEquipment("p1#2", "p1#1")

			// Assert
			if !errors.Is(err, ErrInvalidEquipmentActivation) {
				t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
				t.Errorf("Duel state changed by a rejected activation")
			}
			p1, _ := d.Participant("p1")
			if got := p1.Field().EquipmentBonusFor("p1#1"); got != 0 {
				t.Errorf("EquipmentBonusFor() = %d, want 0 for an incompatible subtype", got)
			}
			assertInvariants(t, p1)
		})
	}
}

func TestActivateEquipment_RejectsATargetThatIsNotAMonster(t *testing.T) {
	// Arrange — two Equipment on the Field; one is used as the target.
	d, _, equipment := duelWithEquipment(t, SubtypeDragon(), SubtypeDragon())
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment(equipment[0], equipment[1])

	// Assert
	if !errors.Is(err, ErrInvalidEquipmentActivation) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
}

func TestActivateEquipment_RejectsAMonsterOfTheOpponent(t *testing.T) {
	// Arrange — p2 summons a compatible Dragon; p1 may not equip it
	// (game-rules §35.2.4: the Monster is on the player's own field).
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()))
	placeInPhase1(t, d, "p1#1")
	endFullTurn(t, d) // -> p2
	p2, _ := d.Participant("p2")
	p2Monster := p2.Hand().Cards()[0].ID()
	placeInPhase1(t, d, p2Monster)
	endFullTurn(t, d) // -> p1
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment("p1#1", p2Monster)

	// Assert
	if !errors.Is(err, ErrInvalidEquipmentActivation) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
}

func TestActivateEquipment_LinksTheEquipmentAndGrantsPlus800ATK(t *testing.T) {
	// Arrange
	d, monster, equipment := duelWithEquipment(t, SubtypeDragon())
	p1, _ := d.Participant("p1")

	// Act
	err := d.ActivateEquipment(equipment[0], monster)

	// Assert
	if err != nil {
		t.Fatalf("ActivateEquipment() returned unexpected error: %v", err)
	}
	if got, want := p1.Field().EquipmentBonusFor(monster), EquipmentATKBonus; got != want {
		t.Errorf("EquipmentBonusFor() = %d, want %d (game-rules §35.2.1)", got, want)
	}
	eq, ok := p1.Field().Get(equipment[0])
	if !ok {
		t.Fatal("the activated Equipment left the Field")
	}
	if got := eq.EquipmentLink(); got != monster {
		t.Errorf("EquipmentLink() = %q, want %q", got, monster)
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_StaysInItsMagicTrapZone(t *testing.T) {
	// Arrange — after the activation the Equipment stays on the Field and
	// keeps occupying its Magic/Trap slot (game-rules §35.2.4).
	d, monster, equipment := duelWithEquipment(t, SubtypeDragon())
	p1, _ := d.Participant("p1")
	zoneBefore := ids(p1.Field().MagicTraps())

	// Act
	if err := d.ActivateEquipment(equipment[0], monster); err != nil {
		t.Fatalf("ActivateEquipment() returned unexpected error: %v", err)
	}

	// Assert
	if got := ids(p1.Field().MagicTraps()); !slices.Equal(got, zoneBefore) {
		t.Errorf("magic/trap area = %v, want %v unchanged by the activation", got, zoneBefore)
	}
	eq, _ := p1.Field().Get(equipment[0])
	if eq.Location() != LocationField {
		t.Errorf("Location() = %s, want %s: the Equipment does not go to the Graveyard on activation", eq.Location(), LocationField)
	}
	if p1.Graveyard().Len() != 0 {
		t.Errorf("Graveyard = %v, want empty", ids(p1.Graveyard().Cards()))
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_BonusesStackOnTheSameMonster(t *testing.T) {
	// Arrange — three compatible Equipment on the Field (game-rules §35.2.2:
	// several Equipment may be linked to the same Monster and the bonuses add
	// up).
	d, monster, equipment := duelWithEquipment(t, SubtypeDragon(), SubtypeDragon(), SubtypeDragon())
	p1, _ := d.Participant("p1")

	// Act and assert — the bonus grows by EquipmentATKBonus per activation.
	for n, id := range equipment {
		if err := d.ActivateEquipment(id, monster); err != nil {
			t.Fatalf("ActivateEquipment(%q) returned unexpected error: %v", id, err)
		}
		if got, want := p1.Field().EquipmentBonusFor(monster), (n+1)*EquipmentATKBonus; got != want {
			t.Errorf("after %d activations EquipmentBonusFor() = %d, want %d", n+1, got, want)
		}
	}
	if got := len(p1.Field().MagicTraps()); got != len(equipment) {
		t.Errorf("magic/trap area holds %d cards, want %d", got, len(equipment))
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_BonusAppliesOnlyToTheLinkedMonster(t *testing.T) {
	// Arrange — p1 holds two face-up Dragons and one Equipment.
	d := duelWithP1Deck(t,
		subtypeMonster(t, 1, SubtypeDragon()),
		subtypeMonster(t, 2, SubtypeDragon()),
		equipmentCard(t, 1, SubtypeDragon()),
	)
	placeInPhase1(t, d, "p1#1")
	endFullTurn(t, d)
	endFullTurn(t, d)
	placeInPhase1(t, d, "p1#2")
	endFullTurn(t, d)
	endFullTurn(t, d)
	placeInPhase1(t, d, "p1#3")
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}

	// Act
	if err := d.ActivateEquipment("p1#3", "p1#1"); err != nil {
		t.Fatalf("ActivateEquipment() returned unexpected error: %v", err)
	}

	// Assert
	p1, _ := d.Participant("p1")
	if got, want := p1.Field().EquipmentBonusFor("p1#1"), EquipmentATKBonus; got != want {
		t.Errorf("bonus of the equipped monster = %d, want %d", got, want)
	}
	if got := p1.Field().EquipmentBonusFor("p1#2"); got != 0 {
		t.Errorf("bonus of the other monster = %d, want 0", got)
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_AllowedInThePhase2OfThePlacementTurn(t *testing.T) {
	// Arrange — the Equipment is placed in Phase 1 and activated in the
	// Phase 2 of the same turn: the rule only requires that the placement
	// precedes the activation (game-rules §35.2.4).
	d := duelWithP1Deck(t, subtypeMonster(t, 1, SubtypeDragon()), equipmentCard(t, 1, SubtypeDragon()))
	placeInPhase1(t, d, "p1#1")
	endFullTurn(t, d) // -> p2
	endFullTurn(t, d) // -> p1
	placeInPhase1(t, d, "p1#2")
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}

	// Act
	err := d.ActivateEquipment("p1#2", "p1#1")

	// Assert
	if err != nil {
		t.Fatalf("ActivateEquipment() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got, want := p1.Field().EquipmentBonusFor("p1#1"), EquipmentATKBonus; got != want {
		t.Errorf("EquipmentBonusFor() = %d, want %d", got, want)
	}
}

func TestActivateEquipment_RejectsASecondActivationOfTheSameEquipment(t *testing.T) {
	// Arrange — the MVP defines no re-link: once linked, the Equipment stays
	// linked to its Monster (game-rules §35.2.4).
	d, monster, equipment := duelWithEquipment(t, SubtypeDragon())
	if err := d.ActivateEquipment(equipment[0], monster); err != nil {
		t.Fatalf("ActivateEquipment() returned unexpected error: %v", err)
	}
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment(equipment[0], monster)

	// Assert
	if !errors.Is(err, ErrInvalidEquipmentActivation) {
		t.Fatalf("second ActivateEquipment() error = %v, want ErrInvalidEquipmentActivation", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected activation")
	}
	p1, _ := d.Participant("p1")
	if got, want := p1.Field().EquipmentBonusFor(monster), EquipmentATKBonus; got != want {
		t.Errorf("EquipmentBonusFor() = %d, want %d: the rejected activation must not add a bonus", got, want)
	}
}

func TestActivateEquipment_LinkSurvivesTheEndOfTheTurn(t *testing.T) {
	// Arrange — the bonus persists while the Equipment stays equipped
	// (game-rules §38): it is not turn state.
	d, monster, equipment := duelWithEquipment(t, SubtypeDragon())
	if err := d.ActivateEquipment(equipment[0], monster); err != nil {
		t.Fatalf("ActivateEquipment() returned unexpected error: %v", err)
	}

	// Act — two full turns later it is p1's turn again.
	if err := d.EndTurn(); err != nil {
		t.Fatalf("EndTurn() returned unexpected error: %v", err)
	}
	endFullTurn(t, d) // -> p1

	// Assert
	p1, _ := d.Participant("p1")
	if got, want := p1.Field().EquipmentBonusFor(monster), EquipmentATKBonus; got != want {
		t.Errorf("EquipmentBonusFor() = %d, want %d after the turn ended", got, want)
	}
	eq, _ := p1.Field().Get(equipment[0])
	if got := eq.EquipmentLink(); got != monster {
		t.Errorf("EquipmentLink() = %q, want %q", got, monster)
	}
	assertInvariants(t, p1)
}

func TestActivateEquipment_RejectedInAFinishedDuel(t *testing.T) {
	// Arrange
	d, _ := finishedDuelWithP1Monster(t)
	before := duelStateOf(d)

	// Act
	err := d.ActivateEquipment("p1#1", "p1#2")

	// Assert
	if !errors.Is(err, ErrDuelEnded) {
		t.Fatalf("ActivateEquipment() error = %v, want ErrDuelEnded", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed after the Duel ended")
	}
}
