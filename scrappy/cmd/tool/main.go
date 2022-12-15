package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"examples/scrappy/internal/csv"
)

func main() {
	// Parse flags
	csvFile := flag.String("f", "", "csv file to load domain names from")

	flag.Parse()

	// Validate flags
	if *csvFile == "" {
		usage()
		checkErr(fmt.Errorf("missing csv file"))
	}

	results, err := csv.LoadFromFile(*csvFile)
	checkErr(err)

	for _, result := range results {
		domain := result.Domain
		fmt.Printf("Domain: %s://%s\n", domain.Scheme, domain.Hostname())
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `
	Usage:
		./tool -f <csv file>
	`)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %q\n", err)
		printExtraErrInfo(err)
		os.Exit(1)
	}
}

func printExtraErrInfo(err error) {
	switch value := errors.Unwrap(err).(type) {
	case csv.ErrInvalidCSVLines:
		for _, invalidLine := range value {
			fmt.Fprintf(os.Stderr, "Invalid line %d %q: %s\n",
				invalidLine.Index, invalidLine.Line, invalidLine.Err)
		}
	}
}
