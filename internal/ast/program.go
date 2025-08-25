package ast

import (
	"bytes"
)

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
