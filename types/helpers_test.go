// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalGQLJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		writer  func() (w *bytes.Buffer, asNil bool)
		want    string
		wantErr bool
	}{
		{
			name:   "valid string",
			input:  "test",
			writer: func() (*bytes.Buffer, bool) { return &bytes.Buffer{}, false },
			want:   `"test"`,
		},
		{
			name:   "valid struct",
			input:  struct{ Name string }{"test"},
			writer: func() (*bytes.Buffer, bool) { return &bytes.Buffer{}, false },
			want:   `{"Name":"test"}`,
		},
		{
			name:    "nil writer",
			input:   "test",
			writer:  func() (*bytes.Buffer, bool) { return nil, true },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, asNil := tt.writer()
			var err error
			if asNil {
				err = marshalGQLJSON[any](nil, tt.input)
				assert.Error(t, err)
				return
			}
			err = marshalGQLJSON(w, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, w.String())
		})
	}
}

func TestUnmarshalGQLJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		target  interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "nil value",
			input:   nil,
			target:  new(string),
			wantErr: true,
		},
		{
			name:   "valid string",
			input:  "test",
			target: new(string),
			want:   "test",
		},
		{
			name:   "valid struct",
			input:  struct{ Name string }{"test"},
			target: new(struct{ Name string }),
			want:   struct{ Name string }{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshalGQLJSON(tt.input, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			// Dereferenzieren, falls Pointer
			switch v := tt.target.(type) {
			case *string:
				assert.Equal(t, tt.want, *v)
			case *struct{ Name string }:
				assert.Equal(t, tt.want, *v)
			default:
				assert.Equal(t, tt.want, v)
			}
		})
	}
}
