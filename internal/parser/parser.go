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
	p.registerPrefix(token.ECHO, p.parseIdentifier)
	p.registerPrefix(token.OPEN_TAG, func() ast.Expr { return nil })

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
			// Ensure we advance the token stream after parsing a statement.
			// Some parse functions may not consume tokens fully, so explicitly move to the next token
			// to avoid an infinite loop that grows the statements slice without bound.
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	var stmt ast.Stmt
	switch p.curTok.Kind {
	/*case token.FUNCTION:
		stmt = p.parseFunctionDeclaration()*/
	default:
		stmt = p.parseExpressionStatement()
	}
	return stmt
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
	return &ast.Identifier{
		Base:  ast.Base{S: p.curTok.Span.Start, E: p.curTok.Span.End},
		Token: p.curTok,
		Value: p.curTok.Lexeme,
	}
}

func (p *Parser) parseCallExpression(function ast.Expr) ast.Expr {
	expr := &ast.CallExpr{
		Base:     ast.Base{S: function.Pos()},
		Token:    p.curTok,
		Function: function,
	}

	expr.Arguments = p.parseCallArguments()
	expr.E = p.curTok.Span.End

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

/*
func (p *Parser) parseFunctionDeclaration() ast.Expr {
	stmt := &ast.FunctionDeclStmt{Token: p.curTok}

	// Expect the function name (an identifier)
	if !p.expectPeek(token.IDENT) {
		return nil // Not a valid function declaration
	}
	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Lexeme}

	// Expect the opening parenthesis for arguments
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// For now, we just skip over the arguments until we find the closing parenthesis.
	// A more advanced parser would parse each argument here.
	for p.curTok.Kind != token.RPAREN && p.curTok.Kind != token.EOF {
		p.nextToken()
	}

	// Expect the opening curly brace for the function body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Smartly skip the function body by tracking the nesting of curly braces.
	// This prevents the parser from stopping on a '}' inside a nested block.
	braceDepth := 1
	for braceDepth > 0 && p.curTok.Kind != token.EOF {
		p.nextToken()
		if p.curTok.Kind == token.LBRACE {
			braceDepth++
		}
		if p.curTok.Kind == token.RBRACE {
			braceDepth--
		}
	}
	
	return stmt
}*/