package domain

import "testing"

func TestPositionString(t *testing.T) {
	tests := []struct {
		position Position
		want     string
	}{
		{PositionAttack, "Attack"},
		{PositionDefense, "Defense"},
		{Position(0), "Position(0)"},
		{Position(9), "Position(9)"},
	}
	for _, tt := range tests {
		if got := tt.position.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestFaceString(t *testing.T) {
	tests := []struct {
		face Face
		want string
	}{
		{FaceUp, "Face-up"},
		{FaceDown, "Face-down"},
		{Face(0), "Face(0)"},
	}
	for _, tt := range tests {
		if got := tt.face.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestValidPositionFace(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		face     Face
		want     bool
	}{
		{name: "face-up attack", position: PositionAttack, face: FaceUp, want: true},
		{name: "face-up defense", position: PositionDefense, face: FaceUp, want: true},
		{name: "face-down defense", position: PositionDefense, face: FaceDown, want: true},
		{name: "face-down attack", position: PositionAttack, face: FaceDown, want: false},
		{name: "zero position", position: Position(0), face: FaceUp, want: false},
		{name: "zero face", position: PositionAttack, face: Face(0), want: false},
		{name: "unknown values", position: Position(9), face: Face(9), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validPositionFace(tt.position, tt.face); got != tt.want {
				t.Errorf("validPositionFace(%v, %v) = %v, want %v", tt.position, tt.face, got, tt.want)
			}
		})
	}
}
