package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codevault-llc/php-lint/internal/linter"
	"github.com/codevault-llc/php-lint/internal/reporter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	l, err := linter.New("config.json", log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize linter")
	}

	paths := l.Config().Paths
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			fmt.Println("Usage: php-lint [options] [paths...]")
			fmt.Println("Options:")
			fmt.Println("  --help, -h       Show this help message")
			return
		}

		paths = os.Args[1:]
	}
	log.Debug().Strs("paths", paths).Msg("Determined target paths")

	filesToLint := make(map[string][]byte)
	for _, path := range paths {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(p) == ".php" {
				content, readErr := os.ReadFile(p)
				if readErr != nil {
					log.Error().Err(readErr).Str("file", p).Msg("Failed to read file")
					return nil
				}
				filesToLint[p] = content
			}
			return nil
		})
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Failed to walk path")
		}
	}

	results := l.LintProject(filesToLint)

	reporter.Render(results)
}