package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Unquabain/thing-namer/internal/mcpsrv"
	"github.com/Unquabain/thing-namer/internal/namer"
	"github.com/Unquabain/thing-namer/internal/server"
)

//go:embed data/words.yaml
var wordsYAML []byte

func main() {
	n, err := namer.New(wordsYAML)
	if err != nil {
		log.Fatalf("load words: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		runStdio(n)
		return
	}

	runHTTP(n)
}

func runStdio(n *namer.Namer) {
	mcpServer := mcpsrv.NewServer(n)
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("mcp stdio: %v", err)
	}
}

func runHTTP(n *namer.Namer) {
	mux := http.NewServeMux()

	// Web routes (HTML/JSON/Go) — wrapped in CORS middleware inside Handler().
	mux.Handle("/", server.New(n).Handler())

	// MCP streamable HTTP — mounted directly, NOT wrapped in the referrer
	// CORS middleware. The MCP transport manages its own protocol-level
	// framing and Origin handling.
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpsrv.NewServer(n)
	}, nil)
	mux.Handle("/mcp", mcpHandler)
	mux.Handle("/mcp/", mcpHandler)

	log.Println("thing-namer listening on :9099")
	if err := http.ListenAndServe(":9099", mux); err != nil {
		log.Fatal(err)
	}
}
