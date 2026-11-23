package domain

import "fmt"

// Phase is a stage of the current Turn (glossary "Phase", game-rules §11,
// Story Mode MVP). Phase 1 handles placement and position-related actions
// (game-rules §14); Phase 2 handles attacks and Magic activation
// (game-rules §23, tasks S5 and S6). The zero value is not a phase.
type Phase int

const (
	// Phase1 is the placement phase (game-rules §14).
	Phase1 Phase = iota + 1

	// Phase2 is the action phase (game-rules §23).
	Phase2
)

// String returns the phase name used by the game rules.
func (p Phase) String() string {
	switch p {
	case Phase1:
		return "Phase 1"
	case Phase2:
		return "Phase 2"
	default:
		return fmt.Sprintf("Phase(%d)", int(p))
	}
}
