// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package krn implements the Kopexa Resource Name (KRN) system, following Google's resource naming design.
// KRN provides a standardized way to identify and reference resources across the Kopexa platform.
//
// Key features:
//   - Canonical resource naming format: //<service-name>/<relative-resource-name>
//   - Support for JSON and YAML serialization
//   - Database integration via sql.Scanner and driver.Valuer
//   - Legacy format support
//   - Resource ID validation
//
// Example usage:
//
//	krn, err := krn.New("//kopexa.com/frameworks/iso-27001-2022")
//	if err != nil {
//	    // handle error
//	}
package krn

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

// Common errors
var (
	ErrInvalidKRNFormat         = errors.New("invalid KRN format")
	ErrInvalidResourceID        = errors.New("invalid resource ID")
	ErrInvalidCollectionID      = errors.New("invalid collection ID")
	ErrCollectionNotFound       = errors.New("collection not found")
	ErrUnsupportedType          = errors.New("unsupported type")
	ErrMissingResourcePath      = errors.New("missing relative resource path")
	ErrInvalidLegacyFormat      = errors.New("invalid legacy KRN format")
	ErrMustStartWithDoubleSlash = errors.New("KRN must start with //")
)

const (
	// PathSeparator is the character used to separate path components
	PathSeparator = "/"
	// MinPathComponents is the minimum number of components in a valid KRN path
	MinPathComponents = 2
)

// ServiceID extracts the service identifier from a full service name by removing the base domain.
// Example: ServiceID("service.example.com", ".example.com") returns "service"
func ServiceID(serviceName string, baseDomain string) string {
	return strings.TrimSuffix(serviceName, baseDomain)
}

// IsValid checks if a string is a valid KRN format.
// It verifies that the string can be parsed as a URL and has no scheme, fragment, or query parameters.
func IsValid(krn string) bool {
	x, err := url.Parse(krn)
	if err != nil {
		return false
	}

	return x.Scheme == "" && x.Fragment == "" && x.RawQuery == ""
}

// KRN represents a Kopexa Resource Name following Google's resource naming design.
// See https://cloud.google.com/apis/design/resource_names for more details.
//
// Format: //<service-name>/<relative-resource-name>
// Example: //kopexa.com/frameworks/iso-27001-2022
type KRN struct {
	ServiceName          string // The service identifier (e.g., "kopexa.com")
	RelativeResourceName string // The resource path (e.g., "frameworks/iso-27001-2022")
}

// New creates a new KRN instance from a full resource name string.
// Returns an error if the input string is not a valid KRN format.
func New(fullResourceName string) (*KRN, error) {
	u, err := url.Parse(fullResourceName)
	if err != nil {
		return nil, fmt.Errorf("invalid KRN format: %w", err)
	}

	path := strings.TrimLeft(u.EscapedPath(), "/")

	return &KRN{
		ServiceName:          u.Host,
		RelativeResourceName: path,
	}, nil
}

// MustNew creates a new KRN instance and panics if the input is invalid.
// This should only be used for constants and tests.
func MustNew(fullResourceName string) *KRN {
	krn, err := New(fullResourceName)
	if err != nil {
		panic(err)
	}

	return krn
}

// NewChildKRN creates a new KRN as a child of an existing KRN.
// The resource and resourceID are appended to the parent's path.
// Returns an error if the resourceID is invalid or if the parent KRN is invalid.
func NewChildKRN(ownerKRN string, resource string, resourceID string) (*KRN, error) {
	if !isValidResourceID(resourceID) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResourceID, resourceID)
	}

	krn, err := New(ownerKRN)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKRNFormat, err)
	}

	krn.RelativeResourceName = path.Join(krn.RelativeResourceName, resource, resourceID)

	return krn, nil
}

// String returns the canonical string representation of the KRN.
func (krn *KRN) String() string {
	return "//" + krn.ServiceName + "/" + krn.RelativeResourceName
}

// CollectionName returns the top-level collection name from the resource path.
// Example: for "frameworks/iso-27001", returns "frameworks"
func (krn *KRN) CollectionName() string {
	parts := strings.Split(krn.RelativeResourceName, "/")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

// Parse parses a canonical KRN string into a KRN struct.
// The input must start with "//" and contain a service name and resource path.
func Parse(input string) (KRN, error) {
	if !strings.HasPrefix(input, "//") {
		return KRN{}, ErrMustStartWithDoubleSlash
	}

	trimmed := strings.TrimPrefix(input, "//")

	parts := strings.SplitN(trimmed, PathSeparator, MinPathComponents)
	if len(parts) != MinPathComponents {
		return KRN{}, ErrMissingResourcePath
	}

	return KRN{
		ServiceName:          parts[0],
		RelativeResourceName: parts[1],
	}, nil
}

// MustParse parses a KRN string and panics if the input is invalid.
// This should only be used for constants and tests.
func MustParse(input string) KRN {
	krn, err := Parse(input)
	if err != nil {
		panic(err)
	}

	return krn
}

// ParseLegacy parses a legacy KRN format that may not include the "//" prefix.
// If the input starts with "//", it is treated as a canonical KRN.
func ParseLegacy(input string) (KRN, error) {
	if strings.HasPrefix(input, "//") {
		return Parse(input)
	}

	parts := strings.SplitN(input, PathSeparator, MinPathComponents)
	if len(parts) != MinPathComponents {
		return KRN{}, ErrInvalidLegacyFormat
	}

	return KRN{
		ServiceName:          parts[0],
		RelativeResourceName: parts[1],
	}, nil
}

// IsZero returns true if the KRN is empty (has no service name or resource path).
func (krn KRN) IsZero() bool {
	return krn.ServiceName == "" || krn.RelativeResourceName == ""
}

// Basename returns the last component of the resource path.
// Example: for "frameworks/iso-27001", returns "iso-27001"
func (krn *KRN) Basename() string {
	keyValues := strings.Split(krn.RelativeResourceName, "/")
	if len(keyValues) == 0 {
		return ""
	}

	return keyValues[len(keyValues)-1]
}

// ResourceID returns the value associated with a collection ID in the resource path.
// The collection ID and its value must be adjacent in the path.
// Returns an error if the collection ID is not found or if the path is malformed.
func (krn *KRN) ResourceID(collectionID string) (string, error) {
	keyValues := strings.Split(krn.RelativeResourceName, PathSeparator)
	for i := 0; i < len(keyValues); i += 2 {
		if keyValues[i] == collectionID {
			if i+1 < len(keyValues) {
				return keyValues[i+1], nil
			}

			return "", fmt.Errorf("%w: %s", ErrInvalidCollectionID, krn.String())
		}
	}

	return "", fmt.Errorf("%w: %s", ErrCollectionNotFound, collectionID)
}

// Equals compares the KRN with another resource string for equality.
// Returns true if both represent the same resource.
func (krn *KRN) Equals(resource string) bool {
	parsed, err := New(resource)
	if err != nil {
		return false
	}

	return parsed.ServiceName == krn.ServiceName && parsed.RelativeResourceName == krn.RelativeResourceName
}

// reResourceID defines the pattern for valid resource IDs:
// - 4-200 characters long
// - Contains only lowercase letters, digits, dots, or hyphens
// - Example: "1.1.2-tmp-configured"
var reResourceID = regexp.MustCompile(`^([\d-_\.]|[a-zA-Z]){4,200}$`)

// isValidResourceID checks if a string matches the resource ID pattern.
func isValidResourceID(id string) bool {
	return reResourceID.MatchString(id)
}

// MarshalJSON implements the json.Marshaler interface.
func (krn KRN) MarshalJSON() ([]byte, error) {
	return json.Marshal(krn.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (krn *KRN) UnmarshalJSON(data []byte) error {
	var krnStr string
	if err := json.Unmarshal(data, &krnStr); err != nil {
		return fmt.Errorf("KRN must be a string: %w", err)
	}

	parsed, err := Parse(krnStr)
	if err != nil {
		return fmt.Errorf("invalid KRN format: %w", err)
	}

	*krn = parsed

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (krn *KRN) UnmarshalYAML(data []byte) error {
	var krnStr string
	if err := yaml.Unmarshal(data, &krnStr); err != nil {
		return fmt.Errorf("KRN must be a string in YAML: %w", err)
	}

	parsed, err := Parse(krnStr)
	if err != nil {
		return fmt.Errorf("invalid KRN format in YAML: %w", err)
	}

	*krn = parsed

	return nil
}

// Scan implements the sql.Scanner interface for database integration.
func (krn *KRN) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedType, value)
	}

	// Handle legacy KRN without leading //
	if !strings.HasPrefix(str, "//") {
		str = "//" + str
	}

	parsed, err := Parse(str)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidKRNFormat, err)
	}

	*krn = parsed

	return nil
}

// Value implements the driver.Valuer interface for database integration.
func (krn KRN) Value() (driver.Value, error) {
	return krn.String(), nil
}
