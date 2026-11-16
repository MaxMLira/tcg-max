package domain

import (
	"errors"
	"math/rand/v2"
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
