// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/openfga/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// AccessCheck represents a permission check request.
// It contains all the necessary information to check if a subject has a specific relation to an object.
// This struct is used to define the parameters for access control checks in the FGA system.
type AccessCheck struct {
	// SubjectType is the type of the subject (e.g., "user", "organization")
	// If not specified, the defaultSubject ("user") will be used.
	SubjectType string
	// SubjectID is the unique identifier of the subject
	// This is required and must not be empty.
	SubjectID string
	// ObjectID is the unique identifier of the object
	// This is required and must not be empty.
	ObjectID string
	// ObjectType is the type of the object (e.g., "document", "space")
	// This is required and must not be empty.
	ObjectType string
	// Relation is the relation to check (e.g., "viewer", "editor")
	// This is required and must not be empty.
	Relation string
	// Context is the context of the request used for conditional relationships
	// This is optional and can be nil.
	Context          *map[string]any
	ContextualTuples []ContextualTupleKey
}

// validate ensures that all required fields are present in the AccessCheck.
// It checks that SubjectID, ObjectID, and Relation are not empty.
//
// Returns:
//   - error: If any required field is missing, returns ErrInvalidArgument
func (ac AccessCheck) validate() error {
	if ac.SubjectID == "" || ac.ObjectID == "" || ac.Relation == "" {
		return fmt.Errorf("%w: subject_id, object_id, and relation are required", ErrInvalidArgument)
	}

	return nil
}

// toBatchCheckItem converts the AccessCheck to a client.ClientBatchCheckItem.
// It validates the AccessCheck and creates a batch check item with a unique correlation ID.
//
// Returns:
//   - *client.ClientBatchCheckItem: The converted batch check item
//   - error: If validation fails
func (ac AccessCheck) toBatchCheckItem() (*client.ClientBatchCheckItem, error) {
	if err := ac.validate(); err != nil {
		log.Error().Err(err).Msg("failed to validate access check")
		return nil, err
	}

	if ac.SubjectType == "" {
		ac.SubjectType = defaultSubject
	}

	sub := Entity{
		Kind:       Kind(ac.SubjectType),
		Identifier: ac.SubjectID,
	}

	obj := Entity{
		Kind:       Kind(ac.ObjectType),
		Identifier: ac.ObjectID,
	}

	return &client.ClientBatchCheckItem{
		User:             sub.String(),
		Relation:         ac.Relation,
		Object:           obj.String(),
		Context:          ac.Context,
		CorrelationId:    ulid.Make().String(), // generate a new correlation ID for each check
		ContextualTuples: ac.ContextualTuples,
	}, nil
}

// toCheckRequest converts the AccessCheck struct to a client.ClientCheckRequest.
// It validates the AccessCheck and creates a check request with the appropriate subject and object.
//
// Returns:
//   - *client.ClientCheckRequest: The converted check request
//   - error: If validation fails
func (ac AccessCheck) toCheckRequest() (*client.ClientCheckRequest, error) {
	if err := ac.validate(); err != nil {
		log.Error().Err(err).Msg("failed to validate access check")
		return nil, err
	}

	sub := Entity{
		Kind:       Kind(ac.SubjectType),
		Identifier: ac.SubjectID,
	}

	obj := Entity{
		Kind:       Kind(ac.ObjectType),
		Identifier: ac.ObjectID,
	}

	return &client.ClientCheckRequest{
		User:             sub.String(),
		Relation:         ac.Relation,
		Object:           obj.String(),
		Context:          ac.Context,
		ContextualTuples: ac.ContextualTuples,
	}, nil
}

// CheckAccess checks if a subject has a specific relation to an object.
// This is the main method used for access control checks in the application.
//
// Example:
//
//	allowed, err := client.CheckAccess(ctx, AccessCheck{
//	    SubjectID: "user123",
//	    Relation: "viewer",
//	    ObjectType: "document",
//	    ObjectID: "doc456",
//	})
//
// Parameters:
//   - ctx: The context for the request
//   - ac: The AccessCheck containing the permission check parameters
//
// Returns:
//   - bool: True if the permission is granted, false otherwise
//   - error: If the check fails
func (c *Client) CheckAccess(ctx context.Context, ac AccessCheck) (bool, error) {
	return c.checkAccess(ctx, ac)
}

// checkAccess performs the actual permission check.
// It converts the AccessCheck to a request and sends it to the FGA service.
//
// Parameters:
//   - ctx: The context for the request
//   - ac: The AccessCheck containing the permission check parameters
//
// Returns:
//   - bool: True if the permission is granted, false otherwise
//   - error: If the check fails
func (c *Client) checkAccess(ctx context.Context, ac AccessCheck) (bool, error) {
	request, err := ac.toCheckRequest()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert access check to request")
		return false, err
	}

	return c.checkTuple(ctx, *request)
}

// checkTuple sends a check request to the FGA service and returns the result.
// This is the low-level method that actually communicates with the FGA service.
//
// Parameters:
//   - ctx: The context for the request
//   - body: The check request to send
//
// Returns:
//   - bool: True if the permission is granted, false otherwise
//   - error: If the check fails
func (c *Client) checkTuple(ctx context.Context, body client.ClientCheckRequest) (bool, error) {
	data, err := c.client.Check(ctx).Body(body).Execute()
	if err != nil {
		log.Error().Err(err).Interface("tuple", body).Msg("failed to check tuple")
		return false, err
	}

	return data.GetAllowed(), nil
}

// BatchCheckObjectAccess performs multiple access checks in a single request.
// This is more efficient than making multiple individual CheckAccess calls.
//
// Parameters:
//   - ctx: The context for the request
//   - checks: A slice of AccessCheck structs to check
//
// Returns:
//   - []string: A list of object IDs that the subject has access to
//   - error: If any of the checks fail
func (c *Client) BatchCheckObjectAccess(ctx context.Context, checks []AccessCheck) ([]string, error) {
	if len(checks) == 0 {
		return []string{}, nil
	}

	checkRequests := make([]client.ClientBatchCheckItem, 0, len(checks))

	for _, check := range checks {
		item, err := check.toBatchCheckItem()
		if err != nil {
			return nil, err
		}

		checkRequests = append(checkRequests, *item)
	}

	results, err := c.client.BatchCheck(ctx).Body(
		client.ClientBatchCheckRequest{
			Checks: checkRequests,
		},
	).Execute()
	if err != nil {
		return nil, err
	}

	allowedObjects := make([]string, 0, len(checks))

	for id, result := range *results.Result {
		if result.HasError() {
			err := result.GetError()
			log.Error().Str("error", err.GetMessage()).Interface("accessCheck", id).Msg("error checking access")

			continue
		}

		if result.GetAllowed() {
			// get id from the correlation ID
			check, ok := getCheckItemByCorrelationID(id, checkRequests)
			if !ok {
				log.Error().Str("correlationID", id).Msg("correlation ID not found in checks")

				continue
			}

			obj, err := ParseEntity(check.Object)
			if err != nil {
				log.Error().Err(err).Str("object", check.Object).Msg("error parsing object")

				return nil, err
			}

			allowedObjects = append(allowedObjects, obj.Identifier)
		}
	}

	return allowedObjects, nil
}

// getCheckItemByCorrelationID retrieves the check by correlation ID from the list of checks.
// This is a helper function used in batch checks to match results with their original requests.
//
// Parameters:
//   - correlationID: The correlation ID to look up
//   - checks: The list of check requests to search through
//
// Returns:
//   - client.ClientBatchCheckItem: The matching check item
//   - bool: True if a match was found, false otherwise
func getCheckItemByCorrelationID(correlationID string, checks []client.ClientBatchCheckItem) (client.ClientBatchCheckItem, bool) {
	for _, check := range checks {
		if check.CorrelationId == correlationID {
			return check, true
		}
	}

	return client.ClientBatchCheckItem{}, false
}
