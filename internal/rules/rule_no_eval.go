package rules

import (
	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/codevault-llc/php-lint/internal/token"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type RuleNoEval struct{}

func (r *RuleNoEval) Name() string { return "security-no-eval" }

func (r *RuleNoEval) Description() string { return "Disallows the use of the eval() function." }

func (r *RuleNoEval) Check(filename string, content []byte, program *ast.Program, symbolTable *stubs.SymbolTable) []types.Issue {
	visitor := &callExprVisitor{
		issues:   []types.Issue{},
		ruleName: r.Name(),
		check: func(node *ast.CallExpr) (*types.Issue, bool) {
			if ident, ok := node.Function.(*ast.Identifier); ok && ident.Token.Kind == token.EVAL {
				issue := &types.Issue{
					RuleName: r.Name(),
					Message:  "Use of eval() is a significant security risk and is strongly discouraged",
					Range:    ident.Token.Span,
					Severity: 2,
					Source:   "php-lint",
				}

				return issue, true
			}

			return nil, false
		},
	}
	ast.Walk(program, visitor)
	return visitor.issues
}