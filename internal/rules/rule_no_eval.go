package rules

import (
	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/token"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type RuleNoEval struct{}

func (r *RuleNoEval) Name() string { return "security-no-eval" }

func (r *RuleNoEval) Description() string { return "Disallows the use of the eval() function." }

func (r *RuleNoEval) Check(filename string, content []byte, program *ast.Program) []types.Issue {
	visitor := &callExprVisitor{
		issues:   []types.Issue{},
		ruleName: r.Name(),
		check: func(node *ast.CallExpr) (bool, string) {
			if ident, ok := node.Function.(*ast.Identifier); ok && ident.Token.Kind == token.EVAL {
				return true, "Use of eval() is a significant security risk and is strongly discouraged"
			}
			return false, ""
		},
	}
	ast.Walk(program, visitor)
	return visitor.issues
}