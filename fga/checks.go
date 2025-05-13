// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"
	"fmt"

	"github.com/openfga/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// AccessCheck represents a permission check request.
// It contains all the necessary information to check if a subject has a specific relation to an object.
type AccessCheck struct {
	// SubjectType is the type of the subject (e.g., "user", "organization")
	SubjectType string
	// SubjectID is the unique identifier of the subject
	SubjectID string
	// ObjectID is the unique identifier of the object
	ObjectID string
	// ObjectType is the type of the object (e.g., "document", "space")
	ObjectType string
	// Relation is the relation to check (e.g., "viewer", "editor")
	Relation string
}

// validate ensures that all required fields are present in the AccessCheck.
// Returns an error if any required field is missing.
func (ac AccessCheck) validate() error {
	if ac.SubjectID == "" || ac.ObjectID == "" || ac.Relation == "" {
		return fmt.Errorf("%w: subject_id, object_id, and relation are required", ErrInvalidArgument)
	}

	return nil
}

// toCheckRequest converts the AccessCheck struct to a client.ClientCheckRequest.
// Returns an error if the AccessCheck is invalid.
func (ac AccessCheck) toCheckRequest() (*client.ClientCheckRequest, error) {
	if err := ac.validate(); err != nil {
		log.Error().Err(err).Msg("failed to validate access check")
		return nil, err
	}

	return &client.ClientCheckRequest{
		User: (&Entity{
			Kind:       Kind(ac.SubjectType),
			Identifier: ac.SubjectID,
		}).String(),
		Relation: ac.Relation,
		Object: (&Entity{
			Kind:       Kind(ac.ObjectType),
			Identifier: ac.ObjectID,
			Relation:   Relation(ac.Relation),
		}).String(),
	}, nil
}

// CheckAccess checks if a subject has a specific relation to an object.
// Returns true if the permission is granted, false otherwise.
// This method is used for access control checks in the application.
//
// Example:
//
//	allowed, err := client.CheckAccess(ctx, AccessCheck{
//	    SubjectID: "user123",
//	    Relation: "viewer",
//	    ObjectType: "document",
//	    ObjectID: "doc456",
//	})
func (c *Client) CheckAccess(ctx context.Context, ac AccessCheck) (bool, error) {
	return c.checkAccess(ctx, ac)
}

// checkAccess performs the actual permission check.
// It converts the AccessCheck to a request and sends it to the FGA service.
func (c *Client) checkAccess(ctx context.Context, ac AccessCheck) (bool, error) {
	request, err := ac.toCheckRequest()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert access check to request")
		return false, err
	}

	return c.checkTuple(ctx, *request)
}

// checkTuple sends a check request to the FGA service and returns the result.
// Returns true if the permission is granted, false otherwise.
func (c *Client) checkTuple(ctx context.Context, body client.ClientCheckRequest) (bool, error) {
	data, err := c.client.Check(ctx).Body(body).Execute()
	if err != nil {
		log.Error().Err(err).Interface("tuple", body).Msg("failed to check tuple")
		return false, err
	}

	return data.GetAllowed(), nil
}
