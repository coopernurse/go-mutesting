package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("branch/case", MutatorCase)
}

// MutatorCase implements a mutator for case clauses.
func MutatorCase(input mutator.MutatorInput) []mutator.Mutation {
	n, ok := input.Node.(*ast.CaseClause)
	if !ok {
		return nil
	}

	old := n.Body

	newStmts, modified := astutil.CreateNoopOfStatements(input.Pkg, input.Info, n.Body, input.Options)
	if !modified {
		return nil
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Body = newStmts
			},
			Reset: func() {
				n.Body = old
			},
		},
	}
}
