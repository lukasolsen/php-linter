package lexer

import (
	"strings"
	"unicode"

	"github.com/codevault-llc/php-lint/internal/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line, column int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) readString() string {
	// Consumes the opening '
	l.readChar()
	start := l.position
	for l.ch != '\'' && l.ch != 0 {
		l.readChar()
	}
	end := l.position
	l.readChar() // Consumes the closing '
	return l.input[start:end]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// newTokenFromPos is a helper to build a token from a starting position.
// The end position is inferred from the lexer's current state.
func (l *Lexer) newTokenFromPos(kind token.Kind, lexeme string, start token.Pos) token.Token {
	return token.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Span: token.Span{
			Start: start,
			End:   token.Pos{Line: l.line, Col: l.column, Offset: l.position},
		},
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func (l *Lexer) readBlockComment() string {
	var sb strings.Builder
	l.readChar() // Consume opening /*
	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // Consume *
			l.readChar() // Consume /
			break
		}
		sb.WriteRune(l.ch)
		l.readChar()
	}
	return sb.String()
}

func (l *Lexer) readLineComment() string {
	var sb strings.Builder
	l.readChar() 
	for l.ch != 0 && l.ch != '\n' {
		sb.WriteRune(l.ch)
		l.readChar()
	}
	return sb.String()
}