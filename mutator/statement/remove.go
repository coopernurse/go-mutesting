package statement

import (
	"fmt"
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
		keep := true

		if input.Options.NameExclude != nil {
			idents := astutil.IdentifiersInStatement(input.Pkg, input.Info, ni)
			for _, id := range idents {
				name := nodeName(id)
				if input.Options.NameExclude.FindString(name) != "" {
					keep = false
					break
				}
				//fmt.Printf("RemoveStmt %d: %v\n", x, name)
			}
		}

		if keep && checkRemoveStatement(ni) {
			li := i
			old := l[li]

			mutations = append(mutations, mutator.Mutation{
				Change: func() {
					l[li] = astutil.CreateNoopOfStatement(input.Pkg, input.Info, old)
				},
				Reset: func() {
					l[li] = old
				},
			})
		}
	}

	return mutations
}

func nodeName(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return v.Sel.Name
	default:
		return fmt.Sprintf("%T", e)
	}
}
