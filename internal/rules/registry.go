package rules

import (
	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/pkg/types"
)

// Rule defines the interface for all linting rules.
// A community member would only need to implement this interface.
type Rule interface {
	Name() string
	Description() string
	Check(filename string, content []byte, program *ast.Program) []types.Issue
}

var registry = make(map[string]Rule)

// Register adds a new rule to the central registry.
// This is called by each rule file's init() function.
func Register(rule Rule) {
	registry[rule.Name()] = rule
}

// GetRegistered returns a slice of all registered rules.
func GetRegistered() []Rule {
	rules := make([]Rule, 0, len(registry))
	for _, rule := range registry {
		rules = append(rules, rule)
	}
	return rules
}

// init is a special Go function that runs when the package is imported.
// We use it to automatically register our core rules.
func init() {
	Register(&RuleNoEval{})
	Register(&RuleNoShellExec{})
	Register(&RuleUndefinedFunction{})
}