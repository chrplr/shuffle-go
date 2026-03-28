# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Does

**Shuffle-Go** is a constraint-based list randomizer written in Go. It shuffles multi-column data files while enforcing sequential constraints (e.g., max consecutive repetitions, minimum gap between identical items). Primary use case is generating experimental stimuli for psychological research.

Three binaries are produced:
- `shuffle-cli` — command-line interface
- `shuffle-gui` — GUI using Fyne toolkit
- `shuffle-gio` — GUI using Gio toolkit (lightweight alternative)

## Commands

### Build
```bash
bash build.sh          # build all binaries with version injection
go build -o shuffle-cli ./cmd/shuffle-cli
go build -o shuffle-gui ./cmd/shuffle-gui
go build -o shuffle-gio ./cmd/shuffle-gio
```

### Test
```bash
go test ./...                          # all tests
go test -run TestShuffleConstructive . # single test
```

### Format
```bash
gofmt -w .
```

## Architecture

**Core library** (`shuffle.go`, package `shuffle`):
- `Shuffler` struct holds input data, constraints, RNG, and iteration config
- `LoadData(reader, delimiter)` parses input into `[][]string`
- `NewShuffler(data, constraints, seed, maxIter, limit)` initializes the shuffler
- `CheckConstraints(data)` validates whether a sequence satisfies all constraints
- `ShuffleConstructive()` — fast greedy algorithm (default); swaps items line-by-line
- `ShuffleEquiprob()` — brute-force permutation filter; equal probability for all valid outputs (slow)

**Constraint encoding** (`type Constraint int`):
- Positive `n`: max consecutive repetitions of the same label in a column
- Negative `-m`: minimum gap (intervening rows) between identical labels in a column
- Zero: no constraint for that column

**CLI** (`cmd/shuffle-cli/main.go`): parses flags, reads stdin or file, calls core library, writes to stdout.

**GUIs** (`cmd/shuffle-gui/`, `cmd/shuffle-gio/`): file open/save dialogs, interactive constraint input, calls same core library.

**Version** (`internal/version/version.go`): `Version` variable set to empty string at compile time; injected via `-ldflags` from git tags during `build.sh` or CI.

**Module path:** `github.com/chrplr/shuffle-go` — use this for all internal imports.

## Coding Conventions

- All `.go` files must include the GPLv3 copyright header for Christophe Pallier.
- Follow standard Go idioms and `gofmt` formatting.
- Version is injected at build time via `-ldflags "-X github.com/chrplr/shuffle-go/internal/version.Version=<tag>"`.
