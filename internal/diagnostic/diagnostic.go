package diagnostic

import (
	"fmt"

	"github.com/codevault-llc/php-lint/internal/token"
)

// Severity of a diagnostic.
type Severity int


const (
	Error Severity = iota
	Warning
	Info
)


type Diagnostic struct {
	Severity Severity
	Message string
	Span token.Span
}


func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: %s at %d:%d", map[Severity]string{Error: "error", Warning: "warning", Info: "info"}[d.Severity], d.Message, d.Span.Start.Line, d.Span.Start.Col)
}
func (d Diagnostic) Location() string {
	return fmt.Sprintf("%d:%d", d.Span.Start.Line, d.Span.Start.Col)
}
