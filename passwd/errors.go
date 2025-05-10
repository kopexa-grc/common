// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package passwd

import "fmt"

// Error types for the derived key algorithm
var (
	ErrCannotCreateDK       = fmt.Errorf("cannot create derived key")
	ErrCouldNotGenerate     = fmt.Errorf("could not generate random salt")
	ErrUnableToVerify       = fmt.Errorf("unable to verify derived key")
	ErrCannotParseDK        = fmt.Errorf("cannot parse derived key")
	ErrCannotParseEncodedEK = fmt.Errorf("cannot parse encoded derived key")
	ErrInvalidArgon2Config  = fmt.Errorf("invalid Argon2Config: all values must be > 0")
)

// newParseError creates a new error for parsing failures
func newParseError(field, value, expected string) error {
	return fmt.Errorf("%w: invalid %s: got %s, expected %s", ErrCannotParseDK, field, value, expected)
}
