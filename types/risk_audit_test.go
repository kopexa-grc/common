// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRiskAudit_String(t *testing.T) {
	tests := []struct {
		name     string
		audit    RiskAudit
		expected string
	}{
		{
			name: "complete audit",
			audit: RiskAudit{
				ScoreComponents: map[string]int{
					"security":   80,
					"compliance": 90,
				},
				CriticalGaps: []string{"gap1", "gap2"},
				UIHints:      []string{"hint1", "hint2"},
			},
			expected: "RiskAudit{ScoreComponents: map[compliance:90 security:80], CriticalGaps: [gap1 gap2], UIHints: [hint1 hint2]}",
		},
		{
			name:     "empty audit",
			audit:    RiskAudit{},
			expected: "RiskAudit{ScoreComponents: map[], CriticalGaps: [], UIHints: []}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.audit.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRiskAudit_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    RiskAudit
		wantErr bool
	}{
		{
			name: "valid json",
			input: map[string]interface{}{
				"scoreComponents": map[string]interface{}{
					"security": 80,
				},
				"criticalGaps": []interface{}{"gap1"},
				"uiHints":      []interface{}{"hint1"},
			},
			want: RiskAudit{
				ScoreComponents: map[string]int{"security": 80},
				CriticalGaps:    []string{"gap1"},
				UIHints:         []string{"hint1"},
			},
			wantErr: false,
		},
		{
			name:    "invalid input",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got RiskAudit

			err := got.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRiskAudit_MarshalGQL(t *testing.T) {
	tests := []struct {
		name     string
		audit    RiskAudit
		expected string
	}{
		{
			name: "complete audit",
			audit: RiskAudit{
				ScoreComponents: map[string]int{"security": 80},
				CriticalGaps:    []string{"gap1"},
				UIHints:         []string{"hint1"},
			},
			expected: `{"scoreComponents":{"security":80},"criticalGaps":["gap1"],"uiHints":["hint1"]}`,
		},
		{
			name:     "empty audit",
			audit:    RiskAudit{},
			expected: `{"scoreComponents":{},"criticalGaps":[],"uiHints":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			tt.audit.MarshalGQL(&buf)
			assert.JSONEq(t, tt.expected, buf.String())
		})
	}
}
