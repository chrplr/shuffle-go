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

// ShuffleConstructive builds a valid permutation by swapping lines as it goes.
func (s *Shuffler) ShuffleConstructive() ([][]string, error) {
	nlines := len(s.Data)
	maxLoops := s.MaxIter
	if maxLoops <= 0 {
		maxLoops = nlines
	}

	loop := 0
	for loop < maxLoops {
		// Initial random shuffle
		table := make([][]string, nlines)
		copy(table, s.Data)
		s.Rand.Shuffle(nlines, func(i, j int) {
			table[i], table[j] = table[j], table[i]
		})

		if len(s.Constraints) == 0 {
			return table[:s.Limit], nil
		}

		badPermut := false
		reps := make([]int, len(s.Constraints))
		for k := range reps {
			reps[k] = 1
		}

		for i := 1; i < s.Limit; i++ {
			passLine := false
			for j := i; j < nlines; j++ {
				// Check if table[j] fits at position i
				fail := false
				for col, c := range s.Constraints {
					if col >= len(table[j]) {
						continue
					}

					// Max rep
					if c > 0 {
						if table[j][col] == table[i-1][col] && reps[col]+1 > int(c) {
							fail = true
							break
						}
					}

					// Min gap
					if c < 0 {
						gap := int(-c)
						start := i - gap
						if start < 0 {
							start = 0
						}
						for k := start; k < i; k++ {
							if table[k][col] == table[j][col] {
								fail = true
								break
							}
						}
						if fail {
							break
						}
					}
				}

				if !fail {
					// Swap table[i] and table[j]
					table[i], table[j] = table[j], table[i]
					passLine = true
					
					// Update repetitions for the next row
					for col := range s.Constraints {
						if col < len(table[i]) && col < len(table[i-1]) {
							if table[i][col] == table[i-1][col] {
								reps[col]++
							} else {
								reps[col] = 1
							}
						}
					}
					break
				}
			}

			if !passLine {
				badPermut = true
				break
			}
		}

		if !badPermut {
			return table[:s.Limit], nil
		}
		loop++
	}

	return nil, fmt.Errorf("could not find a valid permutation after %d loops", loop)
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
