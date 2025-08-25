package rules

import (
	"fmt"

	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type callExprVisitor struct {
	issues    []types.Issue
	ruleName  string
	check     func(node *ast.CallExpr) (bool, string)
}

func (v *callExprVisitor) Visit(node ast.Node) {
	if n, ok := node.(*ast.CallExpr); ok {
		if ident, ok := n.Function.(*ast.Identifier); ok {
			if found, msg := v.check(n); found {
				v.issues = append(v.issues, types.Issue{RuleName: v.ruleName, Message: fmt.Sprintf("%s: %s()", msg, ident.Value)})
			}
		}
	}
}