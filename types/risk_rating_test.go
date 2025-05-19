// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRiskRating_CalculateRatings(t *testing.T) {
	tests := []struct {
		name        string
		likelihood  int
		consequence int
		want        int
	}{
		{
			name:        "valid rating",
			likelihood:  3,
			consequence: 4,
			want:        12,
		},
		{
			name:        "zero values",
			likelihood:  0,
			consequence: 0,
			want:        0,
		},
		{
			name:        "minimum values",
			likelihood:  MinRiskValue,
			consequence: MinRiskValue,
			want:        MinRiskValue * MinRiskValue,
		},
		{
			name:        "maximum values",
			likelihood:  MaxRiskValue,
			consequence: MaxRiskValue,
			want:        MaxRiskValue * MaxRiskValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskRating{
				Likelihood:  tt.likelihood,
				Consequence: tt.consequence,
			}
			r.CalculateRatings()
			assert.Equal(t, tt.want, r.Rating)
		})
	}
}

func TestRiskRating_IsComplete(t *testing.T) {
	tests := []struct {
		name        string
		likelihood  int
		consequence int
		want        bool
	}{
		{
			name:        "complete",
			likelihood:  3,
			consequence: 4,
			want:        true,
		},
		{
			name:        "zero likelihood",
			likelihood:  0,
			consequence: 4,
			want:        false,
		},
		{
			name:        "zero consequence",
			likelihood:  3,
			consequence: 0,
			want:        false,
		},
		{
			name:        "both zero",
			likelihood:  0,
			consequence: 0,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskRating{
				Likelihood:  tt.likelihood,
				Consequence: tt.consequence,
			}
			assert.Equal(t, tt.want, r.IsComplete())
		})
	}
}

func TestRiskRating_IsZero(t *testing.T) {
	tests := []struct {
		name        string
		likelihood  int
		consequence int
		want        bool
	}{
		{
			name:        "zero",
			likelihood:  0,
			consequence: 0,
			want:        true,
		},
		{
			name:        "non-zero",
			likelihood:  3,
			consequence: 4,
			want:        false,
		},
		{
			name:        "partial zero",
			likelihood:  0,
			consequence: 4,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskRating{
				Likelihood:  tt.likelihood,
				Consequence: tt.consequence,
			}
			assert.Equal(t, tt.want, r.IsZero())
		})
	}
}

func TestRiskRating_IsInvalid(t *testing.T) {
	tests := []struct {
		name        string
		likelihood  int
		consequence int
		want        bool
	}{
		{
			name:        "valid",
			likelihood:  3,
			consequence: 4,
			want:        false,
		},
		{
			name:        "zero values",
			likelihood:  0,
			consequence: 0,
			want:        true,
		},
		{
			name:        "likelihood too low",
			likelihood:  MinRiskValue - 1,
			consequence: 3,
			want:        true,
		},
		{
			name:        "likelihood too high",
			likelihood:  MaxRiskValue + 1,
			consequence: 3,
			want:        true,
		},
		{
			name:        "consequence too low",
			likelihood:  3,
			consequence: MinRiskValue - 1,
			want:        true,
		},
		{
			name:        "consequence too high",
			likelihood:  3,
			consequence: MaxRiskValue + 1,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskRating{
				Likelihood:  tt.likelihood,
				Consequence: tt.consequence,
			}
			assert.Equal(t, tt.want, r.IsInvalid())
		})
	}
}

func TestRiskRating_Score(t *testing.T) {
	tests := []struct {
		name        string
		likelihood  int
		consequence int
		want        int
	}{
		{
			name:        "valid score",
			likelihood:  3,
			consequence: 4,
			want:        12,
		},
		{
			name:        "zero values",
			likelihood:  0,
			consequence: 0,
			want:        DefaultScore,
		},
		{
			name:        "minimum values",
			likelihood:  MinRiskValue,
			consequence: MinRiskValue,
			want:        MinRiskValue * MinRiskValue,
		},
		{
			name:        "maximum values",
			likelihood:  MaxRiskValue,
			consequence: MaxRiskValue,
			want:        MaxRiskValue * MaxRiskValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskRating{
				Likelihood:  tt.likelihood,
				Consequence: tt.consequence,
			}
			assert.Equal(t, tt.want, r.Score())
		})
	}
}

func TestRiskRating_String(t *testing.T) {
	tests := []struct {
		name        string
		likelihood  int
		consequence int
		want        string
	}{
		{
			name:        "valid values",
			likelihood:  3,
			consequence: 4,
			want:        "L3C4",
		},
		{
			name:        "zero values",
			likelihood:  0,
			consequence: 0,
			want:        "L0C0",
		},
		{
			name:        "minimum values",
			likelihood:  MinRiskValue,
			consequence: MinRiskValue,
			want:        "L1C1",
		},
		{
			name:        "maximum values",
			likelihood:  MaxRiskValue,
			consequence: MaxRiskValue,
			want:        "L5C5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskRating{
				Likelihood:  tt.likelihood,
				Consequence: tt.consequence,
			}
			assert.Equal(t, tt.want, r.String())
		})
	}
}

func TestRiskRating_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   RiskRating
		want    string
		wantErr bool
	}{
		{
			name: "valid rating",
			input: RiskRating{
				Likelihood:  3,
				Consequence: 4,
				Rating:      12,
				Comment:     "test comment",
			},
			want:    `{"likelihood":3,"consequence":4,"rating":12,"comment":"test comment"}`,
			wantErr: false,
		},
		{
			name: "empty rating",
			input: RiskRating{
				Likelihood:  0,
				Consequence: 0,
				Rating:      0,
			},
			want:    `{"likelihood":0,"consequence":0,"rating":0}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))

			var unmarshaled RiskRating
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}
