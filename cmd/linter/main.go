package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codevault-llc/php-lint/internal/linter"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/codevault-llc/php-lint/internal/workspace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var linterInstance *linter.Linter
var workspaceInstance *workspace.Workspace
var logger zerolog.Logger

func main() {
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var err error
	linterInstance, err = linter.New("config.json", logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create linter")
	}

	paths := linterInstance.Config().Paths
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			fmt.Println("Usage: php-lint [options] [paths...]")
			fmt.Println("Options:")
			fmt.Println("  --help, -h       Show this help message")
			return
		}

		paths = os.Args[1:]
	}
	logger.Debug().Strs("paths", paths).Msg("Determined target paths")

	absPath, err := filepath.Abs(paths[0])
	if err != nil {
		logger.Error().Err(err).Str("path", paths[0]).Msg("Failed to get absolute path")
		return
	}

	logger.Info().Str("path", absPath).Msg("Starting linting for path")

	// Workspace -- Init
	stubsTable := stubs.NewSymbolTable()
	workspaceInstance = workspace.New(absPath, stubsTable, logger)
	workspaceInstance.Build()

	// Run linter
	phpFiles := workspaceInstance.GetPHPFiles()

	logger.Info().Int("files", len(phpFiles)).Msg("Starting linting process")

	linterInstance.LintFiles(phpFiles, workspaceInstance.GetSymbolTable())

	//reporter.Render(results)
}