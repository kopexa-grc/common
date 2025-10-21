// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga

import (
	"fmt"
	"regexp"
	"strings"

	openfga "github.com/openfga/go-sdk"
)

// Kind represents a type of entity in the FGA system.
// It is used to categorize entities and enforce type safety.
// Examples: "user", "organization", "document"
type Kind string

// String returns the lowercase string representation of the Kind.
// This ensures consistent formatting when the kind is used in FGA tuples.
func (k Kind) String() string {
	return strings.ToLower(string(k))
}

// Relation represents a permission or capability in the FGA system.
// It is used to define the type of access or relationship between entities.
// Examples: "member", "owner", "viewer", "editor"
type Relation string

// String returns the lowercase string representation of the Relation.
// This ensures consistent formatting when the relation is used in FGA tuples.
func (r Relation) String() string {
	return strings.ToLower(string(r))
}

// Entity represents an entity in the FGA system.
// It combines a Kind (type) and an Identifier to form a unique entity reference.
// The Relation field is optional and is used when the entity is part of a relationship.
//
// Examples:
// - user:<user_id>
// - organization:<organization_id>#<relation>
// - space:<space_id>#<relation>
type Entity struct {
	// Kind represents the type of the entity (e.g., "user", "organization", "space")
	Kind Kind
	// Identifier is the unique identifier for the entity
	Identifier string
	// Relation is an optional relation that this entity has in a specific context
	Relation Relation
}

// String returns the string representation of the Entity.
// If no Relation is specified, returns "<kind>:<identifier>"
// If a Relation is specified, returns "<kind>:<identifier>#<relation>"
func (e Entity) String() string {
	if e.Kind == "*" && e.Identifier == "*" {
		return "*"
	}

	if e.Kind == "" && e.Identifier == "" {
		return ""
	}

	if e.Relation == "" {
		return fmt.Sprintf("%s:%s", e.Kind, e.Identifier)
	}

	return fmt.Sprintf("%s:%s#%s", e.Kind, e.Identifier, e.Relation)
}

// Condition represents a condition that must be met for a permission to be granted.
// It consists of a condition name and an optional context map containing condition parameters.
type Condition struct {
	// Name is the identifier of the condition to be evaluated
	Name string
	// Context is an optional map of parameters that will be passed to the condition evaluator
	Context *map[string]any
}

// toOpenFgaCondition converts the Condition to an OpenFGA RelationshipCondition.
// Returns nil if no condition name is specified.
func (c Condition) toOpenFgaCondition() *openfga.RelationshipCondition {
	if c.Name == "" {
		return nil
	}

	return &openfga.RelationshipCondition{
		Name:    c.Name,
		Context: c.Context,
	}
}

// TupleKey represents a complete permission tuple in the FGA system.
// It defines who (Subject) has what permission (Relation) on what (Object),
// optionally with additional conditions that must be met.
type TupleKey struct {
	// Subject is the entity that is being granted the permission
	// Example: user:<user_id>
	Subject Entity

	// Object is the entity that is being granted the permission
	// Example: organization:<organization_id>
	Object Entity

	// Relation is the relation that is being granted to the subject
	// Example: member
	Relation Relation

	// Conditions are optional conditions that must be met for the permission to be granted
	Condition Condition
}

// entityRegex is a regular expression for validating entity strings.
// It matches strings in the format "<kind>:<identifier>#<relation>?" where:
// - kind: alphanumeric with underscores and hyphens
// - identifier: alphanumeric with underscores, hyphens, @, ., +, -
// - relation: optional, alphanumeric with underscores and hyphens
var entityRegex = regexp.MustCompile(`([A-za-z0-9_][A-za-z0-9_-]*):([A-za-z0-9_][A-za-z0-9_@.+-]*)(#([A-za-z0-9_][A-za-z0-9_-]*))?`)

// ParseEntity parses a string representation of an entity into an Entity struct.
// The input string must be in the format "<kind>:<identifier>#<relation>?".
// Returns an error if the string format is invalid.
//
// Example:
//
//	entity, err := ParseEntity("user:123#member")
//	if err != nil {
//		// handle error
//	}
func ParseEntity(key string) (Entity, error) {
	c := strings.Count(key, ":")
	if c != 1 {
		return Entity{}, fmt.Errorf("%w: %s", ErrInvalidEntity, key)
	}

	match := entityRegex.FindStringSubmatch(key)
	if match == nil {
		return Entity{}, fmt.Errorf("%w: %s", ErrInvalidEntity, key)
	}

	return Entity{
		Kind:       Kind(match[1]),
		Identifier: match[2],
		Relation:   Relation(match[4]),
	}, nil
}

// parseFGATupleKey converts an OpenFGA TupleKey to our internal TupleKey representation.
// It parses the user and object strings into Entity structs using ParseEntity.
// Returns nil if either the user or object string cannot be parsed.
//
// Example:
//
//	tupleKey := parseFGATupleKey(openfga.TupleKey{
//		User:     "user:123",
//		Relation: "member",
//		Object:   "organization:456",
//	})
func ParseFGATupleKey(t openfga.TupleKey) *TupleKey {
	subject, err := ParseEntity(t.User)
	if err != nil {
		return nil
	}

	object, err := ParseEntity(t.Object)
	if err != nil {
		return nil
	}

	return &TupleKey{
		Subject:  subject,
		Object:   object,
		Relation: Relation(t.Relation),
	}
}
