package domain

import (
	"errors"
	"testing"
)

func TestInitialLifePoints_ByLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   int
		want    int
		wantErr error
	}{
		{name: "negative level", level: -1, wantErr: ErrInvalidLevel},
		{name: "level zero", level: 0, wantErr: ErrInvalidLevel},
		{name: "level 1", level: 1, want: 4000},
		{name: "level 2", level: 2, want: 4000},
		{name: "level 3", level: 3, want: 4000},
		{name: "level 4", level: 4, want: 4000},
		{name: "level 5", level: 5, want: 8000},
		{name: "level 6 not supported yet", level: 6, wantErr: ErrInvalidLevel},
		{name: "level 10 not supported yet", level: 10, wantErr: ErrInvalidLevel},
		{name: "level 11", level: 11, wantErr: ErrInvalidLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, err := InitialLifePoints(tt.level)

			// Assert
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("InitialLifePoints(%d) error = %v, want %v", tt.level, err, tt.wantErr)
				}
				if got != 0 {
					t.Errorf("InitialLifePoints(%d) = %d, want 0 on error", tt.level, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("InitialLifePoints(%d) returned unexpected error: %v", tt.level, err)
			}
			if got != tt.want {
				t.Errorf("InitialLifePoints(%d) = %d, want %d", tt.level, got, tt.want)
			}
		})
	}
}

func TestInitialLifePoints_MaxSupportedLevelIsTheLastAcceptedLevel(t *testing.T) {
	// Act
	_, errAtMax := InitialLifePoints(MaxSupportedLevel)
	_, errAbove := InitialLifePoints(MaxSupportedLevel + 1)

	// Assert
	if errAtMax != nil {
		t.Errorf("InitialLifePoints(MaxSupportedLevel) returned unexpected error: %v", errAtMax)
	}
	if !errors.Is(errAbove, ErrInvalidLevel) {
		t.Errorf("InitialLifePoints(MaxSupportedLevel+1) error = %v, want ErrInvalidLevel", errAbove)
	}
}
