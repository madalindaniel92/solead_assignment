package scrappy

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Website struct {
	Domain string
}

func LoadFromFile(path string) ([]Website, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseCSV(file)
}

func ParseCSV(reader io.Reader) ([]Website, error) {
	results := []Website{}
	scanner := bufio.NewScanner(reader)

	// Parse each line of the CSV, trimming whitespace and validating URLs
	for index := 0; scanner.Scan(); index++ {
		line := strings.TrimSpace(scanner.Text())

		// Check CSV header is "domain", return error otherwise
		if index == 0 {
			err := checkCSVHeader(line, "domain")
			if err != nil {
				return nil, err
			}
			continue
		}

		results = append(results, Website{Domain: line})
	}

	return results, scanner.Err()
}

func checkCSVHeader(line, expected string) error {
	if line != expected {
		return fmt.Errorf("%w: expected '%s'", ErrInvalidCSVHeader, expected)
	}

	return nil
}
