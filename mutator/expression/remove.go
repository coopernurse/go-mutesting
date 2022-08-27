package expression

import (
	"go/ast"
	"go/token"

	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("expression/remove", MutatorRemoveTerm)
}

// MutatorRemoveTerm implements a mutator to remove expression terms.
func MutatorRemoveTerm(input mutator.MutatorInput) []mutator.Mutation {
	n, ok := input.Node.(*ast.BinaryExpr)
	if !ok {
		return nil
	}
	if n.Op != token.LAND && n.Op != token.LOR {
		return nil
	}

	var r *ast.Ident

	switch n.Op {
	case token.LAND:
		r = ast.NewIdent("true")
	case token.LOR:
		r = ast.NewIdent("false")
	}

	x := n.X
	y := n.Y

	return []mutator.Mutation{
		{
			Change: func() {
				n.X = r
			},
			Reset: func() {
				n.X = x
			},
		},
		{
			Change: func() {
				n.Y = r
			},
			Reset: func() {
				n.Y = y
			},
		},
	}
}
