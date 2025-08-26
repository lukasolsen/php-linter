package rules

import (
	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type callExprVisitor struct {
	issues   []types.Issue
	ruleName string
	check    func(node *ast.CallExpr) (*types.Issue, bool)
}

func (v *callExprVisitor) Visit(node ast.Node) {
	if n, ok := node.(*ast.CallExpr); ok {
		if issue, found := v.check(n); found {
			v.issues = append(v.issues, *issue)
		}
	}
}