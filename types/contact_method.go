// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"fmt"
	"io"
	"strings"
)

// ContactMethod represents the method of contact between parties.
// It is used to specify how a contact point should be used for communication.
type ContactMethod string

const (
	// ContactMethodEmail indicates that the contact should be made via email.
	// When this method is used, the contact point must have a valid email address.
	ContactMethodEmail ContactMethod = "EMAIL"

	// ContactMethodPhone indicates that the contact should be made via phone.
	// When this method is used, the contact point must have a valid phone number.
	ContactMethodPhone ContactMethod = "PHONE"

	// ContactMethodWebForm indicates that the contact should be made via a web form.
	// When this method is used, the contact point must have a valid URL.
	ContactMethodWebForm ContactMethod = "WEB_FORM"

	// ContactMethodInvalid is used when an unknown or unsupported value is provided.
	// This value should not be used in normal operation and indicates an error state.
	ContactMethodInvalid ContactMethod = "INVALID"
)

// Values returns a slice of strings representing all valid ContactMethod values.
// This is useful for validation and UI purposes.
//
// Returns:
//   - []string: A slice containing all valid contact method values
func (ContactMethod) Values() []string {
	return []string{
		string(ContactMethodEmail),
		string(ContactMethodPhone),
		string(ContactMethodWebForm),
	}
}

// String returns the string representation of the ContactMethod value.
// This implements the fmt.Stringer interface.
//
// Returns:
//   - string: The string representation of the contact method
func (r ContactMethod) String() string {
	return string(r)
}

// ToContactMethod converts a string to its corresponding ContactMethod enum value.
// The input string is case-insensitive and will be converted to uppercase.
// If the input string does not match any valid contact method, ContactMethodInvalid is returned.
//
// Parameters:
//   - r: The string to convert to a ContactMethod
//
// Returns:
//   - ContactMethod: The corresponding ContactMethod value
func ToContactMethod(r string) ContactMethod {
	switch strings.ToUpper(r) {
	case ContactMethodEmail.String():
		return ContactMethodEmail
	case ContactMethodWebForm.String():
		return ContactMethodWebForm
	case ContactMethodPhone.String():
		return ContactMethodPhone
	default:
		return ContactMethodInvalid
	}
}

// MarshalGQL implements the gqlgen Marshaler interface.
// It allows ContactMethod to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the ContactMethod to
func (r ContactMethod) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(`"` + r.String() + `"`))
}

// UnmarshalGQL implements the gqlgen Unmarshaler interface.
// It allows ContactMethod to be used as a GraphQL scalar type.
// The input value must be a string, otherwise an error is returned.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (r *ContactMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("wrong type for ContactMethod, got: %T", v) //nolint:err113
	}

	*r = ContactMethod(str)

	return nil
}
