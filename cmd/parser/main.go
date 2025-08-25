package main

import (
	"fmt"

	"github.com/codevault-llc/php-lint/internal/lexer"
	"github.com/codevault-llc/php-lint/internal/parser"
)

func main() {
	// The input string has been corrected to 'Hello, World!' to match the lexer's simple string parsing
	input := "<?php echo 'Hello, World!'; ?>"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	// Print the string representation of the entire program AST
	fmt.Println(program.String())
}