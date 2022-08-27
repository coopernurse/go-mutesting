package statement

import (
	"go/ast"
	"go/token"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("statement/remove", MutatorRemoveStatement)
}

func checkRemoveStatement(node ast.Stmt) bool {
	switch n := node.(type) {
	case *ast.AssignStmt:
		if n.Tok != token.DEFINE {
			return true
		}
	case *ast.ExprStmt, *ast.IncDecStmt:
		return true
	}

	return false
}

// MutatorRemoveStatement implements a mutator to remove statements.
func MutatorRemoveStatement(input mutator.MutatorInput) []mutator.Mutation {
	var l []ast.Stmt

	switch n := input.Node.(type) {
	case *ast.BlockStmt:
		l = n.List
	case *ast.CaseClause:
		l = n.Body
	}

	var mutations []mutator.Mutation

	for i, ni := range l {
		if checkRemoveStatement(ni) {
			li := i
			old := l[li]

			newStmt, modified := astutil.CreateNoopOfStatement(input.Pkg, input.Info, old, input.Options)

			if modified {
				mutations = append(mutations, mutator.Mutation{
					Change: func() {
						l[li] = newStmt
					},
					Reset: func() {
						l[li] = old
					},
				})
			}
		}
	}

	return mutations
}
