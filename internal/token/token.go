package token

import "fmt"

// Span represents a region of source code.
type Span struct {
	Start, End Pos
}

// Pos represents a position in a file.
type Pos struct {
	Line, Col, Offset int
}

// Kind is the type of a token.
type Kind string

const (
	// Special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literals
	IDENT  = "IDENT"
	STRING = "STRING"

	// Keywords & Dangerous Functions
	ECHO       = "ECHO"
	EVAL       = "EVAL"
	EXIT       = "EXIT"
	DIE        = "DIE"
	SHELL_EXEC = "SHELL_EXEC"
	EXEC       = "EXEC"
	PASSTHRU   = "PASSTHRU"
	SYSTEM     = "SYSTEM"
	FUNCTION   = "FUNCTION"

	// Delimiters
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	COMMA     = ","
	RBRACE    = "}"
	LBRACE    = "{"

	// Misc
	WHITESPACE = "WHITESPACE"

	// PHP Tags
	OPEN_TAG  = "<?php"
	CLOSE_TAG = "?>"
)

// Token represents a lexical token.
type Token struct {
	Kind   Kind
	Lexeme string
	Span   Span
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Kind: %s, Lexeme: \"%s\"}", t.Kind, t.Lexeme)
}

var keywords = map[string]Kind{
	"echo":       ECHO,
	"eval":       EVAL,
	"exit":       EXIT,
	"die":        DIE,
	"shell_exec": SHELL_EXEC,
	"exec":       EXEC,
	"passthru":   PASSTHRU,
	"system":     SYSTEM,
}

// LookupIdent checks the keywords table to see if a given identifier is a keyword.
func LookupIdent(ident string) Kind {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}