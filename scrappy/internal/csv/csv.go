package csv

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

type Website struct {
	Domain url.URL
}

func LoadFromFile(path string) ([]Website, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results, err := ParseCSV(file)

	// Wrap error with file path information
	if err != nil {
		err = fmt.Errorf("%s - %w", path, err)
	}

	return results, err
}

func ParseCSV(reader io.Reader) ([]Website, error) {
	var results []Website

	// Some CSV lines may be invalid, accumulate them so we can show them in an error message
	var invalidLines ErrInvalidCSVLines

	// Split input into lines using a scanner
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

		// Ignore empty lines
		if line == "" {
			continue
		}

		parsedURL, err := parseURL(line)
		if err != nil {
			invalidLines = invalidLines.Append(err, line, index)
			continue
		}

		results = append(results, Website{Domain: *parsedURL})
	}

	// Check if we have invalid lines
	if len(invalidLines) > 0 {
		return results, invalidLines
	}

	// Check line reader error
	if scanner.Err() != nil {
		return results, scanner.Err()
	}

	// Check we have some results (we consider empty CSVs an error case)
	if len(results) == 0 {
		return nil, ErrEmptyCSV
	}

	return results, scanner.Err()
}

func checkCSVHeader(line, expected string) error {
	if line != expected {
		return fmt.Errorf("%w: expected '%s'", ErrInvalidCSVHeader, expected)
	}

	return nil
}

func parseURL(rawURL string) (*url.URL, error) {
	result, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Since we don't have the scheme, we assume it is "https://"
	if result.Scheme == "" {
		// We need to reparse to extract URL correctly
		result, err = url.Parse("https://" + rawURL)
	}

	if err != nil {
		return nil, err
	}

	if result.Host == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingURLHost, rawURL)
	}

	return result, nil
}