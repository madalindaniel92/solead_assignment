// Package CSV handling parsing and validation of input CSV files.
package csv

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

type Website struct {
	Domain url.URL
}

func (w *Website) URL() string {
	return w.Domain.String()
}

type Company struct {
	Domain            JSONUrl  `json:"domain"`
	CommercialName    string   `json:"commercial_name"`
	LegalName         string   `json:"legal_name"`
	AllAvailableNames []string `json:"all_available_names"`
}

type JSONUrl struct {
	*url.URL
}

func (j JSONUrl) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.URL.String())
}

func (j *JSONUrl) UnmarshalJSON(raw []byte) error {
	// URL field is surrounded by quotes, check we have at least 2 chars
	if len(raw) < 2 {
		return fmt.Errorf("invalid url field")
	}

	// Strip off surrounding quotes
	raw = raw[1 : len(raw)-1]

	// Check we have an URL that can be parsed
	parsedURL, err := url.Parse(string(raw))
	if err != nil {
		return err
	}

	j.URL = parsedURL
	return nil
}

func LoadDomainsFromFile(path string) ([]Website, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results, err := ParseDomainsCSV(file)
	return results, wrapWithPathInfo(err, path)
}

func LoadCompaniesFromFile(path string) ([]Company, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results, err := ParseCompaniesCSV(file)
	return results, wrapWithPathInfo(err, path)
}

// Wrap error with file path information
func wrapWithPathInfo(err error, path string) error {
	if err != nil {
		return fmt.Errorf("%s - %w", path, err)
	}
	return nil
}

func ParseDomainsCSV(reader io.Reader) ([]Website, error) {
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

var companyCSVHeader = []string{
	"domain",
	"company_commercial_name",
	"company_legal_name",
	"company_all_available_names",
}

func ParseCompaniesCSV(reader io.Reader) (companies []Company, err error) {
	csvReader := csv.NewReader(reader)
	// Reuse the same slice for each line, to prevent too many allocations
	csvReader.ReuseRecord = true

	// Some CSV lines may be invalid, accumulate them so we can show them in an error message
	var invalidLines ErrInvalidCSVLines

	// Header lines might appear in any order, so we need to determine the correct field indexes
	domainIndex, commercialIndex, legalIndex, allRawIndex := 0, 0, 0, 0

	// Parse each line of the CSV
	for index := 0; true; index++ {
		var line []string
		line, err = csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			if errors.Is(err, csv.ErrFieldCount) {
				err = wrapWrongNumFieldsErr(err)
				rawLine := strings.Join(line, ",")
				invalidLines = invalidLines.Append(err, rawLine, index)
				continue
			}

			return companies, fmt.Errorf("%w - %s", ErrParseCSV, err)
		}

		if index == 0 {
			// Check that we have the expected headers,
			// and determine the order in which the headers appear
			indexes, err := checkCSVHeaders(line, companyCSVHeader)
			if err != nil {
				return nil, err
			}
			domainIndex, commercialIndex, legalIndex, allRawIndex = indexes[0], indexes[1], indexes[2], indexes[3]
			continue
		}

		// Line length is checked by csvReader,
		// indexes are determined from the csv header line
		domain, commercial, legal, allRaw := line[domainIndex], line[commercialIndex], line[legalIndex], line[allRawIndex]
		parsedURL, err := ParseURL(domain)
		if err != nil {
			invalidLines = invalidLines.Append(err, strings.TrimSpace(domain), index)
			continue
		}

		companies = append(companies, Company{
			Domain:            JSONUrl{URL: parsedURL},
			CommercialName:    strings.TrimSpace(commercial),
			LegalName:         strings.TrimSpace(legal),
			AllAvailableNames: splitAndTrimFields(allRaw, "|"),
		})
	}

	// Check if we have invalid lines
	if len(invalidLines) > 0 {
		return companies, invalidLines
	}

	// Check we have some results (we consider empty CSVs an error case)
	if len(companies) == 0 {
		return nil, ErrEmptyCSV
	}

	return companies, nil
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

// Helpers

func wrapWrongNumFieldsErr(err error) error {
	return fmt.Errorf("%w - expected %d fields", ErrWrongNumberOfFields, len(companyCSVHeader))
}

func splitAndTrimFields(text string, separator string) []string {
	fields := strings.Split(text, separator)

	for index, field := range fields {
		fields[index] = strings.TrimSpace(field)
	}

	return fields
}

func checkCSVHeader(line, expected string) error {
	if line != expected {
		return fmt.Errorf("%w: expected '%s'", ErrInvalidCSVHeader, expected)
	}

	return nil
}

func expectedHeadersErr(actual, expected []string) error {
	return fmt.Errorf("%w:\n    expected headers %v,\n    got %v instead",
		ErrInvalidCSVHeader, expected, actual)
}

// checkCSVHeaders returns an error if the actual headers don't match the expected ones.
// In case the headers match, we return the index of each header from expected in actual
func checkCSVHeaders(actual, expected []string) ([]int, error) {
	if len(actual) != len(expected) {
		return nil, expectedHeadersErr(actual, expected)
	}

	indexes := make([]int, 0, len(expected))

	for _, header := range expected {
		index := indexOf(actual, header)
		if index < 0 {
			return nil, expectedHeadersErr(actual, expected)
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func indexOf(values []string, needle string) int {
	for index, value := range values {
		if value == needle {
			return index
		}
	}

	return -1
}
