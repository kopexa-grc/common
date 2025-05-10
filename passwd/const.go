// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package passwd

var commonPasswords = map[string]struct{}{
	"password":         {},
	"password123":      {},
	"p@ssw0rd123":      {},
	"123456":           {},
	"123456789":        {},
	"qwerty":           {},
	"12345678":         {},
	"1234567890":       {},
	"1234567890123456": {},
	"iloveyou":         {},
	"admin":            {},
	"welcome":          {},
	"monkey":           {},
	"abc123":           {},
	"letmein":          {},
	"football":         {},
	"baseball":         {},
	"dragon":           {},
	"sunshine":         {},
	"princess":         {},
	"trustno1":         {},
	"superman":         {},
	"qazwsx":           {},
	"1qaz2wsx":         {},
}

var leetMap = map[rune]rune{
	'4': 'a', '@': 'a',
	'3': 'e',
	'1': 'l', '!': 'i', '|': 'i',
	'0': 'o',
	'$': 's', '5': 's',
	'7': 't',
	'2': 'z',
	'9': 'g',
}

const (
	fmtPasswordTooShort             = "Password is too short (min %d characters)"
	fmtPasswordTooCommon            = "Password is too common or easily guessed"
	fmtPasswordContainsPersonalInfo = "Password contains personal information"
	fmtPasswordIdealLength          = "Increase length to at least %d characters"
	//nolint:gosec // This is a user-facing feedback string, not a credential
	fmtPasswordTooFewCharacterTypes = "Use both upper and lowercase letters"
	//nolint:gosec // This is a user-facing feedback string, not a credential
	fmtPasswordTooFewNumbers = "Add at least one number"
	//nolint:gosec // This is a user-facing feedback string, not a credential
	fmtPasswordTooFewSymbols = "Add a special character (e.g. !, $, #)"
)
