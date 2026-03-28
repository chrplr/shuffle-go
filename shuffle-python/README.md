# shuffle.py

A Python library for shuffling tabular data while respecting sequential constraints. Useful for generating randomized stimulus lists in psychological experiments.

## Requirements

Python 3.10+, no external dependencies.

## Constraints

Constraints are expressed as dicts mapping a **column index** (0-based) to a number:

- `maxrep` — maximum number of consecutive rows with the same label in that column.
  Example: `maxrep={0: 2}` means column 0 may not have the same value three times in a row.
- `mingap` — minimum number of intervening rows between two rows with the same label in that column.
  Example: `mingap={2: 3}` means at least 3 rows must separate any two rows sharing the same value in column 2.

Both can be specified simultaneously, and multiple columns can be constrained at once.

## API

### `load_csv(filename) → list[list[str]]`

Load a CSV file into a list of rows. The delimiter (`,`, `;`, `\t`, or space) is detected automatically.

```python
table = load_csv("stimuli.csv")
```

### `shuffle_constructive(table, maxrep, mingap, time_limit) → list[list[str]] | None`

Shuffle `table` using a fast greedy algorithm. Builds a valid ordering row by row, swapping the offending row with a later one when a constraint is violated, and restarts from scratch when stuck.

Returns the shuffled table, or `None` if no valid permutation was found within `time_limit` seconds (default: 1.0).

```python
result = shuffle_constructive(table, maxrep={0: 1}, time_limit=5.0)
if result is None:
    print("No valid ordering found — constraints may be too tight.")
```

### `shuffle_equiprob(table, maxrep, mingap, time_limit) → list[list[str]] | None`

Shuffle `table` so that every valid permutation has an equal probability of being selected. Repeatedly draws a random permutation and accepts it if it satisfies the constraints.

Slower than `shuffle_constructive`, especially for large tables or tight constraints. Use when unbiased sampling is required.

```python
result = shuffle_equiprob(table, mingap={2: 2}, time_limit=10.0)
```

### `simple_shuffle(table) → list`

Return a new list with the rows in a uniformly random order, with no constraints applied.

### `check_constraints(table, maxrep, mingap, irow) → tuple[bool, int]`

Check whether `table` satisfies the constraints starting from row `irow`. Returns `(True, len(table))` if all constraints are met, or `(False, i)` where `i` is the index of the first violating row.

## Examples

### Shuffle a CSV file ensuring no two consecutive rows share the same value in column 0

```python
import csv
from shuffle import load_csv, shuffle_constructive

table = load_csv("stimuli.csv")
result = shuffle_constructive(table, maxrep={0: 1})
if result is not None:
    with open("stimuli_shuffled.csv", "wt", newline="") as f:
        csv.writer(f).writerows(result)
```

### Generate multiple shuffled lists

```python
from shuffle import load_csv, shuffle_constructive

table = load_csv("stimuli.csv")
for i in range(10):
    # pass a copy so the original is not modified
    result = shuffle_constructive(table.copy(), mingap={2: 2})
    if result is not None:
        with open(f"list{i+1:02d}.csv", "wt", newline="") as f:
            import csv
            csv.writer(f).writerows(result)
```

### Set a fixed random seed for reproducibility

```python
import random
from shuffle import load_csv, shuffle_constructive

random.seed(42)
table = load_csv("stimuli.csv")
result = shuffle_constructive(table, maxrep={0: 2, 1: 1})
```

## Notes

- Both shuffle functions modify the table in place **and** return it. Pass `table.copy()` if you need to preserve the original.
- If `None` is returned, try relaxing the constraints or increasing `time_limit`. Some constraint combinations make a valid permutation impossible (e.g., more than half the rows sharing a label with `maxrep={0: 1}`).
- `shuffle_constructive` does not guarantee equal probability across all valid permutations; use `shuffle_equiprob` when that property matters.
