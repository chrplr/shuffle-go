package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrplr/shuffle/pkg/shuffle"
)

func main() {
	constrFlag := flag.String("c", "", "Constraints: space-separated numbers (positive for max rep, negative for min gap)")
	seedFlag := flag.Int64("s", 0, "Random seed (0 for current time)")
	limitFlag := flag.Int("n", 0, "Limit output to n lines (0 for all)")
	iterFlag := flag.Int("i", 0, "Max iterations/loops (0 for default)")
	equiprobFlag := flag.Bool("e", false, "Use equiprobable shuffle algorithm (slower)")
	delimFlag := flag.String("d", "", "Field delimiter (empty for whitespace)")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

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
