package ast

// Visitor defines the Visit method for the AST walker.
type Visitor interface {
    Visit(node Node)
}

// Walk traverses an AST in depth-first order.
func Walk(node Node, visitor Visitor) {
    if node == nil || visitor == nil {
        return
    }

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
