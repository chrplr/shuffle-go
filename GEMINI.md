# Gemini Project Context: Shuffle-Go

## Project Overview
**Shuffle-Go** is a Go-based implementation of a constraint-based shuffling utility. It allows users to randomize lists while adhering to specific sequential constraints, such as limiting consecutive repetitions of the same item or ensuring a minimum gap between identical items. This is particularly useful for generating experimental stimuli in psychological research.

The project is structured as follows:
- `shuffle/`: Core library containing the shuffling logic, data structures (`Shuffler`, `Constraint`), and file loading utilities.
- `cmd/shuffle-cli`: A command-line interface for the shuffler.
- `cmd/shuffle-gui`: A graphical user interface built using the [Fyne](https://fyne.io/) toolkit.
- `internal/version`: Internal package for managing the application version, injected at build time.

### Key Technologies
- **Language**: Go (1.24.6)
- **GUI Toolkit**: Fyne v2
- **License**: GPLv3
- **CI/CD**: GitHub Actions for multi-platform builds (Linux, Windows, MacOS).

---

## Building and Running

### Build All Binaries
The project includes a `build.sh` script that handles dependency synchronization (`go mod tidy`) and compiles both the CLI and GUI versions with version injection.

```bash
bash build.sh
```

### Build Manually
To build specific components without the script:

**CLI:**
```bash
go build -o shuffle-cli ./cmd/shuffle-cli
```

**GUI:**
```bash
go build -o shuffle-gui ./cmd/shuffle-gui
```

### Running Tests
The project uses standard Go testing.

```bash
go test ./...
```

---

## Development Conventions

### Versioning
The application version is managed via a Git tag. During the build process (via `build.sh` or GitHub Actions), the version is injected into the `github.com/chrplr/shuffle-go/internal/version.Version` variable using `-ldflags`.

### Constraints Logic
- **Positive `n`**: Maximum consecutive repetitions of an item in a column.
- **Negative `m`**: Minimum gap (number of intervening items) between identical items in a column.
- **Zero**: No constraint.

### Algorithm Choices
- **Constructive (Default)**: A fast, greedy algorithm that swaps items to satisfy constraints.
- **Equiprobable**: A filter-based approach that generates full permutations and checks them, ensuring every valid permutation has an equal probability of being selected (slower for large datasets or tight constraints).

### File Structure & Imports
All internal imports must use the module path `github.com/chrplr/shuffle-go`.

### Coding Style
- Follow standard Go idioms and `gofmt` formatting.
- All `.go` files must include the GPLv3 copyright header for Christophe Pallier.
