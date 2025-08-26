package ast

import "github.com/codevault-llc/php-lint/internal/token"

// Node is any AST node.
type Node interface {
    String() string
    Pos() token.Pos // Start position of the node
    End() token.Pos  // End position of the node
}

// Base is a helper struct embedded in all AST nodes to store their span.
type Base struct {
	S, E token.Pos
}

func (b *Base) Pos() token.Pos { return b.S }
func (b *Base) End() token.Pos { return b.E }

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
