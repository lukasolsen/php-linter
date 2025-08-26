package ast

import (
	"fmt"

	"github.com/codevault-llc/php-lint/internal/token"
)

// StringLiteral represents a string.
type StringLiteral struct {
    Base
    Token token.Token // The string token
    Value string
}

func (sl *StringLiteral) isExpr() {}
func (sl *StringLiteral) String() string {
    return fmt.Sprintf("'%s'", sl.Value)
}

// Identifier represents a variable or function name.
type Identifier struct {
	Base           // Embedded position
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) isExpr()        {}
func (i *Identifier) String() string { return i.Value }

// CallExpr represents a function call.
type CallExpr struct {
    Base
    Token     token.Token // The '(' token
    Function  Expr        // Identifier or another expression
    Arguments []Expr
}

func (ce *CallExpr) isExpr() {}
func (ce *CallExpr) String() string {
    if ce.Function != nil {
        return ce.Function.String() + "()"
    }
    return "<invalid call>"
}
