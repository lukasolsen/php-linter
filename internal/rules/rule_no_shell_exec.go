package rules

import (
	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/codevault-llc/php-lint/internal/token"
	"github.com/codevault-llc/php-lint/pkg/types"
)

type RuleNoShellExec struct{}

func (r *RuleNoShellExec) Name() string { return "security-no-shell-exec" }

func (r *RuleNoShellExec) Description() string { return "Disallows the use of shell_exec() and similar functions." }

func (r *RuleNoShellExec) Check(filename string, content []byte, program *ast.Program, symbolTable *stubs.SymbolTable) []types.Issue {
	visitor := &callExprVisitor{
		ruleName: r.Name(),
		check: func(node *ast.CallExpr) (bool, string) {
			dangerousFunctions := map[token.Kind]bool{
				token.SHELL_EXEC: true, token.EXEC: true, token.PASSTHRU: true, token.SYSTEM: true,
			}
			if ident, ok := node.Function.(*ast.Identifier); ok {
				if _, found := dangerousFunctions[ident.Token.Kind]; found {
					return true, "Execution of shell commands is a security risk"
				}
			}
			return false, ""
		},
	}
	ast.Walk(program, visitor)
	return visitor.issues
}