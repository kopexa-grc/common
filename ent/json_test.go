package ent

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalRawMessage(t *testing.T) {
	tests := []struct {
		name    string
		arg     interface{}
		want    json.RawMessage
		wantErr bool
	}{{
		name: "map",
		arg:  map[string]any{"a": true},
		want: json.RawMessage(`{"a":true}`),
	}, {
		name: "array",
		arg:  []int{1, 2},
		want: json.RawMessage(`[1,2]`),
	}, {
		name: "bytes",
		arg:  []byte{'"', 'a', '"'},
		want: json.RawMessage(`"a"`),
	}, {
		// In practice, this is the way graphql Unmarshal is processing input like {json: "a"}:
		name: "string",
		arg:  "a",
		want: json.RawMessage(`"a"`),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalRawMessage(tt.arg)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalRawMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalRawMessage() = %s, want %s", got, tt.want)
			}
		})
	}
}
