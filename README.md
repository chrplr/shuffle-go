# Shuffle-Go

*See an HTML formatted version of this document [here](https://chrplr.github.io/shuffle-go)*


A Go implementation of the `[shuffle](https://github.com/chrplr/shuffle)` program, providing a core library and both CLI and GUI interfaces for generating randomized sequences with sequential constraints.

This tool is particularly useful for generating stimuli lists for psychological experiments or any task requiring quasi-randomized permutations that avoid specific repetitions or patterns.

## Features

- **GUI App (Gio)**: A lightweight desktop application with native file dialogs (`shuffle-gio`). Recommended for most users.
- **GUI App (Fyne)**: An alternative desktop application built with Fyne (`shuffle-gui`), for users who prefer its interface.
- **CLI Tool**: A command-line interface for power users — scriptable, pipeable, and easy to integrate into automated workflows (`shuffle-cli`).
- **Core Library**: A reusable Go package for constraint-based shuffling, embeddable in your own projects.
- **Constraint Support**:
    - **Max Repetitions**: Limit consecutive occurrences of the same label in a column.
    - **Min Gap**: Ensure a minimum distance between identical labels in a column.
- **Dual Algorithms**:
    - **Constructive**: Fast, line-by-line building (standard).
    - **Equiprobable**: Brute-force filtering to ensure every valid permutation is equally likely.

## Installation

The easiest way to get Shuffle-Go is to download the latest installer directly below.

### 1. Using the Installers (Recommended for most users)

These are standard installers that will set up the application on your computer:

- **Windows**: Download [shuffle-windows-x86_64-setup.exe](https://github.com/chrplr/shuffle-go/releases/latest/download/shuffle-windows-x86_64-setup.exe). Run it to install Shuffle-Go. It will create a desktop shortcut for the GUI version and install the command-line tool as well.
- **macOS**: Download [shuffle-macos-arm64-app.zip](https://github.com/chrplr/shuffle-go/releases/latest/download/shuffle-macos-arm64-app.zip) (Apple Silicon/M1/M2/M3) or [shuffle-macos-x86_64-app.zip](https://github.com/chrplr/shuffle-go/releases/latest/download/shuffle-macos-x86_64-app.zip) (Intel). Extract the archive and drag **Shuffle-Go.app** to your **Applications** folder (or anywhere you like).

  > [!WARNING]
  > macOS may show a security warning the first time you open the app. See [macOS installation and security](https://chrplr.github.io/note-about-macos-unsigned-apps) for an explanation and step-by-step instructions to bypass it.

- **Linux**: Download [shuffle-linux-x86_64.AppImage](https://github.com/chrplr/shuffle-go/releases/latest/download/shuffle-linux-x86_64.AppImage). Right-click the file, go to **Properties > Permissions**, and check **"Allow executing file as program"**. You can then double-click it to run.

### 2. Using Pre-compiled Binaries (Portable version)

If you don't want to install the app, download a `.zip` archive from the [Releases page](https://github.com/chrplr/shuffle-go/releases) for your platform.

1. Download the `.zip` file for your Operating System and Architecture.
2. Extract the archive to a folder of your choice.
3. You will find:
    - `shuffle-gio`: The recommended graphical interface (native file dialogs).
    - `shuffle-gui`: An alternative graphical interface (Fyne version).
    - `shuffle-cli`: The command-line interface for scripting and terminal use.

### 3. Compiling from Source (For developers)

If you have [Go](https://golang.org/doc/install) installed, you can build the project yourself.

#### Prerequisites
- **Go 1.24** or later.
- **C Compiler**: Required for building the GUI components (e.g., GCC on Linux/Windows, or Xcode Command Line Tools on macOS).

#### Building
Run the provided build script:
```bash
bash build.sh
```
Or build specific components manually:
```bash
# Build the CLI
go build -o shuffle-cli ./cmd/shuffle-cli

# Build the Gio GUI
go build -tags novulkan -o shuffle-gio ./cmd/shuffle-gio

# Build the Fyne GUI
go build -o shuffle-gui ./cmd/shuffle-gui
```


## Usage

### GUI Application (Gio version — recommended)

Run the Gio GUI by double-clicking it or from the terminal:
```bash
./shuffle-gio      # on Linux and macOS
./shuffle-gio.exe  # on Windows
```

The Gio version uses your OS's native file dialogs:
1. Click **Import** to load a `.txt` or `.csv` file.
2. Adjust constraints and parameters in the sidebar.
3. Click **Shuffle!** to process the data.

### GUI Application (Fyne version — alternative)

Run the Fyne GUI by double-clicking it (on macOS, use `shuffle-gui.app`) or from the terminal:
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

### CLI Tool (power users)

The CLI is ideal for scripting and automation — it reads from a file or stdin and writes to stdout, making it easy to integrate into shell scripts or experiment generation pipelines.

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

# Use in a pipeline
cat sample.txt | ./shuffle-cli -n 10

# Reproducible shuffle with a fixed seed
./shuffle-cli -c "1 -4" -s 42 sample.txt
```

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
import "github.com/chrplr/shuffle-go"

// ... load data ...
constraints := []shuffle.Constraint{1, -3}
shuffler := shuffle.NewShuffler(data, constraints, 0, 100, 0)
result, err := shuffler.ShuffleConstructive()
```

## License

This project is copyrighted by its author, Christophe Pallier <christophe@pallier.org>

It is licensed under the GPLv3.
