# Shuffle-Go

A Go implementation of the `shuffle` program, providing a core library and both CLI and GUI interfaces for generating randomized sequences with sequential constraints.

This tool is particularly useful for generating stimuli lists for psychological experiments or any task requiring quasi-randomized permutations that avoid specific repetitions or patterns.

## Features

- **Core Library**: A reusable Go package (`pkg/shuffle`) for constraint-based shuffling.
- **CLI Tool**: A powerful command-line interface for batch processing.
- **GUI App**: An interactive desktop application built with Fyne.
- **Constraint Support**:
    - **Max Repetitions**: Limit consecutive occurrences of the same label in a column.
    - **Min Gap**: Ensure a minimum distance between identical labels in a column.
- **Dual Algorithms**:
    - **Constructive**: Fast, line-by-line building (standard).
    - **Equiprobable**: Brute-force filtering to ensure every valid permutation is equally likely.

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (1.21 or later recommended)
- For the GUI: C compiler and development headers for your graphics driver (required by Fyne).

### Building

```bash
cd go-shuffle
go mod tidy

# Build CLI
go build -o shuffle-cli cmd/shuffle-cli/main.go

# Build GUI
go build -o shuffle-gui cmd/shuffle-gui/main.go
```

## Usage

### CLI Tool

```bash
./shuffle-cli [flags] [input_file]
```

**Flags:**
- `-c string`: Constraints as space-separated numbers (e.g., `"1 2 -3"`).
- `-d string`: Field delimiter (defaults to whitespace).
- `-e`: Use equiprobable shuffle algorithm (slower).
- `-i int`: Max iterations/loops for the search.
- `-n int`: Limit output to `n` lines.
- `-s int`: Random seed for reproducibility (0 for random).

**Example:**
```bash
# Shuffle sample.txt, max 1 repetition in col 1, output 10 lines
./shuffle-cli -c "1" -n 10 sample.txt
```

### GUI Application

Run the GUI with:
```bash
./shuffle-gui
```
The GUI allows you to:
1. Load data from `.txt` or `.csv` files.
2. Interactively set constraints and shuffling parameters.
3. Preview the results in a text area.
4. Save the shuffled list to a new file.

## Constraints Explained

Constraints are defined per column:
- **Positive number `n`**: A label in this column cannot be repeated more than `n` times consecutively.
- **Negative number `-m`**: Identical labels in this column must be separated by at least `m` lines.
- **Zero `0`**: No constraints for this column.

**Example: `-c "1 -4"`**
- Column 1: No adjacent identical labels.
- Column 2: Identical labels must have at least 4 other lines between them.

## Library Integration

You can use the core logic in your own Go projects:

```go
import "github.com/chrplr/shuffle-go/pkg/shuffle"

// ... load data ...
constraints := []shuffle.Constraint{1, -3}
shuffler := shuffle.NewShuffler(data, constraints, 0, 100, 0)
result, err := shuffler.ShuffleConstructive()
```

## License

This project is licensed under the same terms as the original shuffle program (GPL).
