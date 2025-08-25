package linter

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/config"
	"github.com/codevault-llc/php-lint/internal/lexer"
	"github.com/codevault-llc/php-lint/internal/parser"
	"github.com/codevault-llc/php-lint/internal/rules"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/codevault-llc/php-lint/pkg/types"
	"github.com/rs/zerolog"
)

var presetConfigs embed.FS

type Linter struct {
	config config.Config
	logger zerolog.Logger
	rules  []rules.Rule
	symbolTable *stubs.SymbolTable
}

func New(configPath string, logger zerolog.Logger) (*Linter, error) {
	cfg := config.New()
	if cfg == nil {
		return nil, fmt.Errorf("could not load configuration")
	}

	// Parse stubs
	symbolTable := stubs.NewSymbolTable()
	logger.Debug().Strs("stubs", cfg.Stubs).Msg("Parsing stubs to build symbol table")
	for _, stubPath := range cfg.Stubs {
		err := filepath.Walk(stubPath, func(p string, info os.FileInfo, err error) error {
			if err != nil { return err }
			if !info.IsDir() && filepath.Ext(p) == ".php" {
				content, _ := os.ReadFile(p)
				lxr := lexer.New(string(content))
				psr := parser.New(lxr)
				program := psr.ParseProgram()
				symbolTable.AddSymbolsFromAST(program)
			}
			return nil
		})
		if err != nil {
			logger.Warn().Err(err).Str("path", stubPath).Msg("Failed to walk stub path")
		}
	}

	logger.Debug().Int("functions", symbolTable.FunctionCount()).Msg("Symbol table built")

	// Collect active rules
	activeRules := []rules.Rule{}
	for ruleName, enabled := range cfg.Rules {
		if enabled {
			found := false
			for _, rule := range rules.GetRegistered() {
				if rule.Name() == ruleName {
					activeRules = append(activeRules, rule)
					found = true
					break
				}
			}
			if !found {
				logger.Warn().Str("rule", ruleName).Msg("Configured rule not found in registry")
			}
		}
	}
	logger.Info().Int("count", len(activeRules)).Msg("Active rules loaded")

	return &Linter{
		config:      *cfg,
		logger:      logger,
		rules:       activeRules,
		symbolTable: symbolTable,
	}, nil
}

func (l *Linter) LintProject(files map[string][]byte) map[string][]types.Issue {
	parsedFiles := make(map[string]*ast.Program)

	// Create and populate the AST for each file
	for filename, content := range files {
		lxr := lexer.New(string(content))
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		// Add locally defined functions
		l.symbolTable.AddSymbolsFromAST(program)
		parsedFiles[filename] = program
	}
	l.logger.Info().Int("functions", l.symbolTable.FunctionCount()).Msg("Symbol collection complete")

	results := make(map[string][]types.Issue)

	for filename, content := range files {
		program := parsedFiles[filename]
		var allIssues []types.Issue

		if enabled, ok := l.config.Rules["require-tags"]; ok && enabled {
			if !strings.HasPrefix(string(content), "<?php") {
				allIssues = append(allIssues, types.Issue{
					RuleName: "require-tags",
					Message:  "File must contain both '<?php' and '?>' tags.",
				})

				//return allIssues
			}
		}
		
		for _, rule := range l.rules {
			issues := rule.Check(filename, content, program, l.symbolTable)
			allIssues = append(allIssues, issues...)
		}
		results[filename] = allIssues
	}

	return results
}

func (l *Linter) Config() *config.Config {
	return &l.config
}