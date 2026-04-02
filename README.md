# Shuffle-Go

Shuffle-Go is a constraint-based list randomizer written in Go. It creates randomized sequences while respecting sequential constraints (e.g., max repetitions, minimum gap), making it ideal for experimental stimuli generation.

- **Constraint Support**: Limit consecutive repetitions (Max Repetitions) or ensure minimum distance between identical items (Min Gap).
- **Dual Algorithms**: Fast **Constructive** (greedy) or **Equiprobable** (brute-force filtering) for unbiased permutations.
- **Multiple Interfaces**: Native GUI (Gio), Cross-platform GUI (Fyne), and a scriptable CLI.

**[Read the Paper](paper/shuffle-paper.pdf)** | **[Download Releases](https://github.com/chrplr/shuffle-go/releases/latest)**

---

## Quick Start

### 1. Download & Install
Download the latest version for your platform from the [Releases page](https://github.com/chrplr/shuffle-go/releases/latest).

| Platform | Recommended File | Note |
| :--- | :--- | :--- |
| **Windows** | `...-setup.exe` | Standard installer; creates desktop shortcut. |
| **macOS** | `...-app.zip` | Extract and move `Shuffle-Go.app` to Applications. [See security note](https://chrplr.github.io/note-about-macos-unsigned-apps). |
| **Linux** | `.AppImage` | Make executable (`chmod +x`) and run. |

### 2. Usage
- **GUI (Recommended)**: Launch `shuffle-gio`. Import your `.txt` or `.csv`, set constraints, and click **Shuffle!**.
- **CLI**: Use `shuffle-cli` for scripts and automation.

---

## Constraints Syntax

Constraints are defined per column (space-separated):
- **`n` (Positive)**: Max consecutive repetitions of the same label.
- **`-m` (Negative)**: Min gap (intervening rows) between identical labels.
- **`0`**: No constraints for this column.

**Example**: `-c "1 -4"` means Column 1 has no adjacent identical labels, and Column 2 requires at least 4 items between repeats.

---

## CLI Tool Reference

```bash
./shuffle-cli [flags] < [input_file]
```

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-c` | Constraints string (e.g., `"1 2 -3"`) | `""` |
| `-d` | Field delimiter (whitespace, `,`, etc.) | whitespace |
| `-e` | Use equiprobable algorithm (slower) | false |
| `-n` | Limit output to `n` lines | all |
| `-s` | Random seed for reproducibility | 0 (random) |

---

## Library Integration

Use the core logic in your own Go projects:

```go
import "github.com/chrplr/shuffle-go"

shuffler := shuffle.NewShuffler(data, []shuffle.Constraint{1, -3}, seed, maxIter, limit)
result, err := shuffler.ShuffleConstructive()
```

---

## Development

Build all binaries using the provided script (requires Go 1.24+):
```bash
bash build.sh
```

## License
Copyright © Christophe Pallier. Licensed under the [GNU GPL v3](LICENSE.txt).
