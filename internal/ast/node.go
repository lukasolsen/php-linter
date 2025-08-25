package ast

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
