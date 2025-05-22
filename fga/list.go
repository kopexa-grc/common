// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/openfga/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// ListRequest represents a request to list objects that a subject has access to.
// It contains all necessary information to query the FGA service for accessible objects.
//
// This struct is used to define the parameters for listing objects in the FGA system.
// It supports both direct object access checks and type-based object listing.
type ListRequest struct {
	// SubjectID is the unique identifier of the subject (e.g., user ID)
	// This is required and must not be empty.
	SubjectID string

	// SubjectType is the type of the subject (e.g., "user", "organization")
	// If not specified, the defaultSubject ("user") will be used.
	SubjectType string

	// ObjectType is the type of objects to list (e.g., "document", "space")
	// This is required and must not be empty.
	ObjectType string

	// ObjectID is the unique identifier of the object
	// This is optional and can be empty for type-based listing.
	ObjectID string

	// Relation is the relation to check (e.g., "viewer", "editor")
	// This is required and must not be empty.
	Relation string

	// Context is the context of the request used for conditional relationships
	// This is optional and can be nil.
	Context *map[string]any
}

// toListObjectsRequest converts the ListRequest to a client.ClientListObjectsRequest.
// It handles the conversion of subject and object types to the format expected by the FGA service.
//
// The method ensures proper formatting of the request by:
// 1. Converting the subject to the correct FGA entity format
// 2. Lowercasing the object type for consistency
// 3. Including any provided context
//
// Returns:
//   - client.ClientListObjectsRequest: The converted request ready for FGA service
func (lr ListRequest) toListObjectsRequest() client.ClientListObjectsRequest {
	sub := Entity{
		Kind:       Kind(lr.SubjectType),
		Identifier: lr.SubjectID,
	}

	listReq := client.ClientListObjectsRequest{
		User:     sub.String(),
		Relation: lr.Relation,
		Type:     strings.ToLower(lr.ObjectType),
	}

	if lr.Context != nil {
		listReq.Context = lr.Context
	}

	return listReq
}

// ListObjectIDsWithAccess returns a list of object IDs for which the given subject has the specified relation.
//
// This method queries the FGA service to determine which objects of the specified type
// the subject has access to through the given relation. It's useful for scenarios like:
// - Listing all documents a user can view
// - Finding all spaces a user is a member of
// - Retrieving all resources a user has specific permissions for
//
// Example:
//
//	objectIDs, err := client.ListObjectIDsWithAccess(ctx, ListRequest{
//	    SubjectID:   "user123",
//	    SubjectType: "user",
//	    ObjectType:  "document",
//	    Relation:    "viewer",
//	})
//
// Parameters:
//   - ctx: Request-scoped context
//   - req: The list request including subject type, subject ID, object type, and relation
//
// Returns:
//   - []string: A list of object IDs the subject has access to
//   - error: If the FGA query failed or was invalid
func (c *Client) ListObjectIDsWithAccess(ctx context.Context, req ListRequest) ([]string, error) {
	list, err := c.listObjects(ctx, req.toListObjectsRequest())
	if err != nil {
		return nil, err
	}

	objectIDs := make([]string, len(list.Objects))

	for i := range list.Objects {
		entity, err := ParseEntity(list.Objects[i])
		if err != nil {
			return nil, err
		}

		objectIDs[i] = entity.Identifier
	}

	return objectIDs, nil
}

// ListAccess is a struct to hold the information needed to list all relations.
// It's used to check which relations a subject has to a specific object.
//
// This struct is particularly useful when you need to:
// - Check all permissions a user has on a specific resource
// - Verify multiple relations in a single request
// - Get a complete picture of a subject's access rights
type ListAccess struct {
	// SubjectType is the type of the subject (e.g., "user", "organization")
	// If not specified, the defaultSubject ("user") will be used.
	SubjectType string

	// SubjectID is the unique identifier of the subject
	// This is required and must not be empty.
	SubjectID string

	// ObjectType is the type of the object (e.g., "document", "space")
	// This is required and must not be empty.
	ObjectType string

	// ObjectID is the unique identifier of the object
	// This is required and must not be empty.
	ObjectID string

	// Relations is a list of specific relations to check
	// If nil, all relations from the model will be checked.
	Relations []string

	// Context is the context of the request used for conditional relationships
	// This is optional and can be nil.
	Context *map[string]any
}

// ListRelations returns a list of relations that the subject has to the specified object.
//
// This method performs a batch check of all possible relations between the subject and object.
// If no specific relations are provided, it will check all relations defined in the FGA model.
//
// Example:
//
//	relations, err := client.ListRelations(ctx, ListAccess{
//	    SubjectID:   "user123",
//	    ObjectType:  "document",
//	    ObjectID:    "doc456",
//	})
//
// Parameters:
//   - ctx: Request-scoped context
//   - ac: The list access request containing subject and object information
//
// Returns:
//   - []string: A list of relations the subject has to the object
//   - error: If the FGA query failed or was invalid
func (c *Client) ListRelations(ctx context.Context, ac ListAccess) ([]string, error) {
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

	relations := ac.Relations
	if relations == nil {
		var err error

		relations, err = c.getRelationsFromModel(ctx, ac.ObjectType)
		if err != nil {
			log.Error().
				Err(err).
				Str("objectType", ac.ObjectType).
				Msg("failed to get relations from model")

			return nil, err
		}
	}

	checks := make([]client.ClientBatchCheckItem, 0, len(relations))
	for _, rel := range relations {
		checks = append(checks, client.ClientBatchCheckItem{
			User:          sub.String(),
			Relation:      rel,
			Object:        obj.String(),
			Context:       ac.Context,
			CorrelationId: ulid.Make().String(),
		})
	}

	return c.batchCheckTuples(ctx, checks)
}

// batchCheckTuples performs a batch check of multiple relations in a single request.
//
// This method is used internally to efficiently check multiple relations at once.
// It processes the results and returns only the relations that are allowed.
//
// Parameters:
//   - ctx: Request-scoped context
//   - checks: A slice of batch check items to process
//
// Returns:
//   - []string: A list of allowed relations
//   - error: If the batch check failed
func (c *Client) batchCheckTuples(ctx context.Context, checks []client.ClientBatchCheckItem) ([]string, error) {
	if len(checks) == 0 {
		return []string{}, nil
	}

	res, err := c.client.BatchCheck(ctx).Body(
		client.ClientBatchCheckRequest{
			Checks: checks,
		},
	).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int("checkCount", len(checks)).
			Msg("failed to execute batch check")

		return nil, err
	}

	if res == nil {
		return nil, ErrEmptyBatchCheckResponse
	}

	relations := make([]string, 0, len(checks))

	for id, r := range *res.Result {
		if !r.GetAllowed() {
			continue
		}

		check, ok := getCheckItemByCorrelationID(id, checks)
		if !ok {
			log.Error().
				Str("correlationID", id).
				Msg("correlation ID not found in checks")

			continue
		}

		relations = append(relations, check.Relation)
	}

	return relations, nil
}

// getRelationsFromModel retrieves all possible relations for a given object type from the FGA model.
//
// This method queries the FGA service to get the authorization model and extracts
// all relations defined for the specified object type.
//
// Parameters:
//   - ctx: Request-scoped context
//   - objectType: The type of object to get relations for
//
// Returns:
//   - []string: A list of all possible relations for the object type
//   - error: If the model query failed
func (c *Client) getRelationsFromModel(ctx context.Context, objectType string) ([]string, error) {
	model, err := c.client.ReadAuthorizationModel(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("objectType", objectType).
			Msg("failed to read authorization model")

		return nil, err
	}

	authorizationModel := model.GetAuthorizationModel()
	typeDefs := authorizationModel.GetTypeDefinitions()

	relations := make([]string, 0, len(typeDefs))

	for _, typeDef := range typeDefs {
		if !strings.EqualFold(typeDef.GetType(), objectType) {
			continue
		}

		for k := range typeDef.GetRelations() {
			relations = append(relations, k)
		}
	}

	return relations, nil
}

// listObjects performs the actual FGA service query to list accessible objects.
//
// This is the low-level method that communicates with the FGA service to get the list
// of objects that match the specified criteria.
//
// Parameters:
//   - ctx: Request-scoped context
//   - req: The list objects request to send to the FGA service
//
// Returns:
//   - *client.ClientListObjectsResponse: The response from the FGA service
//   - error: If the query failed
func (c *Client) listObjects(ctx context.Context, req client.ClientListObjectsRequest) (*client.ClientListObjectsResponse, error) {
	list, err := c.client.ListObjects(ctx).Body(req).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Interface("request", req).
			Msg("failed to list objects")

		return nil, err
	}

	return list, nil
}
