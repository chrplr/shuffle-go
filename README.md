# Shuffle-Go

*See an HTML formatted version of this document [here](https://chrplr.github.io/shuffle-go)*

A Go implementation of the `[shuffle](https://github.com/chrplr/shuffle)` program, providing a core library and both CLI and GUI interfaces for generating randomized sequences with sequential constraints.

This tool is particularly useful for generating stimuli lists for psychological experiments or any task requiring quasi-randomized permutations that avoid specific repetitions or patterns.

## Features

- **Core Library**: A reusable Go package (`pkg/shuffle`) for constraint-based shuffling.
- **CLI Tool**: A powerful command-line interface for batch processing.
- **GUI Apps**: 
    - **Fyne version**: A desktop application built with Fyne (`shuffle-gui`).
    - **Gio version**: A lightweight desktop application built with Gio (`shuffle-gio`), supporting native file import.
- **Constraint Support**:
    - **Max Repetitions**: Limit consecutive occurrences of the same label in a column.
    - **Min Gap**: Ensure a minimum distance between identical labels in a column.
- **Dual Algorithms**:
    - **Constructive**: Fast, line-by-line building (standard).
    - **Equiprobable**: Brute-force filtering to ensure every valid permutation is equally likely.

## Installation

The apps (`shuffle-cli`, `shuffle-gui` and `shuffle-gio`) can be downloaded from <https://github.com/chrplr/shuffle-go/releases>.

### macOS Security Note

Because these binaries are not signed by an Apple Developer account, macOS may prevent them from running.
You will need to explicitly authorize them to run in `System Settings > Privacy & Security` parameters.                                                                              

#### For the GUI apps (`shuffle-gui.app` and `shuffle-gio.app`):
1. **Right-click** (or Control-click) the `.app` icon.
2. Select **Open** from the shortcut menu.
3. Click **Open** in the dialog box that appears.


#### For the CLI (`shuffle-cli`):
Open a terminal in the folder containing the binary and run:
```bash
chmod +x shuffle-cli
xattr -d com.apple.quarantine shuffle-cli
```
Then you can run it normally: `./shuffle-cli`


## Compiling from source.

### Prerequisites

- [Go](https://golang.org/doc/install) (1.21 or later recommended)
- For the GUI: C compiler and development headers for your graphics driver (required by Fyne).

### Building

```bash
go mod tidy

# Build CLI
go build -o shuffle-cli ./cmd/shuffle-cli

# Build GUI based on Fyne
go build -o shuffle-gui ./cmd/shuffle-gui

# Build GUI based on Gio
go build -o shuffle-gui ./cmd/shuffle-gio

```




## Usage

### CLI Tool

```bash
./shuffle-cli [flags] < [input_file]
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
cat sample.txt | ./shuffle-cli -n 10
```

### GUI Application (Fyne version)

Run the GUI by double-clicking it (on macOS, use `shuffle-gui.app`) or from the terminal:
```bash
./shuffle-gui      # on Linux
./shuffle-gui.exe  # on Windows
./shuffle-gui.app  # on macOS
```

The Fyne version allows you to:
1. Load data from `.txt` or `.csv` files using the "Open" button.
2. Interactively set constraints and shuffling parameters.
3. Preview the results in a text area.
4. Save the shuffled list to a new file using the "Save" button.

### GUI Application (Gio version)

Run the Gio GUI by double-clicking it (on macOS, use `shuffle-gio.app`) or from the terminal:
```bash
./shuffle-gio      # on Linux
./shuffle-gio.exe  # on Windows
./shuffle-gio.app  # on macOS
```

The Gio version features a native file selector via the **Import** button.
1. Click **Import** to load a `.txt` or `.csv` file.
2. Adjust constraints and parameters in the sidebar.
3. Click **Shuffle!** to process the data.

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

This project is copyrighted by its author, Christophe Pallier <christophe@pallier.org>

It is licensed under the GPLv3.
