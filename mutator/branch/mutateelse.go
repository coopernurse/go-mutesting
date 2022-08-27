package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("branch/else", MutatorElse)
}

// MutatorElse implements a mutator for else branches.
func MutatorElse(input mutator.MutatorInput) []mutator.Mutation {
	n, ok := input.Node.(*ast.IfStmt)
	if !ok {
		return nil
	}
	// We ignore else ifs and nil blocks
	_, ok = n.Else.(*ast.IfStmt)
	if ok || n.Else == nil {
		return nil
	}

	old := n.Else

	return []mutator.Mutation{
		{
			Change: func() {
				n.Else = astutil.CreateNoopOfStatement(input.Pkg, input.Info, old)
			},
			Reset: func() {
				n.Else = old
			},
		},
	}
}
