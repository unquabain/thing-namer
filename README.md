# thing-namer
Names things like they're action movies from the mid 90s.

## Installing

You know the drill.

```bash
go get github.com/Unquabain/thing-namer
```

## Building

The only dependency not in the standard library is `gopkg.in/yaml/v2`, which is pretty standard.

```bash
go build
```

You now have an executable called `thing-namer` or `thing-namer.exe` (depending on your system) in the current directory.

## Running

```bash
thing-namer
```

Prints 

```
Your project is now called "Wizard Bacon".
```

```
thing-namer -n 20
```

Prints twenty different suggestions.

## History

This is a re-write of a very simple program I wrote a long time ago when my team was having trouble naming things. The original was a JavaScript SPA, and then it was a Python web service.
