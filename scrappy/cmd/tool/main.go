package main

import (
	"flag"
	"fmt"
	"os"

	"examples/scrappy"
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

	results, err := scrappy.LoadFromFile(*csvFile)
	checkErr(err)

	for _, result := range results {
		fmt.Printf("Domain: %q\n", result.Domain)
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
		os.Exit(1)
	}
}
