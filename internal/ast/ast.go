package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/codevault-llc/php-lint/internal/token"
)

// Node is any AST node.
type Node interface {
	String() string
}

// Stmt is any statement/declaration node.
type Stmt interface {
	Node
	isStmt()
}

// Expr is any expression node.
type Expr interface {
	Node
	isExpr()
}

// Program is the root node for a PHP file.
type Program struct {
	Stmts []Stmt
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}

// EchoStmt represents an 'echo' statement, e.g., echo "Hello", "World";
type EchoStmt struct {
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

// StringLiteral represents a string.
type StringLiteral struct {
	Token token.Token // The string token
	Value string
}

func (sl *StringLiteral) isExpr() {}
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("'%s'", sl.Value)
}

// ExpressionStatement holds an expression.
type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expr
}

func (es *ExpressionStatement) isStmt()            {}
func (f *FunctionDeclStmt) isExpr()                {}
func (es *ExpressionStatement) String() string     { return es.Expression.String() }

// FunctionDeclStmt represents a 'function' statement, e.g., function my_func($a) {}
type FunctionDeclStmt struct {
	Token token.Token    // The 'function' token
	Name  *Identifier
}

func (fds *FunctionDeclStmt) isStmt() {}
func (fds *FunctionDeclStmt) String() string {
	return "function " + fds.Name.String()
}

// Identifier represents a variable or function name.
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) isExpr()        {}
func (i *Identifier) String() string { return i.Value }

// CallExpr represents a function call.
type CallExpr struct {
	Token     token.Token // The '(' token
	Function  Expr        // Identifier or another expression
	Arguments []Expr
}

func (ce *CallExpr) isExpr() {}
func (ce *CallExpr) String() string {
	// ... implementation for string representation ...
	return ""
}

// --- AST Walker ---

// Visitor defines the Visit method for the AST walker.
type Visitor interface {
	Visit(node Node)
}

// Walk traverses an AST in depth-first order.
func Walk(node Node, visitor Visitor) {
	visitor.Visit(node)

	switch n := node.(type) {
	case *Program:
		for _, stmt := range n.Stmts {
			Walk(stmt, visitor)
		}
	case *ExpressionStatement:
		Walk(n.Expression, visitor)
	case *CallExpr:
		Walk(n.Function, visitor)
		for _, arg := range n.Arguments {
			Walk(arg, visitor)
		}
	}
}
