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

	newStmt, modified := astutil.CreateNoopOfStatement(input.Pkg, input.Info, old, input.Options)
	if !modified {
		return nil
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Else = newStmt
			},
			Reset: func() {
				n.Else = old
			},
		},
	}
}
