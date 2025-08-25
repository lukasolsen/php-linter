package stubs

import "github.com/codevault-llc/php-lint/internal/ast"


type SymbolTable struct {
	functions map[string]bool
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{functions: make(map[string]bool)}
}

func (st *SymbolTable) AddFunction(name string) {
	st.functions[name] = true
}

func (st *SymbolTable) IsFunctionDefined(name string) bool {
	_, exists := st.functions[name]
	return exists
}

func (st *SymbolTable) AddSymbolsFromAST(program *ast.Program) {
	for _, stmt := range program.Stmts {
		if funcDecl, ok := stmt.(*ast.FunctionDeclStmt); ok {
			st.AddFunction(funcDecl.Name.Value)
		}
	}
}

func (st *SymbolTable) FunctionCount() int {
	return len(st.functions)
}