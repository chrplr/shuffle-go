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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrplr/shuffle-go/internal/version"
	"github.com/chrplr/shuffle-go"
)

func main() {
	constrFlag := flag.String("c", "", "Constraints: space-separated numbers (positive for max rep, negative for min gap)")
	seedFlag := flag.Int64("s", 0, "Random seed (0 for current time)")
	limitFlag := flag.Int("n", 0, "Limit output to n lines (0 for all)")
	iterFlag := flag.Int("i", 0, "Max iterations/loops (0 for default)")
	equiprobFlag := flag.Bool("e", false, "Use equiprobable shuffle algorithm (slower)")
	delimFlag := flag.String("d", "", "Field delimiter (empty for whitespace)")
	versionFlag := flag.Bool("v", false, "Show version")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("shuffle-cli version %s\n", version.Version)
		fmt.Printf("%s\n", version.Info)
		os.Exit(0)
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	var input *os.File
	var err error
	if flag.NArg() > 0 {
		input, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	data, err := shuffle.LoadData(input, *delimFlag)
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	constraints, err := shuffle.ParseConstraints(*constrFlag)
	if err != nil {
		log.Fatalf("Error parsing constraints: %v", err)
	}

	s := shuffle.NewShuffler(data, constraints, *seedFlag, *iterFlag, *limitFlag)

	var result [][]string
	if *equiprobFlag {
		result, err = s.ShuffleEquiprob()
	} else {
		result, err = s.ShuffleConstructive()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Shuffle failed: %v\n", err)
		os.Exit(1)
	}

	delimiter := *delimFlag
	if delimiter == "" {
		delimiter = " "
	}

	for _, row := range result {
		fmt.Println(strings.Join(row, delimiter))
	}
}
