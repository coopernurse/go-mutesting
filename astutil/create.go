package astutil

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/zimmski/go-mutesting/mutator"
)

// CreateNoopOfStatement creates a syntactically safe noop statement out of a given statement.
func CreateNoopOfStatement(pkg *types.Package, info *types.Info, stmt ast.Stmt, opts mutator.MutatorOptions) (ast.Stmt, bool) {
	keepOrigStmt := false
	if opts.NameExclude != nil {
		switch v := stmt.(type) {
		case *ast.AssignStmt:
			name := nodeName(v.Rhs[0])
			if opts.NameExclude.FindString(name) != "" {
				keepOrigStmt = true
				// } else {
				// 	fmt.Printf("AssignStmt keep %v - %T\n", name, stmt)
			}
		case *ast.ExprStmt:
			name := nodeName(v.X)
			if opts.NameExclude.FindString(name) != "" {
				keepOrigStmt = true
				// } else {
				// 	fmt.Printf("ExprStmt keep %v - %T\n", name, stmt)
			}
		default:
			idents := IdentifiersInStatement(pkg, info, stmt)
			for _, id := range idents {
				name := nodeName(id)
				if opts.NameExclude.FindString(name) != "" {
					keepOrigStmt = true
					// } else {
					// 	fmt.Printf("RemoveStmt keep %v - %T\n", name, stmt)
				}
			}
		}
	}

	if keepOrigStmt {
		return stmt, false
	}

	ids := IdentifiersInStatement(pkg, info, stmt)

	if len(ids) == 0 {
		return &ast.EmptyStmt{
			Semicolon: token.NoPos,
		}, true
	}

	lhs := make([]ast.Expr, len(ids))
	for i := range ids {
		lhs[i] = ast.NewIdent("_")
	}

	return &ast.AssignStmt{
		Lhs: lhs,
		Rhs: ids,
		Tok: token.ASSIGN,
	}, true
}

// CreateNoopOfStatements creates a syntactically safe noop statement out of a given statement.
func CreateNoopOfStatements(pkg *types.Package, info *types.Info, stmts []ast.Stmt, opts mutator.MutatorOptions) ([]ast.Stmt, bool) {
	var anyModified, modified bool
	out := make([]ast.Stmt, len(stmts))
	for i, st := range stmts {
		out[i], modified = CreateNoopOfStatement(pkg, info, st, opts)
		if modified {
			anyModified = true
		}
	}
	return out, anyModified
}

func nodeName(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return v.Sel.Name
	case *ast.CallExpr:
		return nodeName(v.Fun)
	default:
		return fmt.Sprintf("%T", e)
	}
}
