package types

import (
	"github.com/codevault-llc/php-lint/internal/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Issue struct {
	RuleName string
	Message  string
	Range    token.Span

	// LSP Information
	Severity protocol.DiagnosticSeverity
	Source   string
}