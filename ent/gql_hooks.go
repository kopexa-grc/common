package ent

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

// Skipping err113 linting since these errors are returned during generation and not runtime
//
//nolint:err113
var (
	addJSONScalar = func(_ *gen.Graph, s *ast.Schema) error {
		s.Types["JSON"] = &ast.Definition{
			Kind:        ast.Scalar,
			Description: "A valid JSON string.",
			Name:        "JSON",
		}
		return nil
	}
)

// import string mutations from entc
var (
	_ entc.Extension = (*Extension)(nil)
)
