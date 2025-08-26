package linter

import (
	"embed"
	"os"
	"path/filepath"

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
}

func New(configPath string, logger zerolog.Logger) (*Linter, error) {
	cfg := config.New()
	if cfg == nil {
		logger.Error().Msg("Failed to create default config")
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
		config: *cfg,
		logger: logger,
		rules:  activeRules,
	}, nil
}

func (l *Linter) LintFile(path string, content []byte, symbolTable *stubs.SymbolTable) []types.Issue {
	lxr := lexer.New(string(content))
	psr := parser.New(lxr)
	program := psr.ParseProgram()

	var allIssues []types.Issue

	for _, rule := range l.rules {
		issues := rule.Check(path, content, program, symbolTable)
		allIssues = append(allIssues, issues...)
	}

	return allIssues
}

func (l *Linter) LintFiles(paths []string, symbolTable *stubs.SymbolTable) []types.Issue {
	var allIssues []types.Issue
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			l.logger.Error().Err(err).Str("path", path).Msg("Failed to read file")
			continue
		}
		issues := l.LintFile(path, content, symbolTable)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

func (l *Linter) Config() *config.Config {
	return &l.config
}