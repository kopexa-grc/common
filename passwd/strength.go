// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package passwd

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

type StrengthLevel uint8

const (
	Rejected StrengthLevel = iota // common, username match, obvious
	TooShort                      // < 16 chars
	Low                           // easy to guess (e.g. only digits)
	Medium                        // some complexity, but predictable
	High                          // >=16 chars, mixed types
	VeryHigh                      // >=24 chars, high entropy, no patterns
)

const (
	requiredMinLength = 8
	minLength         = 16
	idealMinLength    = 24
	minEntropy        = 4.0
)

type Feedback struct {
	Level    StrengthLevel
	Messages []string
}

// Evaluate is a shorthand for EvaluateWithContext using empty username/email/org
func Evaluate(pw string) Feedback {
	return EvaluateWithContext(pw, "", "", "")
}

// EvaluateWithContext evaluates password strength using additional user context
func EvaluateWithContext(pw, username, email, org string) Feedback {
	var messages []string

	pw = strings.TrimSpace(pw)
	pwLower := strings.ToLower(pw)

	if len(pw) < requiredMinLength {
		return Feedback{Rejected, []string{fmt.Sprintf(fmtPasswordTooShort, requiredMinLength)}}
	}

	if isInvalid(pwLower) {
		return Feedback{Rejected, []string{fmtPasswordTooCommon}}
	}

	if (username != "" && strings.Contains(pwLower, strings.ToLower(username))) ||
		(email != "" && strings.Contains(pwLower, strings.ToLower(email))) ||
		(org != "" && strings.Contains(pwLower, strings.ToLower(org))) {
		return Feedback{Rejected, []string{fmtPasswordContainsPersonalInfo}}
	}

	if len(pw) < minLength {
		messages = append(messages, fmt.Sprintf(fmtPasswordIdealLength, minLength))
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	score := 0
	if hasUpper && hasLower {
		score++
	} else {
		messages = append(messages, fmtPasswordTooFewCharacterTypes)
	}
	if hasDigit {
		score++
	} else {
		messages = append(messages, fmtPasswordTooFewNumbers)
	}
	if hasSymbol {
		score++
	} else {
		messages = append(messages, fmtPasswordTooFewSymbols)
	}

	entropy := shannonEntropy(pw)
	if len(pw) >= idealMinLength && score >= 3 && entropy > minEntropy {
		return Feedback{VeryHigh, nil}
	}
	if len(pw) >= minLength && score >= 2 {
		return Feedback{High, messages}
	}
	if score >= 2 {
		return Feedback{Medium, messages}
	}
	return Feedback{Low, messages}
}

// isInvalid checks if the password matches known bad patterns, incl. l33t variants.
func isInvalid(pw string) bool {
	pw = strings.ToLower(strings.TrimSpace(pw))

	if _, ok := commonPasswords[pw]; ok {
		return true
	}

	// Replace common l33t characters and re-check
	var normalized strings.Builder
	for _, r := range pw {
		if repl, ok := leetMap[r]; ok {
			normalized.WriteRune(repl)
		} else {
			normalized.WriteRune(r)
		}
	}

	_, ok := commonPasswords[normalized.String()]
	return ok
}

// shannonEntropy estimates the randomness of a string.
// The result increases with the variety and unpredictability of characters.
// Values above ~4.0 are considered high entropy for human-generated passwords.
func shannonEntropy(s string) float64 {
	freq := make(map[rune]float64)
	for _, r := range s {
		freq[r]++
	}
	var entropy float64
	l := float64(len(s))
	for _, count := range freq {
		p := count / l
		entropy -= p * math.Log2(p)
	}
	return entropy
}
