// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleEvidence_MarshalGQL(t *testing.T) {
	tests := []struct {
		name     string
		evidence ExampleEvidence
		want     string
	}{
		{
			name: "complete evidence",
			evidence: ExampleEvidence{
				DocumentationType: "policy",
				Description:       "Example policy document",
			},
			want: `{"documentationType":"policy","description":"Example policy document"}`,
		},
		{
			name: "empty evidence",
			evidence: ExampleEvidence{
				DocumentationType: "",
				Description:       "",
			},
			want: `{}`,
		},
		{
			name: "only documentation type",
			evidence: ExampleEvidence{
				DocumentationType: "procedure",
				Description:       "",
			},
			want: `{"documentationType":"procedure"}`,
		},
		{
			name: "only description",
			evidence: ExampleEvidence{
				DocumentationType: "",
				Description:       "Example description",
			},
			want: `{"description":"Example description"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.evidence.MarshalGQL(&buf)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestExampleEvidence_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    ExampleEvidence
		wantErr bool
	}{
		{
			name: "valid json object",
			input: map[string]interface{}{
				"documentationType": "policy",
				"description":       "Example policy document",
			},
			want: ExampleEvidence{
				DocumentationType: "policy",
				Description:       "Example policy document",
			},
			wantErr: false,
		},
		{
			name:    "invalid input type",
			input:   "not a map",
			wantErr: true,
		},
		{
			name:    "empty object",
			input:   map[string]interface{}{},
			want:    ExampleEvidence{},
			wantErr: false,
		},
		{
			name: "partial object",
			input: map[string]interface{}{
				"documentationType": "procedure",
			},
			want: ExampleEvidence{
				DocumentationType: "procedure",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var evidence ExampleEvidence
			err := evidence.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, evidence)
		})
	}
}

func TestExampleEvidence_JSON(t *testing.T) {
	tests := []struct {
		name     string
		evidence ExampleEvidence
		want     string
	}{
		{
			name: "complete evidence",
			evidence: ExampleEvidence{
				DocumentationType: "policy",
				Description:       "Example policy document",
			},
			want: `{"documentationType":"policy","description":"Example policy document"}`,
		},
		{
			name: "empty evidence",
			evidence: ExampleEvidence{
				DocumentationType: "",
				Description:       "",
			},
			want: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.evidence)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(jsonData))

			// Test JSON unmarshaling
			var unmarshaled ExampleEvidence
			err = json.Unmarshal(jsonData, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.evidence, unmarshaled)
		})
	}
}

func TestImplementationGuidance_MarshalGQL(t *testing.T) {
	tests := []struct {
		name     string
		guidance ImplementationGuidance
		want     string
	}{
		{
			name: "complete guidance with KRN",
			guidance: ImplementationGuidance{
				ReferenceID: "//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1",
				Guidance: []string{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			want: `{"referenceId":"//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1","guidance":["Implement access control policy","Review access rights regularly"]}`,
		},
		{
			name: "complete guidance with legacy reference",
			guidance: ImplementationGuidance{
				ReferenceID: "ISO27001-A.5.1.1",
				Guidance: []string{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			want: `{"referenceId":"ISO27001-A.5.1.1","guidance":["Implement access control policy","Review access rights regularly"]}`,
		},
		{
			name: "empty guidance",
			guidance: ImplementationGuidance{
				ReferenceID: "",
				Guidance:    nil,
			},
			want: `{}`,
		},
		{
			name: "only reference ID with KRN",
			guidance: ImplementationGuidance{
				ReferenceID: "//kopexa.com/compliance/nist/800-53/controls/AC-1",
				Guidance:    nil,
			},
			want: `{"referenceId":"//kopexa.com/compliance/nist/800-53/controls/AC-1"}`,
		},
		{
			name: "only guidance",
			guidance: ImplementationGuidance{
				ReferenceID: "",
				Guidance: []string{
					"Step 1",
					"Step 2",
				},
			},
			want: `{"guidance":["Step 1","Step 2"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.guidance.MarshalGQL(&buf)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestImplementationGuidance_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    ImplementationGuidance
		wantErr bool
	}{
		{
			name: "valid json object with KRN",
			input: map[string]interface{}{
				"referenceId": "//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1",
				"guidance": []interface{}{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			want: ImplementationGuidance{
				ReferenceID: "//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1",
				Guidance: []string{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			wantErr: false,
		},
		{
			name: "valid json object with legacy reference",
			input: map[string]interface{}{
				"referenceId": "ISO27001-A.5.1.1",
				"guidance": []interface{}{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			want: ImplementationGuidance{
				ReferenceID: "ISO27001-A.5.1.1",
				Guidance: []string{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid input type",
			input:   "not a map",
			wantErr: true,
		},
		{
			name:    "empty object",
			input:   map[string]interface{}{},
			want:    ImplementationGuidance{},
			wantErr: false,
		},
		{
			name: "partial object with KRN",
			input: map[string]interface{}{
				"referenceId": "//kopexa.com/compliance/nist/800-53/controls/AC-1",
			},
			want: ImplementationGuidance{
				ReferenceID: "//kopexa.com/compliance/nist/800-53/controls/AC-1",
			},
			wantErr: false,
		},
		{
			name: "invalid guidance type",
			input: map[string]interface{}{
				"referenceId": "//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1",
				"guidance":    "not an array",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var guidance ImplementationGuidance
			err := guidance.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, guidance)
		})
	}
}

func TestImplementationGuidance_JSON(t *testing.T) {
	tests := []struct {
		name     string
		guidance ImplementationGuidance
		want     string
	}{
		{
			name: "complete guidance with KRN",
			guidance: ImplementationGuidance{
				ReferenceID: "//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1",
				Guidance: []string{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			want: `{"referenceId":"//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1","guidance":["Implement access control policy","Review access rights regularly"]}`,
		},
		{
			name: "complete guidance with legacy reference",
			guidance: ImplementationGuidance{
				ReferenceID: "ISO27001-A.5.1.1",
				Guidance: []string{
					"Implement access control policy",
					"Review access rights regularly",
				},
			},
			want: `{"referenceId":"ISO27001-A.5.1.1","guidance":["Implement access control policy","Review access rights regularly"]}`,
		},
		{
			name: "empty guidance",
			guidance: ImplementationGuidance{
				ReferenceID: "",
				Guidance:    nil,
			},
			want: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.guidance)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(jsonData))

			// Test JSON unmarshaling
			var unmarshaled ImplementationGuidance
			err = json.Unmarshal(jsonData, &unmarshaled)
			assert.NoError(t, err)

			// Compare fields individually to handle nil vs empty slice
			assert.Equal(t, tt.guidance.ReferenceID, unmarshaled.ReferenceID)
			if tt.guidance.Guidance == nil {
				assert.Nil(t, unmarshaled.Guidance)
			} else {
				assert.Equal(t, tt.guidance.Guidance, unmarshaled.Guidance)
			}
		})
	}
}
