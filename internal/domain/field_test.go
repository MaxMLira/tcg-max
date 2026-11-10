package domain

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

// magicCard builds a valid Magic definition.
func magicCard(t *testing.T, n int) *Card {
	t.Helper()

	return mustCard(t, cardID(500+n), fmt.Sprintf("Magic %d", n), CardTypeMagic, 0, 0, 0, MonsterSubtype{})
}

// trapCard builds a valid Trap definition. Trap is a valid domain card type
// even though no Trap is part of the MVP decks.
func trapCard(t *testing.T, n int) *Card {
	t.Helper()

	return mustCard(t, cardID(600+n), fmt.Sprintf("Trap %d", n), CardTypeTrap, 100, 100, 1, MonsterSubtype{})
}

// fieldInstance builds an instance with id "f#<n>" of the given definition.
func fieldInstance(t *testing.T, n int, card *Card) *CardInstance {
	t.Helper()

	return instanceOf(t, InstanceID(fmt.Sprintf("f#%d", n)), card)
}

func TestField_MonsterAreaAcceptsFourAndRejectsTheFifth(t *testing.T) {
	// Arrange
	field := newField()
	for i := 1; i <= MonsterSlots; i++ {
		if err := field.add(fieldInstance(t, i, deckCard(t, i))); err != nil {
			t.Fatalf("add() #%d returned unexpected error: %v", i, err)
		}
	}
	before := ids(field.Monsters())
	fifth := fieldInstance(t, 9, deckCard(t, 9))

	// Act
	err := field.add(fifth)

	// Assert
	if MonsterSlots != 4 {
		t.Errorf("MonsterSlots = %d, want 4 (game-rules §9.1)", MonsterSlots)
	}
	if !errors.Is(err, ErrFieldFull) {
		t.Fatalf("add() error = %v, want ErrFieldFull", err)
	}
	if got := ids(field.Monsters()); !slices.Equal(got, before) {
		t.Errorf("Monsters() = %v, want %v (unchanged)", got, before)
	}
	if len(field.MagicTraps()) != 0 {
		t.Errorf("MagicTraps() = %d cards, want 0", len(field.MagicTraps()))
	}
	if field.HasRoomFor(deckCard(t, 9)) {
		t.Errorf("HasRoomFor(monster) = true on a full monster area")
	}
}

func TestField_MagicTrapAreaAcceptsFourAndRejectsTheFifth(t *testing.T) {
	tests := []struct {
		name  string
		build func(t *testing.T, n int) *Card
	}{
		{name: "magic cards", build: magicCard},
		{name: "trap cards", build: trapCard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			field := newField()
			for i := 1; i <= MagicTrapSlots; i++ {
				if err := field.add(fieldInstance(t, i, tt.build(t, i))); err != nil {
					t.Fatalf("add() #%d returned unexpected error: %v", i, err)
				}
			}
			before := ids(field.MagicTraps())

			// Act
			err := field.add(fieldInstance(t, 9, tt.build(t, 9)))

			// Assert
			if MagicTrapSlots != 4 {
				t.Errorf("MagicTrapSlots = %d, want 4 (game-rules §9.2)", MagicTrapSlots)
			}
			if !errors.Is(err, ErrFieldFull) {
				t.Fatalf("add() error = %v, want ErrFieldFull", err)
			}
			if got := ids(field.MagicTraps()); !slices.Equal(got, before) {
				t.Errorf("MagicTraps() = %v, want %v (unchanged)", got, before)
			}
			if field.HasRoomFor(tt.build(t, 9)) {
				t.Errorf("HasRoomFor() = true on a full magic/trap area")
			}
		})
	}
}

func TestField_MagicAndTrapShareTheMagicTrapArea(t *testing.T) {
	// Arrange
	field := newField()

	// Act
	errs := []error{
		field.add(fieldInstance(t, 1, magicCard(t, 1))),
		field.add(fieldInstance(t, 2, trapCard(t, 2))),
		field.add(fieldInstance(t, 3, magicCard(t, 3))),
		field.add(fieldInstance(t, 4, trapCard(t, 4))),
	}
	errFifth := field.add(fieldInstance(t, 5, trapCard(t, 5)))

	// Assert
	for i, err := range errs {
		if err != nil {
			t.Fatalf("add() #%d returned unexpected error: %v", i+1, err)
		}
	}
	if !errors.Is(errFifth, ErrFieldFull) {
		t.Errorf("fifth magic/trap error = %v, want ErrFieldFull (one shared area of %d)", errFifth, MagicTrapSlots)
	}
	if got, want := ids(field.MagicTraps()), []InstanceID{"f#1", "f#2", "f#3", "f#4"}; !slices.Equal(got, want) {
		t.Errorf("MagicTraps() = %v, want %v (placement order)", got, want)
	}
}

func TestField_AreaIsDerivedFromTheCardType(t *testing.T) {
	tests := []struct {
		name         string
		card         func(t *testing.T) *Card
		wantMonsters int
		wantMagicTrp int
	}{
		{name: "monster goes to the monster area", card: func(t *testing.T) *Card { return deckCard(t, 1) }, wantMonsters: 1},
		{name: "magic goes to the magic/trap area", card: func(t *testing.T) *Card { return magicCard(t, 1) }, wantMagicTrp: 1},
		{name: "trap goes to the magic/trap area", card: func(t *testing.T) *Card { return trapCard(t, 1) }, wantMagicTrp: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			field := newField()

			// Act
			err := field.add(fieldInstance(t, 1, tt.card(t)))

			// Assert
			if err != nil {
				t.Fatalf("add() returned unexpected error: %v", err)
			}
			if got := len(field.Monsters()); got != tt.wantMonsters {
				t.Errorf("len(Monsters()) = %d, want %d", got, tt.wantMonsters)
			}
			if got := len(field.MagicTraps()); got != tt.wantMagicTrp {
				t.Errorf("len(MagicTraps()) = %d, want %d", got, tt.wantMagicTrp)
			}
		})
	}
}

func TestField_MonsterNeverLandsInMagicTrapAreaAndViceVersa(t *testing.T) {
	// Arrange
	field := newField()
	for i := 1; i <= MonsterSlots; i++ {
		if err := field.add(fieldInstance(t, i, deckCard(t, i))); err != nil {
			t.Fatalf("add() monster #%d returned unexpected error: %v", i, err)
		}
	}

	// Act: the monster area is full, but the magic/trap area is still empty.
	errMonster := field.add(fieldInstance(t, 10, deckCard(t, 10)))
	errMagic := field.add(fieldInstance(t, 11, magicCard(t, 11)))
	errTrap := field.add(fieldInstance(t, 12, trapCard(t, 12)))

	// Assert
	if !errors.Is(errMonster, ErrFieldFull) {
		t.Errorf("fifth monster error = %v, want ErrFieldFull (it must not overflow into the magic/trap area)", errMonster)
	}
	if errMagic != nil || errTrap != nil {
		t.Fatalf("magic/trap add errors = %v, %v, want nil (the monster area being full is irrelevant)", errMagic, errTrap)
	}
	for _, c := range field.Monsters() {
		if c.Card().Type() != CardTypeMonster {
			t.Errorf("monster area holds %s %q", c.Card().Type(), c.ID())
		}
	}
	for _, c := range field.MagicTraps() {
		if c.Card().Type() == CardTypeMonster {
			t.Errorf("magic/trap area holds a Monster %q", c.ID())
		}
	}
	if got := len(field.Monsters()); got != MonsterSlots {
		t.Errorf("len(Monsters()) = %d, want %d", got, MonsterSlots)
	}
	if got := len(field.MagicTraps()); got != 2 {
		t.Errorf("len(MagicTraps()) = %d, want 2", got)
	}
}

func TestField_HasRoomForIsPerArea(t *testing.T) {
	// Arrange
	field := newField()
	for i := 1; i <= MonsterSlots; i++ {
		if err := field.add(fieldInstance(t, i, deckCard(t, i))); err != nil {
			t.Fatalf("add() returned unexpected error: %v", err)
		}
	}

	// Act and Assert
	if field.HasRoomFor(deckCard(t, 9)) {
		t.Errorf("HasRoomFor(monster) = true, want false (monster area full)")
	}
	if !field.HasRoomFor(magicCard(t, 1)) {
		t.Errorf("HasRoomFor(magic) = false, want true")
	}
	if !field.HasRoomFor(trapCard(t, 1)) {
		t.Errorf("HasRoomFor(trap) = false, want true")
	}
	if field.HasRoomFor(nil) {
		t.Errorf("HasRoomFor(nil) = true, want false")
	}
}

func TestField_GetAndRemove(t *testing.T) {
	// Arrange
	field := newField()
	monster := fieldInstance(t, 1, deckCard(t, 1))
	magic := fieldInstance(t, 2, magicCard(t, 2))
	for _, inst := range []*CardInstance{monster, magic} {
		if err := field.add(inst); err != nil {
			t.Fatalf("add() returned unexpected error: %v", err)
		}
	}

	// Act
	gotMonster, okMonster := field.Get("f#1")
	_, okMissing := field.Get("f#9")
	removed, err := field.remove("f#2")

	// Assert
	if !okMonster || gotMonster != monster {
		t.Errorf("Get(f#1) = %v, %v, want the monster instance", gotMonster, okMonster)
	}
	if okMissing {
		t.Errorf("Get(f#9) found a card that is not on the field")
	}
	if err != nil || removed != magic {
		t.Errorf("remove(f#2) = %v, %v, want the magic instance", removed, err)
	}
	if len(field.MagicTraps()) != 0 || len(field.Monsters()) != 1 {
		t.Errorf("after remove: monsters=%d magicTraps=%d, want 1 and 0", len(field.Monsters()), len(field.MagicTraps()))
	}
}

func TestField_RemoveOfAbsentCardFailsAndStaysUnchanged(t *testing.T) {
	// Arrange
	field := newField()
	if err := field.add(fieldInstance(t, 1, deckCard(t, 1))); err != nil {
		t.Fatalf("add() returned unexpected error: %v", err)
	}

	// Act
	removed, err := field.remove("f#9")

	// Assert
	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("remove() error = %v, want ErrCardNotFound", err)
	}
	if removed != nil {
		t.Errorf("remove() = %v, want nil on error", removed)
	}
	if got, want := ids(field.Monsters()), []InstanceID{"f#1"}; !slices.Equal(got, want) {
		t.Errorf("Monsters() = %v, want %v (unchanged)", got, want)
	}
}

func TestField_ReadersReturnCopies(t *testing.T) {
	// Arrange
	field := newField()
	if err := field.add(fieldInstance(t, 1, deckCard(t, 1))); err != nil {
		t.Fatalf("add() returned unexpected error: %v", err)
	}
	if err := field.add(fieldInstance(t, 2, magicCard(t, 2))); err != nil {
		t.Fatalf("add() returned unexpected error: %v", err)
	}

	// Act
	monsters := field.Monsters()
	magicTraps := field.MagicTraps()
	monsters[0] = nil
	magicTraps[0] = nil

	// Assert
	if got, want := ids(field.Monsters()), []InstanceID{"f#1"}; !slices.Equal(got, want) {
		t.Errorf("Monsters() = %v, want %v", got, want)
	}
	if got, want := ids(field.MagicTraps()), []InstanceID{"f#2"}; !slices.Equal(got, want) {
		t.Errorf("MagicTraps() = %v, want %v", got, want)
	}
}
