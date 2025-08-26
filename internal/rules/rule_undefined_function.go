package rules

import (
	"fmt"
	"log"

	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type RuleUndefinedFunction struct{}

func (r *RuleUndefinedFunction) Name() string { return "undefined-function" }
func (r *RuleUndefinedFunction) Description() string {
	return "Reports calls to functions that are not defined."
}

func (r *RuleUndefinedFunction) Check(filename string, content []byte, program *ast.Program, symbolTable *stubs.SymbolTable) []types.Issue {
	visitor := &callExprVisitor{
		issues:   []types.Issue{},
		ruleName: r.Name(),
		check: func(node *ast.CallExpr) (*types.Issue, bool) {
			if ident, ok := node.Function.(*ast.Identifier); ok {
				if !symbolTable.IsFunctionDefined(ident.Value) {
					log.Printf("Undefined function %s() called", ident.Token)

					issue := types.Issue{
						RuleName: r.Name(),
						Message:  fmt.Sprintf("Call to undefined function %s()", ident.Value),
						Range:    ident.Token.Span,
						Severity: 2,
						Source:   "php-lint",
					}

					return &issue, true
				}
			}

			return nil, false
		},
	}
	ast.Walk(program, visitor)
	return visitor.issues
}