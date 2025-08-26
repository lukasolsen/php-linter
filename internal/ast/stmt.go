package ast

import (
	"bytes"
	"strings"

	"github.com/codevault-llc/php-lint/internal/token"
)

// EchoStmt represents an 'echo' statement, e.g., echo "Hello", "World";
type EchoStmt struct {
    Base
    Token       token.Token // The 'echo' token
    Expressions []Expr
}

func (es *EchoStmt) isStmt() {}
func (es *EchoStmt) String() string {
    var out bytes.Buffer
    out.WriteString(es.Token.Lexeme + " ")
    expressions := []string{}
    for _, e := range es.Expressions {
        expressions = append(expressions, e.String())
    }
    out.WriteString(strings.Join(expressions, ", "))
    out.WriteString(";")
    return out.String()
}

// ExpressionStatement holds an expression.
type ExpressionStatement struct {
    Base
    Token      token.Token // The first token of the expression
    Expression Expr
}

func (es *ExpressionStatement) isStmt() {}
func (es *ExpressionStatement) String() string {
    if es.Expression != nil {
        return es.Expression.String()
    }
    return ""
}

// FunctionDeclStmt represents a 'function' statement, e.g., function my_func($a) {}.
type FunctionDeclStmt struct {
    Base
    Token token.Token // The 'function' token
    Name  *Identifier
}

func (fds *FunctionDeclStmt) isStmt() {}
func (fds *FunctionDeclStmt) String() string {
    if fds.Name != nil {
        return "function " + fds.Name.String()
    }
    return "<invalid function decl>"
}
