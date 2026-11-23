package domain

import "fmt"

// Position is the position of a Monster on the Field (game-rules §21). A
// Monster in Defense cannot attack. The zero value is not a position.
type Position int

const (
	// PositionAttack is the position of a Monster that can attack
	// (game-rules §21).
	PositionAttack Position = iota + 1

	// PositionDefense is the position of a Monster that cannot attack
	// (game-rules §21).
	PositionDefense
)

// String returns the position name used by the game rules.
func (p Position) String() string {
	switch p {
	case PositionAttack:
		return "Attack"
	case PositionDefense:
		return "Defense"
	default:
		return fmt.Sprintf("Position(%d)", int(p))
	}
}

// Face is whether a Monster on the Field is face-up or face-down
// (game-rules §18). A face-down Monster hides all of its information from the
// opponent (game-rules §18.2). The zero value is not a face.
type Face int

const (
	// FaceUp is the face of a Monster whose information is visible to both
	// Participants (game-rules §18.1).
	FaceUp Face = iota + 1

	// FaceDown is the face of a Monster whose information is hidden from the
	// opponent (game-rules §18.2).
	FaceDown
)

// String returns the face name used by the game rules.
func (f Face) String() string {
	switch f {
	case FaceUp:
		return "Face-up"
	case FaceDown:
		return "Face-down"
	default:
		return fmt.Sprintf("Face(%d)", int(f))
	}
}

// validPositionFace reports whether the pair is a coherent Monster position
// and face: both values must be valid, and a face-down Monster is always in
// Defense (game-rules §18.2). The callers wrap a failure with their own
// sentinel.
func validPositionFace(position Position, face Face) bool {
	if position != PositionAttack && position != PositionDefense {
		return false
	}
	if face != FaceUp && face != FaceDown {
		return false
	}
	return face != FaceDown || position == PositionDefense
}
