package main

import (
	"log"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)


const lsName = "php-linter"

var version string = "0.0.1"
var handler protocol.Handler
func main() {
	log.Println("PHP Linter LSP is starting...")

	/*l, err := linter.New("config.json", log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize linter")
	}*/

	commonlog.Configure(1, nil)

	log.Println("PHP Linter LSP is starting... handler finished")

	handler = protocol.Handler{
		Initialize:  onInitialize,
		Initialized: onInitialized,
		Shutdown:    onShutdown,
		TextDocumentDidOpen:  onDidOpen,
		TextDocumentDidChange: onDidChange,
		SetTrace: setTrace,
	}

	log.Println("PHP Linter LSP is starting... server initialized")

	server := server.NewServer(&handler, lsName, false)

	log.Println("PHP Linter LSP is starting... before run stdio")

	server.RunStdio()

	log.Println("Running LSP server on stdio...")

	//results := l.LintProject(filesToLint)
}

func onInitialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
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
	log.Println("Document opened, linting...")

	return lintDocument(ctx, params.TextDocument.URI, params.TextDocument.Text)
}

func onDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Println("Document changed, re-linting...")

	// turn it into a string
	text := ""
	for _, change := range params.ContentChanges {
		if whole, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			text += whole.Text
		}
	}
	
	return lintDocument(ctx, params.TextDocument.URI, text)
}

func lintDocument(ctx *glsp.Context, uri string, text string) error {
	var diagnostics []protocol.Diagnostic


	lines := splitLines(text)
	log.Println("Linting document:", uri, len(lines), "lines")

	for i, line := range lines {
        if len(line) > 80 {
            sev := protocol.DiagnosticSeverityWarning
            diagnostics = append(diagnostics, protocol.Diagnostic{
                Range: protocol.Range{
                    Start: protocol.Position{Line: uint32(i), Character: 80},
                    End:   protocol.Position{Line: uint32(i), Character: uint32(len(line))},
                },
                Severity: &sev,
                Message:  "Line exceeds 80 characters",
            })
        }
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