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
	"reflect"
	"strings"
	"testing"
)

func TestCheckConstraints(t *testing.T) {
	data := [][]string{
		{"A", "1"},
		{"A", "2"},
		{"B", "1"},
		{"B", "2"},
	}

	tests := []struct {
		name        string
		constraints []Constraint
		wantOk      bool
	}{
		{"No constraints", []Constraint{0, 0}, true},
		{"Max rep 2 OK col 0", []Constraint{2, 0}, true},
		{"Max rep 1 Fail col 0", []Constraint{1, 0}, false},
		{"Min gap 2 Fail col 0", []Constraint{-2, 0}, false}, // A, A is gap 1
		{"Min gap 1 Fail col 0", []Constraint{-1, 0}, false}, // gap 1 means no same labels adjacent
		{"Col 1 Max rep 1 OK", []Constraint{0, 1}, true},    // 1, 2, 1, 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShuffler(data, tt.constraints, 1, 0, 0)
			gotOk, _ := s.CheckConstraints(data)
			if gotOk != tt.wantOk {
				t.Errorf("CheckConstraints() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestShuffleConstructive(t *testing.T) {
	data := [][]string{
		{"A"}, {"A"}, {"B"}, {"B"}, {"C"}, {"C"},
	}
	// Max 1 repetition means no two identical labels can be adjacent
	constraints := []Constraint{1}
	s := NewShuffler(data, constraints, 42, 100, 0)

	got, err := s.ShuffleConstructive()
	if err != nil {
		t.Fatalf("ShuffleConstructive() error = %v", err)
	}

	if len(got) != len(data) {
		t.Errorf("ShuffleConstructive() length = %d, want %d", len(got), len(data))
	}

	ok, row := s.CheckConstraints(got)
	if !ok {
		t.Errorf("ShuffleConstructive() produced invalid sequence at row %d: %v", row, got)
	}
}

func TestShuffleEquiprob(t *testing.T) {
	data := [][]string{
		{"A"}, {"B"}, {"C"},
	}
	// No constraints, should always succeed
	s := NewShuffler(data, []Constraint{0}, 42, 10, 0)

	got, err := s.ShuffleEquiprob()
	if err != nil {
		t.Fatalf("ShuffleEquiprob() error = %v", err)
	}

	if len(got) != len(data) {
		t.Errorf("ShuffleEquiprob() length = %d, want %d", len(got), len(data))
	}
}

func TestLoadData(t *testing.T) {
	input := `row1 col1
row2 col2

row3 col3`
	
	t.Run("Whitespace delimiter", func(t *testing.T) {
		got, err := LoadData(strings.NewReader(input), "")
		if err != nil {
			t.Fatalf("LoadData() error = %v", err)
		}
		want := [][]string{
			{"row1", "col1"},
			{"row2", "col2"},
			{"row3", "col3"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("LoadData() got = %v, want %v", got, want)
		}
	})

	t.Run("Comma delimiter", func(t *testing.T) {
		csvInput := "a,b,c\nd,e,f"
		got, err := LoadData(strings.NewReader(csvInput), ",")
		if err != nil {
			t.Fatalf("LoadData() error = %v", err)
		}
		want := [][]string{
			{"a", "b", "c"},
			{"d", "e", "f"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("LoadData() got = %v, want %v", got, want)
		}
	})
}

func TestParseConstraints(t *testing.T) {
	tests := []struct {
		input   string
		want    []Constraint
		wantErr bool
	}{
		{"1 2 -3", []Constraint{1, 2, -3}, false},
		{"", []Constraint{}, false},
		{"  1   -2  ", []Constraint{1, -2}, false},
		{"abc", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseConstraints(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConstraints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConstraints() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffleImpossible(t *testing.T) {
	data := [][]string{
		{"A"}, {"A"}, {"A"},
	}
	// Impossible: max 1 repetition of "A" when there are 3 "A"s
	constraints := []Constraint{1}
	s := NewShuffler(data, constraints, 1, 10, 0)

	_, err := s.ShuffleConstructive()
	if err == nil {
		t.Error("ShuffleConstructive() should have failed for impossible constraints")
	}

	_, err = s.ShuffleEquiprob()
	if err == nil {
		t.Error("ShuffleEquiprob() should have failed for impossible constraints")
	}
}
