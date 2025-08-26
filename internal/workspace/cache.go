package workspace

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/codevault-llc/php-lint/internal/ast"
	"github.com/codevault-llc/php-lint/internal/lexer"
	"github.com/codevault-llc/php-lint/internal/parser"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/rs/zerolog"
)

// CacheEntry stores the parsed AST and metadata for a single file.
type CacheEntry struct {
	AST     *ast.Program
	ModTime time.Time
}

// Workspace holds the state for the entire project.
type Workspace struct {
	rootDir     string
	logger      zerolog.Logger
	cache       map[string]CacheEntry
	symbolTable *stubs.SymbolTable
	mu          sync.RWMutex // To protect concurrent access to cache and symbols
}

// New creates and initializes a new workspace for the given root directory.
func New(rootDir string, stubsTable *stubs.SymbolTable, logger zerolog.Logger) *Workspace {
	return &Workspace{
		rootDir:     rootDir,
		logger:      logger,
		cache:       make(map[string]CacheEntry),
		symbolTable: stubsTable, // Start with stubs (WordPress, etc.)
	}
}

// Build performs the initial scan of the entire workspace.
func (w *Workspace) Build() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.logger.Info().Str("path", w.rootDir).Msg("Building initial workspace cache...")
	startTime := time.Now()

	filepath.Walk(w.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".php" {
			content, _ := os.ReadFile(path)
			w.updateCacheEntry(path, content, info.ModTime())
		}
		return nil
	})

	// After caching all ASTs, build the symbol table from them
	w.rebuildSymbolTableFromCache()

	w.logger.Info().
		Int("files", len(w.cache)).
		Int("functions", w.symbolTable.FunctionCount()).
		Dur("duration", time.Since(startTime)).
		Msg("Workspace cache built.")
}

// UpdateFile is called by the LSP when a file changes. It's fast because it only re-parses one file.
func (w *Workspace) UpdateFile(path string, content []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.updateCacheEntry(path, content, time.Now())
	w.rebuildSymbolTableFromCache()
	w.logger.Debug().Str("file", path).Msg("Workspace updated for single file")
}

// GetSymbolTable provides thread-safe access to the complete symbol table.
func (w *Workspace) GetSymbolTable() *stubs.SymbolTable {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.symbolTable
}

// updateCacheEntry is an internal helper to parse and cache a file.
func (w *Workspace) updateCacheEntry(path string, content []byte, modTime time.Time) {
	lxr := lexer.New(string(content))
	psr := parser.New(lxr)
	program := psr.ParseProgram()
	w.cache[path] = CacheEntry{
		AST:     program,
		ModTime: modTime,
	}
}

// rebuildSymbolTableFromCache rebuilds the entire symbol table from the cached ASTs.
// This is very fast as it operates on in-memory data.
func (w *Workspace) rebuildSymbolTableFromCache() {
	// Reset local symbols, but keep external stubs
	w.symbolTable.ClearLocalSymbols()
	for _, entry := range w.cache {
		w.symbolTable.AddSymbolsFromAST(entry.AST)
	}
}

func (w *Workspace) GetPHPFiles() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var phpFiles []string
	for path := range w.cache {
		phpFiles = append(phpFiles, path)
	}
	return phpFiles
}