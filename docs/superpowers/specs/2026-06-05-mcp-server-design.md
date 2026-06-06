# MCP Server + Comprehensive Unit Tests

**Date:** 2026-06-05
**Branch:** `mcp-server`

## Goal

Expose the project-name generator as an MCP (Model Context Protocol) server so LLM clients can call it as a tool. Add comprehensive unit tests across the existing and new code.

## Non-goals

- Auth / rate limiting on the MCP transport.
- Persisting or seeding randomness across calls.
- Renaming or relocating the existing routes (`/`, `/api.json`, `/api.html`, `/api.go`).
- Multi-binary split.

## Architecture

Extract three internal packages from the current single-file layout. `main.go` becomes a thin command dispatcher.

```
thing-namer/
├── main.go                          # subcommand dispatch
├── internal/
│   ├── namer/
│   │   ├── wordlist.go              # WeightedWord, WordList, Choose, YAML loading
│   │   ├── namer.go                 # Namer: Suggest(), SuggestN(n)
│   │   ├── theme.go                 # RenderContext, createTheme
│   │   └── *_test.go
│   ├── server/
│   │   ├── server.go                # routes, render handlers (HTML/JSON/Go)
│   │   ├── cors.go                  # ReferrerCORSMiddleware
│   │   └── *_test.go
│   └── mcpsrv/
│       ├── server.go                # NewServer(*namer.Namer) *mcp.Server
│       └── server_test.go
├── data/words.yaml                  # unchanged
├── templates/                       # unchanged
```

### Why this split

- `namer` is pure logic. It owns `*rand.Rand` so tests can seed it. No I/O beyond YAML loading from an `[]byte`.
- `server` owns HTTP-specific concerns (handlers, middleware, template execution). It depends on `namer` but knows nothing about MCP.
- `mcpsrv` owns MCP tool registration. It depends on `namer` but is transport-agnostic — the same `*mcp.Server` is used for stdio and HTTP.
- `main.go` wires everything: parses the subcommand, constructs `*namer.Namer` (loading the embedded YAML once), then either runs the web server (with MCP mounted at `/mcp`) or runs the stdio MCP server.

## Commands

| Invocation | Behavior |
|---|---|
| `thing-namer` | Web server on `:9099`. Existing routes unchanged. MCP added at `POST /mcp` via streamable HTTP transport. |
| `thing-namer mcp` | Stdio MCP server. No HTTP listener. |

The default behavior (no args → web server) is preserved exactly. Adding `mcp` as a subcommand is backwards-compatible.

## MCP tools

Both tools take no auth and call into the same `*namer.Namer`.

### `suggest_name`

- **Description:** "Suggests a single random project name in the style of a mid-90s action movie."
- **Input schema:** `{}` (no parameters)
- **Output:** `{ "name": "Wizard Bacon" }`

### `suggest_names`

- **Description:** "Suggests multiple random project names."
- **Input schema:** `{ "count": integer, 1 ≤ count ≤ 100 }`
- **Output:** `{ "names": ["Wizard Bacon", "Sparkle Dolphin", ...] }`
- **Validation:** `count` is required. Values outside `[1, 100]` return an MCP tool error with a clear message. The cap exists because the HTTP transport is internet-exposed.

### CORS interaction

The `/mcp` endpoint must NOT pass through `ReferrerCORSMiddleware` — the MCP streamable HTTP transport manages its own protocol-level framing and Origin handling. Mount it on a sub-mux or register the route ahead of the middleware wrapper.

## Dependencies

Add: `github.com/modelcontextprotocol/go-sdk` (official Go MCP SDK).

Use:
- `mcp.NewServer` for the server instance.
- `mcp.AddTool` for tool registration with typed input/output structs.
- `mcp.StdioTransport` for the `mcp` subcommand.
- `mcp.NewStreamableHTTPHandler` mounted at `/mcp` for the HTTP transport.

## Testing strategy

### `internal/namer`

- `WordList.Choose` returns a word from the list (table-driven, seeded `*rand.Rand`).
- `WordList.Choose` weight distribution: with a seeded RNG and many iterations, frequencies skew toward higher-weighted words (loose statistical assertion — e.g., heaviest word > 2× lightest over N=10000).
- `WordList.Add` concatenates word lists without mutating inputs.
- `WordList.UnmarshalYAML` round-trips a representative YAML document.
- `Namer.Suggest` returns non-empty string composed of words from the loaded list.
- `Namer.SuggestN(n)` returns slice of length `n`, all non-empty. Duplicates are allowed (each suggestion is independent).
- `createTheme` produces colors parseable by `colorful.Hex` and yields adequate contrast between background and title text (e.g., delta-E above a threshold).

### `internal/server`

- `httptest`-based tests for each render path:
  - `GET /api.json` → 200, `Content-Type: application/json`, body parses, has a `ProjectName` field that is non-empty.
  - `GET /api.html` → 200, `Content-Type: text/html`, body contains the rendered name.
  - `GET /api.go` → 200, `Content-Type: text/plain` (or similar), body contains valid Go syntax (parse with `go/parser`).
- `ReferrerCORSMiddleware` matrix:
  - `OPTIONS` with `Origin` → 204 with ACAO/ACAM/ACAH set.
  - `OPTIONS` without `Origin` → 204 (current behavior; documented in issue #5 for later fix).
  - `GET` with no `Referer` → handler runs, `X-Error: missing referrer` set, no ACAO.
  - `GET` with non-`http(s)://` `Referer` → handler runs, `X-Error: invalid referrer`, no ACAO.
  - `GET` with valid `Referer: https://example.com/foo` → ACAO is `https://example.com`.

### `internal/mcpsrv`

- Use `mcp.NewInMemoryTransports` to wire a client+server in-process.
- `suggest_name` → result contains a non-empty name.
- `suggest_names` with `count=5` → result contains exactly 5 non-empty names.
- `suggest_names` with `count=0` → tool returns error.
- `suggest_names` with `count=101` → tool returns error.
- `suggest_names` with missing `count` → tool returns error (schema enforcement).

### `main`

Skipped. Pure dispatch with no logic worth testing.

### Determinism

`namer.Namer` accepts an optional `*rand.Rand` via a constructor option. The constructor takes the embedded `data/words.yaml` bytes and returns a fully-loaded `*Namer`:

```go
n, err := namer.New(wordsYAML)              // default: time-seeded RNG
n, err := namer.New(wordsYAML, namer.WithRand(r))  // fixed-seed for tests
```

All randomness inside `namer` flows through this `*rand.Rand` — no calls to the package-level `math/rand` functions. Tests use a fixed seed for reproducibility.

## Documentation

`README.md` must be updated as part of this PR:

1. Remove the stale CLI claims (`thing-namer` printing a name, `thing-namer -n 20`) — those don't reflect current behavior.
2. Document the actual subcommands: bare `thing-namer` runs the web server on `:9099`; `thing-namer mcp` runs the stdio MCP server.
3. Add a **"Using as an MCP server"** section with Claude Code setup for both transports:

   **Stdio (local binary):**
   ```bash
   go install github.com/Unquabain/thing-namer@latest
   claude mcp add thing-namer thing-namer mcp
   ```

   **HTTP (hosted instance):**
   ```bash
   claude mcp add --transport http thing-namer https://wizard-bacon.unquabain.com/mcp
   ```

   **Invocation from Claude Code:**
   Once registered, the LLM can call the tools directly when asked something like _"suggest a project name"_ or _"give me 5 codename options for this project"_. The tools available are `suggest_name` and `suggest_names`.

4. Briefly list the two tools and their parameters (mirroring the table in this spec).

## Out of scope (deferred)

These are real issues but not part of this PR:
- Issues #3 (`Vary: Origin`), #4 (empty-host Referer), #5 (OPTIONS swallowing), #6 (Add vs Set) — tracked separately.

## Acceptance criteria

1. `thing-namer` runs the existing web server unchanged; `/mcp` responds to MCP protocol requests.
2. `thing-namer mcp` runs an MCP server over stdio that responds to `tools/list` and `tools/call`.
3. Both transports expose `suggest_name` and `suggest_names` with identical behavior.
4. `go test ./...` passes with all the tests above.
5. The existing Docker build and deployment still work.
6. `README.md` reflects current behavior and documents Claude Code MCP setup for both transports.
