// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package fga provides a client for interacting with OpenFGA (Fine Grained Authorization).
// It offers a type-safe and idiomatic Go interface for managing authorization tuples,
// checking permissions, and handling authorization-related operations.
package fga

import (
	"context"
	"fmt"
	"strings"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// ListTuplesRequest represents a request to list tuples from the FGA service.
// It allows filtering tuples by user, relation, and object.
type ListTuplesRequest struct {
	// Subject filters tuples by subject identifier
	Subject Entity
	// Relation filters tuples by relation type
	Relation Relation
	// Object filters tuples by object identifier
	Object Entity
}

// ListTuplesResponse represents the response from a list tuples operation.
// It contains the list of tuples matching the request criteria.
type ListTuplesResponse struct {
	// Tuples is the list of matching tuples
	Tuples []TupleKey
}

// ListTuples retrieves a list of tuples from the FGA service based on the provided filters.
// Returns a list of tuples matching the request criteria.
//
// Example:
//
//	tuples, err := client.ListTuples(ctx, ListTuplesRequest{
//	    Subject: Entity{Kind: "user", Identifier: "123"},
//	    Relation: "viewer",
//	    Object: Entity{Kind: "document", Identifier: "456"},
//	})
func (c *Client) ListTuples(ctx context.Context, req ListTuplesRequest) (*ListTuplesResponse, error) {
	if err := c.validateListTuplesRequest(req); err != nil {
		return nil, fmt.Errorf("invalid list tuples request: %w", err)
	}

	subjectStr := req.Subject.String()
	relationStr := req.Relation.String()
	objectStr := req.Object.String()

	body := client.ClientReadRequest{
		User:     &subjectStr,
		Relation: &relationStr,
		Object:   &objectStr,
	}

	resp, err := c.client.Read(ctx).Body(body).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Interface("request", req).
			Str("subject", subjectStr).
			Str("relation", relationStr).
			Str("object", objectStr).
			Msg("failed to list tuples")

		return nil, fmt.Errorf("failed to list tuples: %w", err)
	}

	return &ListTuplesResponse{
		Tuples: convertToTuples(resp.Tuples),
	}, nil
}

// validateListTuplesRequest validates the list tuples request parameters.
// Returns an error if the request is invalid.
func (c *Client) validateListTuplesRequest(req ListTuplesRequest) error {
	// At least one filter must be provided
	if req.Subject.Identifier == "" && req.Relation == "" && req.Object.Identifier == "" {
		return fmt.Errorf("%w: at least one filter (subject, relation, or object) must be provided", ErrInvalidArgument)
	}

	return nil
}

// convertToTuples converts a slice of openfga.Tuple to a slice of TupleKey.
// This function parses the OpenFGA tuple format into our internal TupleKey structure.
// It uses ParseFGATupleKey to handle the conversion of individual tuples.
func convertToTuples(clientTuples []openfga.Tuple) []TupleKey {
	tuples := make([]TupleKey, len(clientTuples))

	for i, t := range clientTuples {
		tupleKey := ParseFGATupleKey(t.Key)
		if tupleKey != nil {
			tuples[i] = *tupleKey
		}
	}

	return tuples
}

// WriteTupleKeys writes and/or deletes multiple tuples in a single operation.
// Returns the response from the FGA service and any error that occurred.
//
// Example:
//
//	resp, err := client.WriteTupleKeys(ctx, []TupleKey{
//	    {Subject: user, Relation: "member", Object: org},
//	}, []TupleKey{
//	    {Subject: user, Relation: "viewer", Object: doc},
//	})
func (c *Client) WriteTupleKeys(ctx context.Context, writes []TupleKey, deletes []TupleKey) (*client.ClientWriteResponse, error) {
	opts := client.ClientWriteOptions{}

	body := client.ClientWriteRequest{
		Writes:  tupleKeyToWriteRequest(writes),
		Deletes: tupleKeyToDeleteRequest(deletes),
	}

	return c.handleWrite(c.client.
		Write(ctx).
		Body(body).
		Options(opts).
		Execute(),
	)
}

// handleWrite processes the response from a write operation.
// It handles duplicate key errors based on the client's configuration.
// If IgnoreDuplicateKeyError is true, duplicate key errors are logged and skipped.
// Otherwise, all errors are returned to the caller.
func (c *Client) handleWrite(resp *client.ClientWriteResponse, err error) (*client.ClientWriteResponse, error) {
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, ErrEmptyResponse
	}

	// Avoid any tuple parsing if we don't allow duplicates
	if !c.IgnoreDuplicateKeyError {
		log.Info().
			Err(err).
			Interface("writes", resp.Writes).
			Interface("deletes", resp.Deletes).
			Msg("error writing tuples")

		return resp, nil
	}

	for _, entry := range collectWriteResults(resp) {
		if entry.Error == nil {
			continue
		}

		ll := log.With().
			Err(entry.Error).
			Str("user", entry.User).
			Str("relation", entry.Relation).
			Str("object", entry.Object).
			Str("operation", entry.Operation).
			Logger()

		// avoid string allocations by using HasPrefix directly on error messages (lightweight)
		msg := entry.Error.Error()
		if strings.HasPrefix(msg, ErrDuplicateKey) || strings.Contains(msg, ErrDuplicateKey) {
			ll.Warn().Msg("duplicate relation, skipping")
			continue
		}

		ll.Error().Msg("error writing tuple")

		return nil, &WriteError{
			User:          entry.User,
			Relation:      entry.Relation,
			Object:        entry.Object,
			Operation:     entry.Operation,
			ErrorResponse: entry.Error,
		}
	}

	return resp, nil
}

// tupleResult represents the result of a single tuple operation (write or delete).
// It contains information about the operation and any error that occurred.
type tupleResult struct {
	// User is the user identifier involved in the operation
	User string
	// Relation is the relation that was modified
	Relation string
	// Object is the object identifier involved in the operation
	Object string
	// Error contains any error that occurred during the operation
	Error error
	// Operation indicates whether this was a write or delete operation
	Operation string
}

// collectWriteResults combines the results of write and delete operations into a single slice.
// It processes both successful and failed operations, maintaining the order of operations.
func collectWriteResults(resp *client.ClientWriteResponse) []tupleResult {
	n := len(resp.Writes) + len(resp.Deletes)
	out := make([]tupleResult, 0, n)

	for i := range resp.Writes {
		t := resp.Writes[i]
		out = append(out, tupleResult{
			User:      t.TupleKey.User,
			Relation:  t.TupleKey.Relation,
			Object:    t.TupleKey.Object,
			Error:     t.Error,
			Operation: OpWrite,
		})
	}

	for i := range resp.Deletes {
		t := resp.Deletes[i]
		out = append(out, tupleResult{
			User:      t.TupleKey.User,
			Relation:  t.TupleKey.Relation,
			Object:    t.TupleKey.Object,
			Error:     t.Error,
			Operation: OpDelete,
		})
	}

	return out
}

// tupleKeyToWriteRequest converts a slice of TupleKey to a slice of client.ClientTupleKey.
// This is used for write operations. It includes any conditions specified in the TupleKey.
func tupleKeyToWriteRequest(tupleKeys []TupleKey) []client.ClientTupleKey {
	req := make([]client.ClientTupleKey, len(tupleKeys))
	for i := range tupleKeys {
		req[i] = client.ClientTupleKey{
			User:      tupleKeys[i].Subject.String(),
			Relation:  tupleKeys[i].Relation.String(),
			Object:    tupleKeys[i].Object.String(),
			Condition: tupleKeys[i].Condition.toOpenFgaCondition(),
		}
	}

	return req
}

// tupleKeyToDeleteRequest converts a slice of TupleKey to a slice of openfga.TupleKeyWithoutCondition.
// This is used for delete operations. Conditions are not included in delete operations.
func tupleKeyToDeleteRequest(tupleKeys []TupleKey) []openfga.TupleKeyWithoutCondition {
	req := make([]openfga.TupleKeyWithoutCondition, len(tupleKeys))
	for i := range tupleKeys {
		req[i] = openfga.TupleKeyWithoutCondition{
			User:     tupleKeys[i].Subject.String(),
			Relation: tupleKeys[i].Relation.String(),
			Object:   tupleKeys[i].Object.String(),
		}
	}

	return req
}

// TupleRequest is the fields needed to check a tuple in the FGA store
type TupleRequest struct {
	// ObjectID is the identifier of the object that the subject is related to
	ObjectID string
	// ObjectType is the type of object that the subject is related to
	ObjectType string
	// ObjectRelation is the tuple set relation for the object (e.g #member)
	ObjectRelation string
	// SubjectID is the identifier of the subject that is related to the object
	SubjectID string
	// SubjectType is the type of subject that is related to the object
	SubjectType string
	// SubjectRelation is the tuple set relation for the subject (e.g #member)
	SubjectRelation string
	// Relation is the relationship between the subject and object
	Relation string
	// ConditionName for the relationship
	ConditionName string
	// ConditionContext for the relationship
	ConditionContext *map[string]any
}

// WithSubject sets the subject for the tuple request.
func (r *TupleRequest) WithSubjectType(subjectType string) *TupleRequest {
	r.SubjectType = subjectType
	return r
}

// GetTupleKey creates a Tuple key with the provided subject, object, and role
func GetTupleKey(req TupleRequest) TupleKey {
	sub := Entity{
		Kind:       Kind(req.SubjectType),
		Identifier: req.SubjectID,
	}

	if req.SubjectRelation != "" {
		sub.Relation = Relation(req.SubjectRelation)
	}

	object := Entity{
		Kind:       Kind(req.ObjectType),
		Identifier: req.ObjectID,
	}

	if req.ObjectRelation != "" {
		object.Relation = Relation(req.ObjectRelation)
	}

	k := TupleKey{
		Subject:  sub,
		Object:   object,
		Relation: Relation(req.Relation),
	}

	if req.ConditionName != "" {
		k.Condition = Condition{
			Name:    req.ConditionName,
			Context: req.ConditionContext,
		}
	}

	return k
}

// CreatePublicWildcardTuples creates a slice of TupleKey for public access with the specified relation.
// This is an internal method that requires all parameters.
//
// Args:
//   - relation: The relation to use for the tuples
//   - objectType: The type of object to create tuples for
//   - objectID: The ID of the object
//
// Returns:
//   - []TupleKey: A slice of TupleKey with wildcard subject and the specified relation and object
func CreatePublicWildcardTuples(relation Relation, objectType string, objectID string) []TupleKey {
	userTuple := &TupleRequest{
		ObjectID:    objectID,
		ObjectType:  objectType,
		SubjectID:   Wildcard,
		Relation:    string(relation),
		SubjectType: userSubject,
	}
	serviceTuple := &TupleRequest{
		ObjectID:    objectID,
		ObjectType:  objectType,
		SubjectID:   Wildcard,
		Relation:    string(relation),
		SubjectType: serviceSubject,
	}

	return []TupleKey{
		GetTupleKey(*userTuple),
		GetTupleKey(*serviceTuple),
	}
}

// CreatePublicViewTuples creates a slice of TupleKey for public view access.
// This is a convenience method that uses "canView" as the default relation.
//
// Args:
//   - objectType: The type of object to create tuples for
//   - objectID: The ID of the object
//
// Returns:
//   - []TupleKey: A slice of TupleKey with wildcard subject and "canView" relation
func CreatePublicViewTuples(objectType string, objectID string) []TupleKey {
	return CreatePublicWildcardTuples(CanView, objectType, objectID)
}
