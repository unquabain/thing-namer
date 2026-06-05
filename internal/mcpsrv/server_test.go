package mcpsrv

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Unquabain/thing-namer/internal/namer"
)

func loadNamer(t *testing.T) *namer.Namer {
	t.Helper()
	wd, _ := os.Getwd()
	root := filepath.Join(wd, "..", "..")
	data, err := os.ReadFile(filepath.Join(root, "data", "words.yaml"))
	if err != nil {
		t.Fatalf("read words.yaml: %v", err)
	}
	n, err := namer.New(data, namer.WithRand(rand.New(rand.NewSource(1))))
	if err != nil {
		t.Fatalf("namer.New: %v", err)
	}
	return n
}

func connect(t *testing.T) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()

	srv := NewServer(loadNamer(t))
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil)

	clientT, serverT := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, serverT, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	sess, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { sess.Close() })
	return sess
}

func TestSuggestName(t *testing.T) {
	sess := connect(t)
	res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "suggest_name",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool reported error: %+v", res.Content)
	}
	if res.StructuredContent == nil {
		t.Fatal("StructuredContent nil")
	}
	out, ok := res.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("StructuredContent type = %T", res.StructuredContent)
	}
	name, _ := out["name"].(string)
	if name == "" {
		t.Fatalf("name empty in %+v", out)
	}
}

func TestSuggestNamesValid(t *testing.T) {
	sess := connect(t)
	res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "suggest_names",
		Arguments: map[string]any{"count": 5},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool reported error: %+v", res.Content)
	}
	out := res.StructuredContent.(map[string]any)
	names, _ := out["names"].([]any)
	if len(names) != 5 {
		t.Fatalf("got %d names, want 5", len(names))
	}
	for i, n := range names {
		if s, _ := n.(string); s == "" {
			t.Fatalf("names[%d] empty", i)
		}
	}
}

func TestSuggestNamesZeroRejected(t *testing.T) {
	sess := connect(t)
	res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "suggest_names",
		Arguments: map[string]any{"count": 0},
	})
	if err != nil {
		t.Fatalf("CallTool transport error: %v", err)
	}
	if !res.IsError {
		t.Fatalf("expected IsError=true for count=0, got %+v", res)
	}
}

func TestSuggestNamesOverCapRejected(t *testing.T) {
	sess := connect(t)
	res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "suggest_names",
		Arguments: map[string]any{"count": 101},
	})
	if err != nil {
		t.Fatalf("CallTool transport error: %v", err)
	}
	if !res.IsError {
		t.Fatalf("expected IsError=true for count=101, got %+v", res)
	}
}

func TestSuggestNamesMissingCountRejected(t *testing.T) {
	sess := connect(t)
	res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "suggest_names",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool transport error: %v", err)
	}
	if !res.IsError {
		t.Fatalf("expected IsError=true for missing count, got %+v", res)
	}
}
