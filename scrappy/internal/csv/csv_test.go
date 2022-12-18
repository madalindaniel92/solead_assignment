package csv_test

import (
	"errors"
	"net/url"
	"strings"
	"testing"

	"examples/scrappy/internal/csv"
)

func TestParseDomainsCSV(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		expected []csv.Website
	}{
		{
			name: "valid domains",
			body: `domain
				https://en.wikipedia.org
				https://google.com
				http://example.com
			`,
			expected: []csv.Website{
				{Domain: url.URL{Host: "en.wikipedia.org", Scheme: "https"}},
				{Domain: url.URL{Host: "google.com", Scheme: "https"}},
				{Domain: url.URL{Host: "example.com", Scheme: "http"}},
			},
		},
		{
			name: "domains without URI schemes",
			body: `domain
				bostonzen.org
				mazautoglass.com
				melatee.com
				timent.com`,
			expected: []csv.Website{
				{Domain: url.URL{Host: "bostonzen.org", Scheme: "https"}},
				{Domain: url.URL{Host: "mazautoglass.com", Scheme: "https"}},
				{Domain: url.URL{Host: "melatee.com", Scheme: "https"}},
				{Domain: url.URL{Host: "timent.com", Scheme: "https"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			results, err := csv.ParseDomainsCSV(reader)
			checkNoErr(t, err)

			if len(results) != len(tc.expected) {
				t.Fatalf("Expected %d results, received %d instead",
					len(tc.expected), len(results))
			}

			// zip across results and expected, checking Host and Scheme
			for index, result := range results {
				expected := tc.expected[index]

				if result.Domain.Host != expected.Domain.Host {
					t.Errorf("Expected host %q, got host %q instead (index %d)",
						expected.Domain.Host, result.Domain.Host, index)
				}

				expectedScheme := expected.Domain.Scheme
				if expectedScheme != "" && result.Domain.Scheme != expectedScheme {
					t.Errorf("Expected scheme %q, got scheme %q instead (index %d)",
						expectedScheme, result.Domain.Scheme, index)
				}
			}
		})
	}
}

func TestParseCSV_failure(t *testing.T) {
	testCases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name: "empty file",
			body: "",
			// we expect the file to have the "domain" header
			expectedErr: csv.ErrEmptyCSV,
		},
		{
			name: "invalid header",
			body: `first_name, last_name, address
				Daniel, Smith, Someplace Nice 42
			`,
			expectedErr: csv.ErrInvalidCSVHeader,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			_, err := csv.ParseDomainsCSV(reader)
			checkErrIs(t, err, tc.expectedErr)
		})
	}
}

func TestParseDomainsCSV_invalidLines(t *testing.T) {
	testCases := []struct {
		name        string
		body        string
		expectedErr error
		// which lines are invalid
		expectedIndexes []int
		// even though we have invalid lines,
		// we still return the results that are valid
		expectedResults []csv.Website
	}{
		{
			name: "invalid domains",
			body: `domain
				bostonzen.org
				invalid right here
				mazautoglass.com
				dragons-are-awesome.com
				not quite valid either
				melatee.com`,
			expectedErr:     csv.ErrInvalidCSVLines{},
			expectedIndexes: []int{2, 5}, // lines 2 and 5 are invalid
			// even though we have invalid lines,
			expectedResults: []csv.Website{
				{Domain: url.URL{Host: "bostonzen.org", Scheme: "https"}},
				{Domain: url.URL{Host: "mazautoglass.com", Scheme: "https"}},
				{Domain: url.URL{Host: "mazautoglass.com", Scheme: "https"}},
				{Domain: url.URL{Host: "melatee.com", Scheme: "https"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			_, err := csv.ParseDomainsCSV(reader)
			checkErrIs(t, err, tc.expectedErr)

			errLines, _ := err.(csv.ErrInvalidCSVLines)

			if len(errLines) != len(tc.expectedIndexes) {
				t.Fatalf("Expected %d invalid lines, got %d instead",
					len(tc.expectedIndexes), len(errLines))
			}

			// zip across errLines and expectedIndexes and compare index values
			for index, errLine := range errLines {
				expectedIndex := tc.expectedIndexes[index]

				if errLine.Index != expectedIndex {
					t.Errorf("Expected error line %d, got %d instead (index %d)",
						expectedIndex, errLine.Index, index)
				}
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected url.URL
	}{
		{
			name:     "valid URL",
			url:      "https://en.wikipedia.org",
			expected: url.URL{Scheme: "https", Host: "en.wikipedia.org"},
		},
		{
			name:     "http URL",
			url:      "http://why_no_tsl.org",
			expected: url.URL{Scheme: "http", Host: "why_no_tsl.org"},
		},
		{
			name:     "domain without scheme",
			url:      "example.com",
			expected: url.URL{Scheme: "https", Host: "example.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := csv.ParseURL(tc.url)
			checkNoErr(t, err)

			if *result != tc.expected {
				t.Errorf("Expected %+v, got %+v instead", tc.expected, *result)
			}
		})
	}
}

func TestParseURL_failure(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		expectedErr error
	}{
		{
			name:        "invalid URL",
			url:         "https://en wikipedia dot org",
			expectedErr: csv.ErrInvalidURL,
		},
		{
			name:        "invalid URL scheme",
			url:         "redis://some-host.com",
			expectedErr: csv.ErrInvalidURLScheme,
		},
		{
			name:        "missing URL",
			url:         "",
			expectedErr: csv.ErrMissingURLHost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := csv.ParseURL(tc.url)
			checkErrIs(t, err, tc.expectedErr)
		})
	}
}

// Helpers

func checkNoErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error: %q", err)
	}
}

func checkErrIs(t *testing.T, err error, expected error) {
	t.Helper()

	if !errors.Is(err, expected) {
		t.Fatalf("Expected error %q, got %q (%T) instead", expected, err, err)
	}
}
