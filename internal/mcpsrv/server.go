package mcpsrv

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Unquabain/thing-namer/internal/namer"
)

const maxCount = 100

type suggestNameInput struct{}

type suggestNameOutput struct {
	Name string `json:"name"`
}

type suggestNamesInput struct {
	Count int `json:"count" jsonschema:"Number of names to generate (1-100)."`
}

type suggestNamesOutput struct {
	Names []string `json:"names"`
}

// NewServer returns an *mcp.Server with suggest_name and suggest_names
// registered. The same server can be served over stdio or streamable HTTP.
func NewServer(n *namer.Namer) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "thing-namer",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "suggest_name",
		Description: "Suggests a single random project name in the style of a mid-90s action movie.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ suggestNameInput) (*mcp.CallToolResult, suggestNameOutput, error) {
		return nil, suggestNameOutput{Name: n.Suggest()}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "suggest_names",
		Description: "Suggests multiple random project names.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in suggestNamesInput) (*mcp.CallToolResult, suggestNamesOutput, error) {
		if in.Count < 1 || in.Count > maxCount {
			return nil, suggestNamesOutput{}, fmt.Errorf("count must be between 1 and %d, got %d", maxCount, in.Count)
		}
		return nil, suggestNamesOutput{Names: n.SuggestN(in.Count)}, nil
	})

	return s
}
