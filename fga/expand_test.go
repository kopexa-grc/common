// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga_test

import (
	"context"
	"testing"

	"github.com/kopexa-grc/common/fga"
	"github.com/kopexa-grc/common/fga/internal/fgamock"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// helper to build a leaf users node with given subjects
func leafUsers(users ...string) *openfga.Node {
	leaf := openfga.NewLeaf()
	leaf.SetUsers(*openfga.NewUsers(users))
	nd := openfga.NewNode("leaf")
	nd.SetLeaf(*leaf)

	return nd
}

func union(nodes ...*openfga.Node) *openfga.Node {
	child := make([]openfga.Node, len(nodes))
	for i, n := range nodes { // copy values
		child[i] = *n
	}

	n := openfga.NewNode("union")
	n.SetUnion(*openfga.NewNodes(child))

	return n
}

func intersection(nodes ...*openfga.Node) *openfga.Node {
	child := make([]openfga.Node, len(nodes))
	for i, n := range nodes {
		child[i] = *n
	}

	n := openfga.NewNode("intersection")
	n.SetIntersection(*openfga.NewNodes(child))

	return n
}

func difference(base, subtract *openfga.Node) *openfga.Node {
	d := openfga.NewUsersetTreeDifference(*base, *subtract)
	n := openfga.NewNode("difference")
	n.SetDifference(*d)

	return n
}

// fakeExpandReq implements client.SdkClientExpandRequestInterface for testing.
type fakeExpandReq struct {
	body    *client.ClientExpandRequest
	execute func() (*client.ClientExpandResponse, error)
}

func (f *fakeExpandReq) Options(_ client.ClientExpandOptions) client.SdkClientExpandRequestInterface {
	return f
}
func (f *fakeExpandReq) Body(b client.ClientExpandRequest) client.SdkClientExpandRequestInterface {
	f.body = &b
	return f
}
func (f *fakeExpandReq) Execute() (*client.ClientExpandResponse, error) { return f.execute() }

// revive:disable-next-line var-naming (method name fixed by upstream interface)
func (f *fakeExpandReq) GetAuthorizationModelIdOverride() *string { return nil }

// revive:disable-next-line var-naming (method name fixed by upstream interface)
func (f *fakeExpandReq) GetStoreIdOverride() *string             { return nil }
func (f *fakeExpandReq) GetContext() context.Context             { return context.Background() }
func (f *fakeExpandReq) GetBody() *client.ClientExpandRequest    { return f.body }
func (f *fakeExpandReq) GetOptions() *client.ClientExpandOptions { return nil }

func TestClient_ListUsersWithAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)

	tests := []struct {
		name        string
		objectType  string
		objectID    string
		relation    string
		resp        *client.ClientExpandResponse
		err         error
		expected    []string
		expectError bool
	}{
		{
			name:       "empty params returns empty slice",
			objectType: "",
			objectID:   "",
			relation:   "",
			resp:       nil,
			expected:   []string{},
		},
		{
			name:       "single leaf users",
			objectType: "space",
			objectID:   "123",
			relation:   "member",
			resp: func() *client.ClientExpandResponse {
				root := leafUsers("user:alice", "user:bob")
				tree := openfga.NewUsersetTree()
				tree.SetRoot(*root)
				resp := openfga.ExpandResponse{Tree: tree}
				return (*client.ClientExpandResponse)(&resp)
			}(),
			expected: []string{"alice", "bob"},
		},
		{
			name:       "union + duplicates",
			objectType: "doc",
			objectID:   "42",
			relation:   "viewer",
			resp: func() *client.ClientExpandResponse {
				l1 := leafUsers("user:alice", "group:eng#member")
				l2 := leafUsers("user:bob", "user:alice")
				root := union(l1, l2)
				tree := openfga.NewUsersetTree()
				tree.SetRoot(*root)
				resp := openfga.ExpandResponse{Tree: tree}
				return (*client.ClientExpandResponse)(&resp)
			}(),
			expected: []string{"alice", "bob"},
		},
		{
			name:       "intersection nested (still collects all users)",
			objectType: "repo",
			objectID:   "999",
			relation:   "reader",
			resp: func() *client.ClientExpandResponse {
				l1 := leafUsers("user:alice")
				l2 := leafUsers("user:carol")
				root := intersection(l1, l2)
				tree := openfga.NewUsersetTree()
				tree.SetRoot(*root)
				resp := openfga.ExpandResponse{Tree: tree}
				return (*client.ClientExpandResponse)(&resp)
			}(),
			expected: []string{"alice", "carol"},
		},
		{
			name:       "difference traverses both sides",
			objectType: "board",
			objectID:   "ab1",
			relation:   "editor",
			resp: func() *client.ClientExpandResponse {
				base := leafUsers("user:alice", "user:dave")
				sub := leafUsers("user:bob")
				root := difference(base, sub)
				tree := openfga.NewUsersetTree()
				tree.SetRoot(*root)
				resp := openfga.ExpandResponse{Tree: tree}
				return (*client.ClientExpandResponse)(&resp)
			}(),
			expected: []string{"alice", "dave", "bob"}, // we collect all encountered users regardless of set semantics
		},
		{
			name:        "error from SDK",
			objectType:  "space",
			objectID:    "err",
			relation:    "member",
			err:         assert.AnError,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		// separate run to avoid gomock expectation bleed
		t.Run(tc.name, func(t *testing.T) {
			c := fga.NewMockFGAClient(mockSdk)

			if tc.objectType != "" && tc.objectID != "" && tc.relation != "" {
				fe := &fakeExpandReq{execute: func() (*client.ClientExpandResponse, error) { return tc.resp, tc.err }}

				mockSdk.EXPECT().Expand(gomock.Any()).DoAndReturn(func(_ context.Context) client.SdkClientExpandRequestInterface { return fe }).Times(1)
			}

			users, err := c.ListUsersWithAccess(context.Background(), tc.objectType, tc.objectID, tc.relation)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)

				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, users)
		})
	}
}
