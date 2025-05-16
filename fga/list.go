// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"
	"strings"

	"github.com/openfga/go-sdk/client"
)

// ListRequest represents a request to list objects that a subject has access to.
// It contains all necessary information to query the FGA service for accessible objects.
type ListRequest struct {
	// SubjectID is the unique identifier of the subject
	SubjectID string
	// SubjectType is the type of the subject (e.g., "user", "organization")
	SubjectType string
	// ObjectType is the type of objects to list (e.g., "document", "space")
	ObjectType string
	// Relation is the relation to check (e.g., "viewer", "editor")
	Relation string
	// Context is the context of the request used for conditional relationships
	Context *map[string]any
}

// toListObjectsRequest converts the ListRequest to a client.ClientListObjectsRequest.
// It handles the conversion of subject and object types to the format expected by the FGA service.
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
// the subject has access to through the given relation.
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

// listObjects performs the actual FGA service query to list accessible objects.
// It sends the request to the FGA service and returns the response.
func (c *Client) listObjects(ctx context.Context, req client.ClientListObjectsRequest) (*client.ClientListObjectsResponse, error) {
	list, err := c.client.ListObjects(ctx).Body(req).Execute()
	if err != nil {
		return nil, err
	}

	return list, nil
}
