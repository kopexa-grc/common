// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    DateTime
		wantErr bool
	}{
		{
			name:  "nil value",
			input: nil,
			want:  DateTime{},
		},
		{
			name:  "valid time",
			input: time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC),
			want:  DateTime(time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC)),
		},
		{
			name:    "invalid type",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DateTime
			err := dt.Scan(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, dt)
		})
	}
}

func TestDateTime_Value(t *testing.T) {
	now := time.Now()
	dt := DateTime(now)

	value, err := dt.Value()
	assert.NoError(t, err)
	assert.Equal(t, now, value)
}

func TestDateTime_UnmarshalCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DateTime
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  DateTime{},
		},
		{
			name:  "valid ISO date",
			input: "2024-03-20T15:04:05Z",
			want:  DateTime(time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC)),
		},
		{
			name:  "valid simple date",
			input: "2024-03-20",
			want:  DateTime(time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:    "invalid date",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DateTime
			err := dt.UnmarshalCSV(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, dt)
		})
	}
}

func TestDateTime_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    DateTime
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  DateTime{},
		},
		{
			name:  "valid ISO date",
			input: "2024-03-20T15:04:05Z",
			want:  DateTime(time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC)),
		},
		{
			name:  "valid simple date",
			input: "2024-03-20",
			want:  DateTime(time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DateTime
			err := dt.UnmarshalGQL(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, dt)
		})
	}
}

func TestDateTime_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		input DateTime
		want  string
	}{
		{
			name:  "zero time",
			input: DateTime{},
			want:  `""`,
		},
		{
			name:  "valid time",
			input: DateTime(time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC)),
			want:  `"2024-03-20T15:04:05Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			tt.input.MarshalGQL(&buf)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestDateTime_String(t *testing.T) {
	tests := []struct {
		name  string
		input DateTime
		want  string
	}{
		{
			name:  "zero time",
			input: DateTime{},
			want:  "",
		},
		{
			name:  "valid time",
			input: DateTime(time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC)),
			want:  "2024-03-20T15:04:05Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.String())
		})
	}
}

func TestToDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DateTime
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:  "valid ISO date",
			input: "2024-03-20T15:04:05Z",
			want:  DateTime(time.Date(2024, 3, 20, 15, 4, 5, 0, time.UTC)),
		},
		{
			name:  "valid simple date",
			input: "2024-03-20",
			want:  DateTime(time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:    "invalid date",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDateTime(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, *got)
		})
	}
}
