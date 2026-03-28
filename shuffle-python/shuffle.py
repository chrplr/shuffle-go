# Time-stamp: <2021-12-09 11:59:09 christophe@pallier.org>
"""
Constraint-based list shuffler.

Each row is tokenized into fields (columns). Constraints are expressed as:
  - maxrep: maximum consecutive repetitions of the same label in a column
  - mingap: minimum gap (number of intervening rows) between repetitions
            of the same label in a column

Both are dicts mapping column index (0-based) to a number.
Example: maxrep={0: 2, 3: 1} means column 0 allows at most 2 consecutive
repetitions, and column 3 allows no consecutive repetitions.
"""

import csv
import random
import time
from typing import Optional


def load_csv(filename: str) -> list[list[str]]:
    """Load a table from a CSV file, auto-detecting the delimiter.

    Args:
        filename: path to the CSV file

    Returns:
        List of rows, each row being a list of strings.

    Raises:
        ValueError: if the delimiter cannot be detected.
    """
    with open(filename, "rt") as f:
        sample = f.readline()
        dialect = csv.Sniffer().sniff(sample, [',', ';', '\t', ' '])
        if dialect is None:
            raise ValueError(f"Cannot detect delimiter in {filename!r}")
        f.seek(0)
        return [row for row in csv.reader(f, dialect)]


def simple_shuffle(table: list) -> list:
    """Return a new list with the elements of *table* in a random order."""
    result = table.copy()
    random.shuffle(result)
    return result


def check_constraints(
    table: list[list[str]],
    maxrep: Optional[dict[int, int]] = None,
    mingap: Optional[dict[int, int]] = None,
    irow: int = 0,
) -> tuple[bool, int]:
    """Check whether *table* satisfies the given constraints from *irow* onward.

    Args:
        table: list of rows to check.
        maxrep: maps column index → max allowed consecutive repetitions.
        mingap: maps column index → min gap between identical labels.
        irow: starting row index (rows before this are assumed valid).

    Returns:
        (ok, row): *ok* is True if constraints are satisfied; *row* is the
        index of the first violating row (or len(table) when ok is True).
    """
    if irow < 0:
        irow = 0

    repetitions: dict[int, int] = {}
    if maxrep is not None:
        repetitions = {col: 1 for col in maxrep}

    previous = table[irow]
    irow += 1
    ok = True

    while ok and irow < len(table):
        row = table[irow]

        if maxrep is not None:
            for col, limit in maxrep.items():
                if previous[col] == row[col]:
                    repetitions[col] += 1
                    if repetitions[col] > limit:
                        ok = False
                        break
                else:
                    repetitions[col] = 1

        if ok and mingap is not None:
            for col, gap in mingap.items():
                back = max(0, irow - gap)
                while back < irow:
                    if table[back][col] == row[col]:
                        ok = False
                        break
                    back += 1
                if not ok:
                    break

        previous = row
        if ok:
            irow += 1

    return (ok, irow)


def shuffle_equiprob(
    table: list[list[str]],
    maxrep: Optional[dict[int, int]] = None,
    mingap: Optional[dict[int, int]] = None,
    time_limit: float = 1.0,
) -> Optional[list[list[str]]]:
    """Shuffle *table* so that every valid permutation is equally likely.

    Repeatedly shuffles the table at random until a permutation that satisfies
    the constraints is found, or until *time_limit* seconds have elapsed.

    Args:
        table: list of rows to shuffle (modified in place; pass a copy if
               the original must be preserved).
        maxrep: max consecutive repetitions per column.
        mingap: min gap between identical labels per column.
        time_limit: maximum wall-clock seconds to spend searching.

    Returns:
        The shuffled table, or None if no valid permutation was found in time.
    """
    deadline = time.time() + time_limit
    ok = False
    while not ok and time.time() < deadline:
        random.shuffle(table)
        ok, _ = check_constraints(table, maxrep, mingap)
    return table if ok else None


def shuffle_constructive(
    table: list[list[str]],
    maxrep: Optional[dict[int, int]] = None,
    mingap: Optional[dict[int, int]] = None,
    time_limit: float = 1.0,
) -> Optional[list[list[str]]]:
    """Shuffle *table* using a constructive (greedy) algorithm.

    Builds a valid permutation row by row, swapping the current violating row
    with a later one. Falls back to a full re-shuffle when stuck.

    Args:
        table: list of rows to shuffle (modified in place; pass a copy if
               the original must be preserved).
        maxrep: max consecutive repetitions per column.
        mingap: min gap between identical labels per column.
        time_limit: maximum wall-clock seconds to spend searching.

    Returns:
        The shuffled table, or None if no valid permutation was found in time.
    """
    if maxrep is None and mingap is None:
        raise ValueError("At least one of maxrep or mingap must be specified.")

    n = len(table)
    backtrack = max(
        max(mingap.values()) if mingap else 0,
        max(maxrep.values()) if maxrep else 0,
    )

    deadline = time.time() + time_limit
    random.shuffle(table)
    ok = False
    irow = 0
    nfailure = 0

    while not ok and time.time() < deadline:
        ok, irow = check_constraints(table, maxrep, mingap, irow - backtrack)
        if not ok:
            nfailure += 1
            if irow >= n - 1 or nfailure > n * 100:
                random.shuffle(table)
                irow = 0
                nfailure = 0
            else:
                i2 = random.choice(range(irow + 1, n))
                table[irow], table[i2] = table[i2], table[irow]

    return table if ok else None
