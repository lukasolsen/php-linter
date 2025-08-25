package parser

import (
	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/lexer"
	"github.com/codevault-llc/php-lint/internal/token"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

type Parser struct {
	l      *lexer.Lexer
	curTok token.Token
	peekTok token.Token

	prefixParseFns map[token.Kind]prefixParseFn
	infixParseFns  map[token.Kind]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFns = make(map[token.Kind]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.EVAL, p.parseIdentifier) // Treat keywords like identifiers for parsing
	p.registerPrefix(token.EXIT, p.parseIdentifier)
	p.registerPrefix(token.DIE, p.parseIdentifier)
	p.registerPrefix(token.SHELL_EXEC, p.parseIdentifier)
	p.registerPrefix(token.EXEC, p.parseIdentifier)
	p.registerPrefix(token.PASSTHRU, p.parseIdentifier)
	p.registerPrefix(token.SYSTEM, p.parseIdentifier)
	p.registerPrefix(token.FUNCTION, p.parseFunctionDeclaration)

	p.infixParseFns = make(map[token.Kind]infixParseFn)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(kind token.Kind, fn prefixParseFn) {
	p.prefixParseFns[kind] = fn
}

func (p *Parser) registerInfix(kind token.Kind, fn infixParseFn) {
	p.infixParseFns[kind] = fn
}


func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Stmts: []ast.Stmt{}}
	if p.curTok.Kind != token.OPEN_TAG {
		return program
	}

	p.nextToken()
	
	for p.curTok.Kind != token.EOF && p.curTok.Kind != token.CLOSE_TAG {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
	
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curTok.Kind {
	case token.FUNCTION:
		if expr := p.parseExpression(); expr != nil {
			if stmt, ok := expr.(ast.Stmt); ok {
				return stmt
			}
		}
	default:
		return p.parseExpressionStatement()
	}
	return nil
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curTok}
	stmt.Expression = p.parseExpression()

	if p.peekTok.Kind == token.SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression() ast.Expr {
	prefix := p.prefixParseFns[p.curTok.Kind]
	if prefix == nil {
		return nil // No prefix parsing function found
	}
	leftExp := prefix()

	// This is the core of the Pratt parser for infix operators (like function calls)
	for p.peekTok.Kind != token.SEMICOLON {
		infix := p.infixParseFns[p.peekTok.Kind]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.curTok, Value: p.curTok.Lexeme}
}

func (p *Parser) parseCallExpression(function ast.Expr) ast.Expr {
	expr := &ast.CallExpr{Token: p.curTok, Function: function}
	expr.Arguments = p.parseCallArguments()
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expr {
	args := []ast.Expr{}
	if p.peekTok.Kind == token.RPAREN {
		p.nextToken()
		return args
	}
	p.nextToken()
	// For now, we don't parse arguments, just acknowledge the call
	for p.curTok.Kind != token.RPAREN && p.curTok.Kind != token.EOF {
		// A real implementation would call p.parseExpression() here
		p.nextToken()
	}
	return args
}

func (p *Parser) parseFunctionDeclaration() ast.Expr {
	stmt := &ast.FunctionDeclStmt{Token: p.curTok}

	if p.peekTok.Kind != token.IDENT {
		return nil // Invalid function declaration
	}
	p.nextToken()
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Lexeme}
	
    // We don't need to parse the body `() {}` for stubs, 
    // just skip until the next statement for performance.
	for p.curTok.Kind != token.SEMICOLON && p.curTok.Kind != token.RBRACE && p.curTok.Kind != token.EOF {
		p.nextToken()
	}
	
	return stmt
}