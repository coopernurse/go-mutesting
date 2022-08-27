package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("branch/if", MutatorIf)
}

// MutatorIf implements a mutator for if and else if branches.
func MutatorIf(input mutator.MutatorInput) []mutator.Mutation {
	n, ok := input.Node.(*ast.IfStmt)
	if !ok {
		return nil
	}

	old := n.Body.List

	return []mutator.Mutation{
		{
			Change: func() {
				n.Body.List = []ast.Stmt{
					astutil.CreateNoopOfStatement(input.Pkg, input.Info, n.Body),
				}
			},
			Reset: func() {
				n.Body.List = old
			},
		},
	}
}
