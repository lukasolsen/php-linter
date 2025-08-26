package main

import (
	"log"
	"net/url"

	"github.com/codevault-llc/php-lint/internal/linter"
	"github.com/codevault-llc/php-lint/internal/stubs"
	"github.com/codevault-llc/php-lint/internal/workspace"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)


const lsName = "php-linter"

var version string = "0.0.1"
var handler protocol.Handler

var linterInstance *linter.Linter
var workspaceInstance *workspace.Workspace
var serverLogger commonlog.Logger

func main() {
	commonlog.Configure(1, nil)
	serverLogger = commonlog.GetLogger("php-linter")

	var err error
	linterInstance, err = linter.New("config.json", serverLogger)
	if err != nil {
		serverLogger.Criticalf("Failed to create linter: %v", err)
	}

	handler = protocol.Handler{
		Initialize:          onInitialize,
		Initialized:         onInitialized,
		Shutdown:            onShutdown,
		TextDocumentDidOpen: onDidOpen,
		TextDocumentDidChange: onDidChange,
		SetTrace:            setTrace,
	}

	server := server.NewServer(&handler, lsName, false)
	server.RunStdio()
}

func onInitialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	if params.RootURI != nil {
		uri, err := url.Parse(*params.RootURI)
		if err == nil {
			stubsTable := stubs.NewSymbolTable()
			// You would parse configured stubs here and pass them to the workspace
			// stubsTable.AddFromPath("/path/to/wordpress-stubs")

			workspaceInstance = workspace.New(uri.Path, stubsTable, serverLogger)
			
			go workspaceInstance.Build()
		}
	}

	capabilities := handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: &protocol.True,
		Change: func() *protocol.TextDocumentSyncKind { kind := protocol.TextDocumentSyncKindFull; return &kind }(),
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func onInitialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	log.Println("LSP server initialized")

	return nil
}

func onShutdown(ctx *glsp.Context) error {
	log.Println("Shutting down LSP server...")

	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func onDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	return lintDocument(ctx, params.TextDocument.URI, []byte(params.TextDocument.Text))
}

func onDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Println("Document changed, re-linting...")

	text := ""
	for _, change := range params.ContentChanges {
		if whole, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			text += whole.Text
		}
	}

	return lintDocument(ctx, params.TextDocument.URI, []byte(text))
}

func lintDocument(ctx *glsp.Context, uri string, text []byte) error {
	if workspaceInstance == nil {
		serverLogger.Warning("Workspace not initialized, cannot lint.")
		return nil
	}
	
	path, err := url.Parse(uri)
	if err != nil {
		return nil
	}

	// 1. Update the workspace with the latest file content from the editor
	workspaceInstance.UpdateFile(path.Path, text)

	// 2. Lint the file using the complete, up-to-date symbol table from the workspace
	issues := linterInstance.LintFile(path.Path, text, workspaceInstance.GetSymbolTable())

	diagnostics := []protocol.Diagnostic{}
	for _, issue := range issues {
		positionLine := protocol.UInteger(issue.Pos.Line - 1)
		positionColumn := protocol.UInteger(issue.Pos.Col - 1)

		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: positionLine, Character: positionColumn},
				End:   protocol.Position{Line: positionLine, Character: positionColumn},
			},
			Severity:  issue.Severity,
			Message:   issue.Message,
			Source:    issue.Source,
			CodeDescription: &protocol.CodeDescription{
				HRef: "https://example.com/rules/line-length",
			},
			RelatedInformation: nil,
			Data:               nil,
		})
	}
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	return nil
}

func splitLines(text string) []string {
	var lines []string
	start := 0
	for i, r := range text {
		if r == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}
	lines = append(lines, text[start:])
	return lines
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}