// Package CSV handling parsing and validation of input CSV files.
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

		parsedURL, err := ParseURL(line)
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

	return results, nil
}

func checkCSVHeader(line, expected string) error {
	if line != expected {
		return fmt.Errorf("%w: expected '%s'", ErrInvalidCSVHeader, expected)
	}

	return nil
}

func ParseURL(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, ErrMissingURLHost
	}

	// If we don't have the URL scheme, we assume it is "https://"
	if !strings.Contains(rawURL, "://") {
		rawURL = fmt.Sprintf("https://%s", rawURL)
	}

	result, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, err)
	}

	// Check that the domain host is present
	if result.Host == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingURLHost, rawURL)
	}

	// Only allow http and https URLs
	if result.Scheme != "http" && result.Scheme != "https" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURLScheme, result.Scheme)
	}

	return result, nil
}
