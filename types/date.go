// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"time"
)

// DateTime is a custom GraphQL scalar that converts to/from time.Time.
// It supports both "YYYY-MM-DD" and ISO8601 formats.
//
// Example:
//
//	dt := DateTime(time.Now())
//	str := dt.String() // Returns "2024-03-20T15:04:05Z"
type DateTime time.Time

const (
	// dateLayout represents the simple date format "YYYY-MM-DD"
	dateLayout = "2006-01-02"
	// isoDateLayout represents the full ISO8601 format
	isoDateLayout = time.RFC3339
)

var (
	// ErrUnsupportedDateTimeType is returned when an unsupported time format is provided
	ErrUnsupportedDateTimeType = errors.New("unsupported time format")
	// ErrInvalidTimeType is returned when the date format is invalid
	ErrInvalidTimeType = errors.New("invalid date format, expected YYYY-MM-DD or full ISO8601")
)

// Scan implements the sql.Scanner interface for DateTime.
// It converts a database value into a DateTime.
//
// Parameters:
//   - value: The database value to scan
//
// Returns:
//   - error: If the value cannot be converted to DateTime
func (d *DateTime) Scan(value any) error {
	if value == nil {
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return ErrUnsupportedDateTimeType
	}

	*d = DateTime(t)

	return nil
}

// Value implements the driver.Valuer interface for DateTime.
// It converts a DateTime into a database value.
//
// Returns:
//   - driver.Value: The database value
//   - error: If the conversion fails
func (d DateTime) Value() (driver.Value, error) {
	return time.Time(d), nil
}

// UnmarshalCSV implements the csv.Unmarshaler interface for DateTime.
// It converts a CSV string into a DateTime.
//
// Parameters:
//   - s: The CSV string to unmarshal
//
// Returns:
//   - error: If the string cannot be converted to DateTime
func (d *DateTime) UnmarshalCSV(s string) error {
	if s == "" {
		*d = DateTime{}
		return nil
	}

	if t, err := time.Parse(isoDateLayout, s); err == nil {
		*d = DateTime(t)
		return nil
	}

	if t, err := time.Parse(dateLayout, s); err == nil {
		*d = DateTime(t)
		return nil
	}

	return ErrUnsupportedDateTimeType
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for DateTime.
// It allows DateTime to accept both "YYYY-MM-DD" and "YYYY-MM-DDTHH:MM:SSZ".
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If the value cannot be converted to DateTime
func (d *DateTime) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return ErrUnsupportedDateTimeType
	}

	if str == "" {
		*d = DateTime{}
		return nil
	}

	if t, err := time.Parse(isoDateLayout, str); err == nil {
		*d = DateTime(t)
		return nil
	}

	if t, err := time.Parse(dateLayout, str); err == nil {
		*d = DateTime(t)
		return nil
	}

	return ErrInvalidTimeType
}

// MarshalGQL implements the graphql.Marshaler interface for DateTime.
// It writes the datetime as "YYYY-MM-DDTHH:MM:SSZ".
//
// Parameters:
//   - w: The writer to write the datetime to
func (d DateTime) MarshalGQL(w io.Writer) {
	t := time.Time(d)
	if t.IsZero() {
		_, _ = io.WriteString(w, `""`)
		return
	}

	formatted := fmt.Sprintf("%q", t.Format(isoDateLayout))
	_, _ = io.WriteString(w, formatted)
}

// String returns a human-readable string representation of the DateTime.
// It formats the datetime in ISO8601 format.
//
// Returns:
//   - string: The formatted datetime string
func (d DateTime) String() string {
	t := time.Time(d)
	if t.IsZero() {
		return ""
	}

	return t.Format(isoDateLayout)
}

// ToDateTime converts a string to a DateTime.
// It supports both "YYYY-MM-DD" and ISO8601 formats.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - *DateTime: The converted DateTime
//   - error: If the string cannot be converted
func ToDateTime(s string) (*DateTime, error) {
	if s == "" {
		return nil, ErrInvalidTimeType
	}

	if t, err := time.Parse(isoDateLayout, s); err == nil {
		dt := DateTime(t)
		return &dt, nil
	}

	if t, err := time.Parse(dateLayout, s); err == nil {
		dt := DateTime(t)
		return &dt, nil
	}

	return nil, ErrInvalidTimeType
}
