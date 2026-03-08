// Copyright (C) 2026 Christophe Pallier
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package shuffle

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

// LoadData reads lines from r and splits them by delimiter.
func LoadData(r io.Reader, delimiter string) ([][]string, error) {
	var data [][]string
	
	if delimiter != "" && len(delimiter) == 1 {
		reader := csv.NewReader(r)
		reader.Comma = rune(delimiter[0])
		reader.FieldsPerRecord = -1 // Allow variable number of fields
		reader.TrimLeadingSpace = true
		
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		data = records
	} else {
		// Fallback to manual scanner for multi-char delimiters or whitespace
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			var fields []string
			if delimiter == "" {
				fields = strings.Fields(line)
			} else {
				fields = strings.Split(line, delimiter)
			}
			data = append(data, fields)
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// Constraint represents a constraint on a specific column.
// Positive value: max consecutive repetitions.
// Negative value: min gap between identical labels.
// Zero: no constraint.
type Constraint int

// Shuffler holds the data and configuration for the shuffle operation.
type Shuffler struct {
	Data        [][]string
	Constraints []Constraint
	Rand        *rand.Rand
	MaxIter     int
	Limit       int // Number of lines to output
}

// NewShuffler creates a new Shuffler instance.
func NewShuffler(data [][]string, constraints []Constraint, seed int64, maxIter int, limit int) *Shuffler {
	source := rand.NewSource(seed)
	if seed == 0 {
		source = rand.NewSource(time.Now().UnixNano())
	}
	
	if limit <= 0 || limit > len(data) {
		limit = len(data)
	}

	return &Shuffler{
		Data:        data,
		Constraints: constraints,
		Rand:        rand.New(source),
		MaxIter:     maxIter,
		Limit:       limit,
	}
}

// CheckConstraints checks if the data (up to limit) satisfies the constraints.
// Returns true if valid, and the index of the first row that violates a constraint.
func (s *Shuffler) CheckConstraints(data [][]string) (bool, int) {
	if len(data) == 0 {
		return true, 0
	}

	n := len(data)
	if s.Limit > 0 && s.Limit < n {
		n = s.Limit
	}

	// Track repetitions for each column
	reps := make([]int, len(s.Constraints))
	for i := range reps {
		reps[i] = 1
	}

	for i := 1; i < n; i++ {
		for col, c := range s.Constraints {
			if col >= len(data[i]) || col >= len(data[i-1]) {
				continue
			}

			// Max repetition constraint (positive)
			if c > 0 {
				if data[i][col] == data[i-1][col] {
					reps[col]++
					if reps[col] > int(c) {
						return false, i
					}
				} else {
					reps[col] = 1
				}
			}

			// Min gap constraint (negative)
			if c < 0 {
				gap := int(-c)
				start := i - gap
				if start < 0 {
					start = 0
				}
				for j := start; j < i; j++ {
					if data[j][col] == data[i][col] {
						return false, i
					}
				}
			}
		}
	}

	return true, 0
}

// ShuffleEquiprob performs an equiprobable shuffle by generating random
// permutations and filtering them.
func (s *Shuffler) ShuffleEquiprob() ([][]string, error) {
	iter := 0
	maxIter := s.MaxIter
	if maxIter <= 0 {
		maxIter = 1000 // Default for equiprob
	}

	for iter < maxIter {
		shuffled := make([][]string, len(s.Data))
		copy(shuffled, s.Data)
		s.Rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		ok, _ := s.CheckConstraints(shuffled)
		if ok {
			return shuffled[:s.Limit], nil
		}
		iter++
	}

	return nil, fmt.Errorf("could not find a valid permutation after %d iterations", iter)
}

// fitsAtPosition checks if the row at table[candidateIdx] can be placed at table[pos]
// without violating any constraints, given the current repetition counts.
func (s *Shuffler) fitsAtPosition(table [][]string, pos int, candidateIdx int, currentReps []int) bool {
	candidateRow := table[candidateIdx]

	for col, constraint := range s.Constraints {
		if col >= len(candidateRow) {
			continue
		}

		// Check Max Repetition Constraint (positive values)
		if constraint > 0 && pos > 0 {
			prevRow := table[pos-1]
			if col < len(prevRow) && candidateRow[col] == prevRow[col] {
				if currentReps[col]+1 > int(constraint) {
					return false
				}
			}
		}

		// Check Minimum Gap Constraint (negative values)
		if constraint < 0 {
			gap := int(-constraint)
			start := pos - gap
			if start < 0 {
				start = 0
			}
			for i := start; i < pos; i++ {
				if col < len(table[i]) && table[i][col] == candidateRow[col] {
					return false
				}
			}
		}
	}
	return true
}

// updateRepetitionCount calculates the repetition count for the row at table[pos]
// for each constrained column, updating the reps slice.
func (s *Shuffler) updateRepetitionCount(table [][]string, pos int, reps []int) {
	if pos == 0 {
		for i := range reps {
			reps[i] = 1
		}
		return
	}

	for col := range s.Constraints {
		if col < len(table[pos]) && col < len(table[pos-1]) && table[pos][col] == table[pos-1][col] {
			reps[col]++
		} else {
			reps[col] = 1
		}
	}
}

// ShuffleConstructive builds a valid permutation by swapping lines as it goes.
// It uses a greedy approach, trying to find a valid row for each position sequentially.
func (s *Shuffler) ShuffleConstructive() ([][]string, error) {
	numRows := len(s.Data)
	maxAttempts := s.MaxIter
	if maxAttempts <= 0 {
		maxAttempts = numRows
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Create a local copy of the data and perform an initial random shuffle.
		table := make([][]string, numRows)
		copy(table, s.Data)
		s.Rand.Shuffle(numRows, func(i, j int) {
			table[i], table[j] = table[j], table[i]
		})

		// If there are no constraints, the initial shuffle is enough.
		if len(s.Constraints) == 0 {
			return table[:s.Limit], nil
		}

		reps := make([]int, len(s.Constraints))
		s.updateRepetitionCount(table, 0, reps)

		success := true
		// Sequentially build the table up to the limit.
		for pos := 1; pos < s.Limit; pos++ {
			foundValidCandidate := false
			// Search for a candidate row (from the remaining rows) that fits at the current position.
			for candidateIdx := pos; candidateIdx < numRows; candidateIdx++ {
				if s.fitsAtPosition(table, pos, candidateIdx, reps) {
					// Swap the valid candidate into the current position.
					table[pos], table[candidateIdx] = table[candidateIdx], table[pos]
					s.updateRepetitionCount(table, pos, reps)
					foundValidCandidate = true
					break
				}
			}

			if !foundValidCandidate {
				success = false
				break
			}
		}

		if success {
			return table[:s.Limit], nil
		}
	}

	return nil, fmt.Errorf("could not find a valid permutation after %d attempts", maxAttempts)
}

// ParseConstraints parses a string like "1 2 -3" or "1,2,-3" into a slice of Constraint.
func ParseConstraints(s string) ([]Constraint, error) {
	// Handle both spaces and commas
	normalized := strings.ReplaceAll(s, ",", " ")
	fields := strings.Fields(normalized)
	constraints := make([]Constraint, len(fields))
	for i, f := range fields {
		var val int
		_, err := fmt.Sscanf(f, "%d", &val)
		if err != nil {
			return nil, fmt.Errorf("invalid constraint: %s", f)
		}
		constraints[i] = Constraint(val)
	}
	return constraints, nil
}
