// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package auth

import (
	"fmt"
	"io"
	"strings"
)

// AuthLevel is a custom type representing the various states of AuthLevel.
type Level string

const (
	// LevelNotAuthenticated indicates the not authenticated.
	LevelNotAuthenticated Level = "NOT_AUTHENTICATED"
	// LevelOneFactor indicates the one factor.
	LevelOneFactor Level = "ONE_FACTOR"
	// LevelTwoFactor indicates the two factor.
	LevelTwoFactor Level = "TWO_FACTOR"
	// LevelInvalid is used when an unknown or unsupported value is provided.
	LevelInvalid Level = "INVALID"
)

// Values returns a slice of strings representing all valid AuthLevel values.
func (Level) Values() []string {
	return []string{
		string(LevelNotAuthenticated),
		string(LevelOneFactor),
		string(LevelTwoFactor),
	}
}

// String returns the string representation of the AuthLevel value.
func (r Level) String() string {
	return string(r)
}

// ToLevel converts a string to its corresponding Level enum value.
func ToLevel(r string) Level {
	switch strings.ToUpper(r) {
	case LevelNotAuthenticated.String():
		return LevelNotAuthenticated
	case LevelOneFactor.String():
		return LevelOneFactor
	case LevelTwoFactor.String():
		return LevelTwoFactor
	default:
		return LevelInvalid
	}
}

// MarshalGQL implements the gqlgen Marshaler interface.
func (r Level) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(`"` + r.String() + `"`))
}

// UnmarshalGQL implements the gqlgen Unmarshaler interface.
func (r *Level) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("wrong type for AuthLevel, got: %T", v) //nolint:err113
	}

	*r = Level(str)

	return nil
}

// ToPtr returns a pointer to the AuthLevel value.
func (r Level) ToPtr() *Level {
	return &r
}
