// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package passwd

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     StrengthLevel
	}{
		{
			name:     "Too short password",
			password: "short",
			want:     Rejected,
		},
		{
			name:     "Common password",
			password: "password123",
			want:     Rejected,
		},
		{
			name:     "Low strength - only digits",
			password: "1234567890123456",
			want:     Rejected,
		},
		{
			name:     "Medium strength - letters and numbers",
			password: "Password123456",
			want:     Medium,
		},
		{
			name:     "High strength - mixed types",
			password: "Password123!@#",
			want:     Medium,
		},
		{
			name:     "Very high strength - long and complex",
			password: "Super53%l;-:s7dUncommonSecure!@#$%^&*()",
			want:     VeryHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Evaluate(tt.password)
			if got.Level != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got.Level, tt.want)
			}
		})
	}
}

func TestEvaluateWithContext(t *testing.T) {
	tests := []struct {
		name     string
		password string
		username string
		email    string
		org      string
		want     StrengthLevel
	}{
		{
			name:     "Password contains username",
			password: "myusername123!@#",
			username: "myusername",
			want:     Rejected,
		},
		{
			name:     "Password contains email",
			password: "test@example.com123",
			email:    "test@example.com",
			want:     Rejected,
		},
		{
			name:     "Password contains organization",
			password: "MyOrg123!@#",
			org:      "MyOrg",
			want:     Rejected,
		},
		{
			name:     "Valid password with context",
			password: "SuperSecurePassword123!@#$%^&*()",
			username: "user123",
			email:    "test@example.com",
			org:      "MyOrg",
			want:     VeryHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateWithContext(tt.password, tt.username, tt.email, tt.org)
			if got.Level != tt.want {
				t.Errorf("EvaluateWithContext() = %v, want %v", got.Level, tt.want)
			}
		})
	}
}

func TestEvaluateWithContext_AllPaths(t *testing.T) {
	tests := []struct {
		name     string
		password string
		username string
		email    string
		org      string
		want     StrengthLevel
	}{
		{
			name:     "Rejected - too short",
			password: "abc",
			want:     Rejected,
		},
		{
			name:     "Rejected - common password",
			password: "password123",
			want:     Rejected,
		},
		{
			name:     "Rejected - contains personal info (username)",
			password: "myuser123!@#",
			username: "myuser",
			want:     Rejected,
		},
		{
			name:     "Low - only lowercase, long genug",
			password: "abcdefghijklmnoq",
			want:     Low,
		},
		{
			name:     "Medium - lower+upper+digit, aber zu kurz für High",
			password: "Abcdef123",
			want:     Medium,
		},
		{
			name:     "High - >=16, mind. 2 Typen",
			password: "Abcdefghijklmn12",
			want:     High,
		},
		{
			name:     "VeryHigh - >=24, 3 Typen, hohe Entropie",
			password: "Abcdefghijklmnop1234!@#$",
			want:     VeryHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateWithContext(tt.password, tt.username, tt.email, tt.org)
			if got.Level != tt.want {
				t.Errorf("EvaluateWithContext() = %v, want %v", got.Level, tt.want)
			}
		})
	}
}

func TestIsInvalid(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Common password",
			password: "password123",
			want:     true,
		},
		{
			name:     "L33t speak password",
			password: "p@ssw0rd123",
			want:     true,
		},
		{
			name:     "Valid password",
			password: "SecurePass123!@#",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isInvalid(tt.password)
			if got != tt.want {
				t.Errorf("isInvalid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShannonEntropy(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     float64
	}{
		{
			name:     "Low entropy - repeated characters",
			password: "aaaaaaaa",
			want:     0.0,
		},
		{
			name:     "Medium entropy - mixed characters",
			password: "Password123",
			want:     3.0,
		},
		{
			name:     "High entropy - complex password",
			password: "P@ssw0rd!@#$%^&*()",
			want:     3.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shannonEntropy(tt.password)
			if got < tt.want {
				t.Errorf("shannonEntropy() = %v, want >= %v", got, tt.want)
			}
		})
	}
}

func TestFeedbackMessages(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantMsgs int
	}{
		{
			name:     "Missing character types",
			password: "password",
			wantMsgs: 1, // Missing numbers, symbols, and length warning
		},
		{
			name:     "Missing symbols",
			password: "Password123",
			wantMsgs: 1, // Missing symbols
		},
		{
			name:     "Missing numbers",
			password: "Password!@#",
			wantMsgs: 2, // Missing numbers und Mindestlänge
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Evaluate(tt.password)
			if len(got.Messages) != tt.wantMsgs {
				t.Errorf("Evaluate() messages count = %v, want %v", len(got.Messages), tt.wantMsgs)
			}
		})
	}
}

func TestMinLengthFeedback(t *testing.T) {
	pw := "XyZabc12!uvw" // 12 Zeichen, garantiert nicht common, keine personal info

	feedback := EvaluateWithContext(pw, "", "", "")
	if feedback.Level == Rejected {
		t.Errorf("Expected not Rejected, got %v", feedback.Level)
	}

	found := false

	for _, msg := range feedback.Messages {
		if msg == "Increase length to at least 16 characters" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected minLength feedback message, got %v", feedback.Messages)
	}
}
