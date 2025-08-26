package lexer

import (
	"log"

	"github.com/codevault-llc/php-lint/internal/token"
)

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	startPos := token.Pos{Line: l.line, Col: l.column, Offset: l.position}
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
				tok = l.newTokenFromPos(token.OPEN_TAG, "<?php", startPos)
			}
		}
	case '?':
		if l.peekChar() == '>' {
			l.readChar()
			tok = l.newTokenFromPos(token.CLOSE_TAG, "?>", startPos)
		}
	case '(':
		tok = l.newTokenFromPos(token.LPAREN, "(", startPos)
	case ')':
		tok = l.newTokenFromPos(token.RPAREN, ")", startPos)
	case '{':
		tok = l.newTokenFromPos(token.LBRACE, "{", startPos)
	case ',':
		tok = l.newTokenFromPos(token.COMMA, ",", startPos)
	case '"':
		tok = l.newTokenFromPos(token.STRING, l.readString(), startPos)
	case '}':
		tok = l.newTokenFromPos(token.RBRACE, "}", startPos)
	case '/':
		if l.peekChar() == '*' {
			l.readChar()
			tok = l.newTokenFromPos(token.BLOCK_COMMENT, l.readBlockComment(), startPos)
		} else if l.peekChar() == '/' {
			l.readChar()
			tok = l.newTokenFromPos(token.LINE_COMMENT, l.readLineComment(), startPos)
		}
	case ';':
		tok = l.newTokenFromPos(token.SEMICOLON, ";", startPos)
	case '\'':
		tok = l.newTokenFromPos(token.STRING, l.readString(), startPos)
	case 0:
		tok = l.newTokenFromPos(token.EOF, "", startPos)
	default:
		if isLetter(l.ch) {
    ident := l.readIdentifier()
    tok = l.newTokenFromPos(token.LookupIdent(ident), ident, startPos)
    
		return tok
	} else {
			log.Println("Illegal character listed !!", string(l.ch))
			tok = l.newTokenFromPos(token.ILLEGAL, string(l.ch), startPos)
		}
	}

	l.readChar()
	tok.Span.End = token.Pos{Line: l.line, Col: l.column, Offset: l.position}
	return tok
}