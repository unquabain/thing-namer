# thing-namer
Names things like they're action movies from the mid 90s.

## Installing

```bash
go install github.com/Unquabain/thing-namer@latest
```

## Building from source

```bash
go build
```

You now have an executable called `thing-namer` (or `thing-namer.exe`) in the current directory.

## Running

### Web server (default)

```bash
thing-namer
```

Starts an HTTP server on `:9099` that serves a rotating, themed project-name page. Routes:

- `GET /` — HTML page with a randomized color theme
- `GET /api.json` — JSON `{ "projectName": ..., "intro": ..., "outro": ... }`
- `GET /api.html` — Same as `/` but explicit
- `GET /api.go` — A small Go client that calls `/api.json` at the current host
- `POST /mcp` — Model Context Protocol streamable HTTP endpoint (see below)

### Stdio MCP server

```bash
thing-namer mcp
```

Speaks the Model Context Protocol over stdin/stdout. No HTTP listener is started.

## Using as an MCP server

`thing-namer` exposes its name generator as MCP tools so LLM clients can call it directly.

### Tools

| Tool | Parameters | Returns |
|---|---|---|
| `suggest_name` | _(none)_ | `{ "name": "Wizard Bacon" }` |
| `suggest_names` | `count` (integer, 1–100, required) | `{ "names": ["Wizard Bacon", "Sparkle Dolphin", ...] }` |

### Claude Code setup

**Stdio (local binary):**

```bash
go install github.com/Unquabain/thing-namer@latest
claude mcp add thing-namer thing-namer mcp
```

**HTTP (hosted instance):**

```bash
claude mcp add --transport http thing-namer https://wizard-bacon.unquabain.com/mcp
```

Once registered, Claude Code can call the tools whenever you ask for project name suggestions — for example, _"suggest a project name"_ or _"give me 5 codename options for this project"_.

## History

This is a re-write of a very simple program I wrote a long time ago when my team was having trouble naming things. The original was a JavaScript SPA, and then it was a Python web service.
