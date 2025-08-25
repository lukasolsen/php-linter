package lexer

import (
	"fmt"
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

// peekChar looks ahead in the input without consuming the character.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPosition])
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
	case '<':
		if l.peekChar() == '?' {
			// Check for '<?php'
			if len(l.input) > l.position+4 && l.input[l.position:l.position+5] == "<?php" {
				// Consume '<?php'
				l.readChar() // ?
				l.readChar() // p
				l.readChar() // h
				l.readChar() // p
				tok = l.newToken(token.OPEN_TAG, "<?php")
			}
		}
	case '?':
		if l.peekChar() == '>' {
			l.readChar()
			tok = l.newToken(token.CLOSE_TAG, "?>")
		}
	case '(':
		tok = l.newToken(token.LPAREN, "(")
	case ')':
		tok = l.newToken(token.RPAREN, ")")
	case ';':
		tok = l.newToken(token.SEMICOLON, ";")
	case '\'':
		tok.Kind = token.STRING
		tok.Lexeme = l.readString()
	case 0:
		tok.Lexeme = ""
		tok.Kind = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Lexeme = l.readIdentifier()
			tok.Kind = token.LookupIdent(tok.Lexeme)
			return tok
		} else {
			fmt.Println("Illegal character listed !!", string(l.ch))
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	var sb strings.Builder
	l.readChar() // Consume opening '
	for l.ch != '\'' && l.ch != 0 {
		sb.WriteRune(l.ch)
		l.readChar()
	}
	return sb.String()
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func (l *Lexer) newToken(kind token.Kind, lexeme string) token.Token {
	// Note: Span calculation can be improved for multi-char tokens
	return token.Token{Kind: kind, Lexeme: lexeme}
}