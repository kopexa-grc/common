// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"
	"strings"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// default initial capacity for user slice (avoids magic number lint)
const defaultUserCap = 8

// Expand returns an OpenFGA expand request builder
// See: https://openfga.dev/docs/interacting/relationship-queries#expand
func (c *Client) Expand(ctx context.Context) client.SdkClientExpandRequestInterface {
	return c.client.Expand(ctx)
}

// uses go-sdk: github.com/openfga/go-sdk/client
// ListUsersWithAccess returns the list of user IDs that have the given relation to the object.
//
// It performs an Expand query (object#relation) and walks the returned userset tree collecting
// all leaf user subjects (type == "user"). Duplicate user IDs are de-duplicated while preserving
// discovery order.
//
// Parameters:
//   - ctx: Request context
//   - ot: Object type (e.g. "space")
//   - oid: Object identifier (e.g. "123")
//   - rel: Relation name (e.g. "member")
//
// Returns:
//   - []string of user IDs or empty slice if none
//   - error when the expand call fails
//
// nolint:gocyclo
func (c *Client) ListUsersWithAccess(ctx context.Context, ot, oid, rel string) ([]string, error) {
	if ot == "" || oid == "" || rel == "" {
		return []string{}, nil
	}

	object := strings.ToLower(ot) + ":" + oid

	resp, err := c.client.Expand(ctx).
		Body(client.ClientExpandRequest{Object: object, Relation: rel}).
		Execute()

	if err != nil {
		log.Error().Err(err).Str("object", object).Str("relation", rel).Msg("failed to expand userset")

		return nil, err
	}

	if resp == nil {
		return []string{}, nil
	}

	tree, ok := resp.GetTreeOk()
	if !ok || tree == nil {
		return []string{}, nil
	}

	root, ok := tree.GetRootOk()
	if !ok || root == nil {
		return []string{}, nil
	}

	return traverseUserset(root), nil
}

// traverseUserset walks a userset tree root and returns unique user IDs.
func traverseUserset(root *openfga.Node) []string {
	if root == nil {
		return []string{}
	}

	seen := make(map[string]struct{})
	out := make([]string, 0, defaultUserCap)

	stack := []*openfga.Node{root}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n == nil {
			continue
		}

		if collectLeafUsers(n, seen, &out) {
			continue
		}

		if pushUnionNodes(n, &stack) {
			continue
		}

		if pushIntersectionNodes(n, &stack) {
			continue
		}

		if pushDifferenceNodes(n, &stack) {
			continue
		}
	}

	return out
}

// collectLeafUsers extracts users from a leaf node. Returns true if node was a leaf.
func collectLeafUsers(n *openfga.Node, seen map[string]struct{}, out *[]string) bool {
	leaf, ok := n.GetLeafOk()
	if !ok || leaf == nil {
		return false
	}

	users, oku := leaf.GetUsersOk()
	if !oku || users == nil {
		return true
	}

	list, okl := users.GetUsersOk()
	if !okl || list == nil {
		return true
	}

	for _, subject := range *list {
		if !strings.HasPrefix(subject, "user:") {
			continue
		}

		id := strings.TrimPrefix(subject, "user:")
		if id == "" {
			continue
		}

		if _, dup := seen[id]; dup {
			continue
		}

		seen[id] = struct{}{}

		*out = append(*out, id)
	}

	return true
}

// pushUnionNodes pushes child nodes for a union node.
func pushUnionNodes(n *openfga.Node, stack *[]*openfga.Node) bool {
	union, ok := n.GetUnionOk()
	if !ok || union == nil {
		return false
	}

	nodes, okn := union.GetNodesOk()
	if !okn || nodes == nil {
		return true
	}

	for i := range *nodes {
		*stack = append(*stack, &(*nodes)[i])
	}

	return true
}

// pushIntersectionNodes pushes child nodes for an intersection node.
func pushIntersectionNodes(n *openfga.Node, stack *[]*openfga.Node) bool {
	inter, ok := n.GetIntersectionOk()
	if !ok || inter == nil {
		return false
	}

	nodes, okn := inter.GetNodesOk()
	if !okn || nodes == nil {
		return true
	}

	for i := range *nodes {
		*stack = append(*stack, &(*nodes)[i])
	}

	return true
}

// pushDifferenceNodes pushes base and subtract nodes for a difference node.
func pushDifferenceNodes(n *openfga.Node, stack *[]*openfga.Node) bool {
	diff, ok := n.GetDifferenceOk()
	if !ok || diff == nil {
		return false
	}

	if base, okb := diff.GetBaseOk(); okb && base != nil {
		*stack = append(*stack, base)
	}

	if sub, oks := diff.GetSubtractOk(); oks && sub != nil {
		*stack = append(*stack, sub)
	}

	return true
}
