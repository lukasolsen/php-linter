package rules

import (
	"fmt"

	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type RuleUndefinedFunction struct{}

func (r *RuleUndefinedFunction) Name() string { return "undefined-function" }
func (r *RuleUndefinedFunction) Description() string {
	return "Reports calls to functions that are not defined."
}

// This rule's Check method now accepts the symbol table from the linter.
func (r *RuleUndefinedFunction) Check(filename string, content []byte, program *ast.Program) []types.Issue {
	visitor := &callExprVisitor{
		issues:   []types.Issue{},
		ruleName: r.Name(),
		check: func(node *ast.CallExpr) (bool, string) {
			if ident, ok := node.Function.(*ast.Identifier); ok {
				/*if !symbolTable.IsFunctionDefined(ident.Value) {
					return true, fmt.Sprintf("Call to undefined function %s()", ident.Value)
				}*/

				// Ignore symbols now
				return true, fmt.Sprintf("Ignoring undefined function %s() (symbol table not yet integrated)", ident.Value)
			}
			return false, ""
		},
	}
	ast.Walk(program, visitor)
	return visitor.issues
}