package csv_test

import (
	"errors"
	"net/url"
	"reflect"
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

				checkDomainUrl(t, &result.Domain, &expected.Domain, index)
			}
		})
	}
}

func TestParseDomainsCSV_failure(t *testing.T) {
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
			expectedErr: csv.ErrInvalidCSVLines([]csv.InvalidCSVLine{
				{Index: 2, Line: "invalid right here", Err: csv.ErrInvalidURL},
				{Index: 5, Line: "not quite valid either", Err: csv.ErrInvalidURL},
			}),
			// even though we have invalid lines,
			expectedResults: []csv.Website{
				{Domain: url.URL{Host: "bostonzen.org", Scheme: "https"}},
				{Domain: url.URL{Host: "mazautoglass.com", Scheme: "https"}},
				{Domain: url.URL{Host: "dragons-are-awesome.com", Scheme: "https"}},
				{Domain: url.URL{Host: "melatee.com", Scheme: "https"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			results, err := csv.ParseDomainsCSV(reader)
			checkErrIs(t, err, tc.expectedErr)

			// Check results
			if len(results) != len(tc.expectedResults) {
				t.Fatalf("Expected %d results, received %d instead",
					len(tc.expectedResults), len(results))
			}

			// zip across results and expected, checking Host and Scheme
			for index, result := range results {
				expected := tc.expectedResults[index]
				checkDomainUrl(t, &result.Domain, &expected.Domain, index)
			}

			// Check error lines
			errLines, _ := err.(csv.ErrInvalidCSVLines)
			expectedLines, _ := tc.expectedErr.(csv.ErrInvalidCSVLines)

			if len(errLines) != len(expectedLines) {
				t.Fatalf("Expected %d invalid lines, got %d instead",
					len(expectedLines), len(errLines))
			}

			// zip across errLines and expectedIndexes and compare index values
			for index, errLine := range errLines {
				expected := expectedLines[index]
				checkErrLine(t, &errLine, &expected, index)
			}
		})
	}
}

func TestParseCompaniesCSV(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		expected []csv.Company
	}{
		{
			name: "valid company info",
			body: `domain,company_commercial_name,company_legal_name,company_all_available_names
				bostonzen.org,Greater Boston Zen Center,GREATER BOSTON ZEN CENTER INC.,Greater Boston Zen Center | Boston Zen | GREATER BOSTON ZEN CENTER INC.
				mazautoglass.com,MAZ Auto Glass,,MAZ Auto Glass
				melatee.com,Melatee,,Melatee
				timent.com,Timent Technologies,,Timent Technologies | Timent`,
			expected: []csv.Company{
				{
					Domain:         url.URL{Host: "bostonzen.org", Scheme: "https"},
					CommercialName: "Greater Boston Zen Center",
					LegalName:      "GREATER BOSTON ZEN CENTER INC.",
					AllAvailableNames: []string{
						"Greater Boston Zen Center",
						"Boston Zen",
						"GREATER BOSTON ZEN CENTER INC.",
					},
				},
				{
					Domain:         url.URL{Host: "mazautoglass.com", Scheme: "https"},
					CommercialName: "MAZ Auto Glass",
					AllAvailableNames: []string{
						"MAZ Auto Glass",
					},
				},
				{
					Domain:         url.URL{Host: "melatee.com", Scheme: "https"},
					CommercialName: "Melatee",
					AllAvailableNames: []string{
						"Melatee",
					},
				},
				{
					Domain:         url.URL{Host: "timent.com", Scheme: "https"},
					CommercialName: "Timent Technologies",
					AllAvailableNames: []string{
						"Timent Technologies",
						"Timent",
					},
				},
			},
		},
		{
			name: "headers in a different order",
			body: `company_legal_name,domain,company_all_available_names,company_commercial_name
				GREATER BOSTON ZEN CENTER INC.,bostonzen.org,Greater Boston Zen Center | Boston Zen | GREATER BOSTON ZEN CENTER INC.,Greater Boston Zen Center
				,melatee.com,Melatee,Melatee`,
			expected: []csv.Company{
				{
					Domain:         url.URL{Host: "bostonzen.org", Scheme: "https"},
					CommercialName: "Greater Boston Zen Center",
					LegalName:      "GREATER BOSTON ZEN CENTER INC.",
					AllAvailableNames: []string{
						"Greater Boston Zen Center",
						"Boston Zen",
						"GREATER BOSTON ZEN CENTER INC.",
					},
				},
				{
					Domain:         url.URL{Host: "melatee.com", Scheme: "https"},
					CommercialName: "Melatee",
					AllAvailableNames: []string{
						"Melatee",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			results, err := csv.ParseCompaniesCSV(reader)
			checkNoErr(t, err)

			if len(results) != len(tc.expected) {
				t.Fatalf("Expected %d results, received %d instead",
					len(tc.expected), len(results))
			}

			// zip across results and expected, checking fields
			for index, result := range results {
				expected := tc.expected[index]

				checkDomainUrl(t, &result.Domain, &expected.Domain, index)
				checkCompanyNames(t, &result, &expected, index)
			}
		})
	}
}

func TestParseCompaniesCSV_failure(t *testing.T) {
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

			_, err := csv.ParseCompaniesCSV(reader)
			checkErrIs(t, err, tc.expectedErr)
		})
	}
}

func TestParseCompaniesCSV_invalidLines(t *testing.T) {
	testCases := []struct {
		name        string
		body        string
		expectedErr error
		// even though we have invalid lines,
		// we still return the results that are valid
		expectedResults []csv.Company
	}{
		{
			name: "invalid lines",
			body: `domain,company_commercial_name,company_legal_name,company_all_available_names
				bostonzen.org,Greater Boston Zen Center,GREATER BOSTON ZEN CENTER INC.,Greater Boston Zen Center | Boston Zen | GREATER BOSTON ZEN CENTER INC.
				acme.com,too, many, fields, on, this, line
				melatee.com,Melatee,,Melatee

				invalid url,Timent Technologies,,Timent Technologies | Timent
				xkcd.com, XKCD, XKCD Comics, xkcd | The awesome stick figure comic`,
			expectedErr: csv.ErrInvalidCSVLines([]csv.InvalidCSVLine{
				{Index: 2, Line: "\t\t\t\tacme.com,too, many, fields, on, this, line", Err: csv.ErrWrongNumberOfFields},
				{Index: 4, Line: "invalid url", Err: csv.ErrInvalidURL},
			}),
			expectedResults: []csv.Company{
				{
					Domain:         url.URL{Host: "bostonzen.org", Scheme: "https"},
					CommercialName: "Greater Boston Zen Center",
					LegalName:      "GREATER BOSTON ZEN CENTER INC.",
					AllAvailableNames: []string{
						"Greater Boston Zen Center",
						"Boston Zen",
						"GREATER BOSTON ZEN CENTER INC.",
					},
				},
				{
					Domain:         url.URL{Host: "melatee.com", Scheme: "https"},
					CommercialName: "Melatee",
					AllAvailableNames: []string{
						"Melatee",
					},
				},
				{
					Domain:         url.URL{Host: "xkcd.com", Scheme: "https"},
					CommercialName: "XKCD",
					LegalName:      "XKCD Comics",
					AllAvailableNames: []string{
						"xkcd",
						"The awesome stick figure comic",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			results, err := csv.ParseCompaniesCSV(reader)
			checkErrIs(t, err, tc.expectedErr)

			// Check results
			if len(results) != len(tc.expectedResults) {
				t.Fatalf("Expected %d results, received %d instead",
					len(tc.expectedResults), len(results))
			}

			// zip across results and expected, checking Host and Scheme
			for index, result := range results {
				expected := tc.expectedResults[index]
				checkDomainUrl(t, &result.Domain, &expected.Domain, index)
				checkCompanyNames(t, &result, &expected, index)
			}

			// Check error lines
			errLines, _ := err.(csv.ErrInvalidCSVLines)
			expectedLines, _ := tc.expectedErr.(csv.ErrInvalidCSVLines)

			if len(errLines) != len(expectedLines) {
				t.Fatalf("Expected %d invalid lines, got %d instead",
					len(expectedLines), len(errLines))
			}

			// zip across errLines and expectedIndexes and compare index values
			for index, errLine := range errLines {
				expected := expectedLines[index]
				checkErrLine(t, &errLine, &expected, index)
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

func checkDomainUrl(t *testing.T, result *url.URL, expected *url.URL, index int) {
	t.Helper()

	if result.Host != expected.Host {
		t.Errorf("Expected host %q, got host %q instead (index %d)",
			expected.Host, result.Host, index)
	}

	if expected.Scheme != "" && result.Scheme != expected.Scheme {
		t.Errorf("Expected scheme %q, got scheme %q instead (index %d)",
			expected.Scheme, result.Scheme, index)
	}
}

func checkCompanyNames(t *testing.T, result *csv.Company, expected *csv.Company, index int) {
	t.Helper()

	if result.CommercialName != expected.CommercialName {
		t.Errorf("Expected commercial name %q, got %q instead (index %d)",
			expected.CommercialName, result.CommercialName, index)
	}

	if result.LegalName != expected.LegalName {
		t.Errorf("Expected legal name %q, got %q instead (index %d)",
			expected.LegalName, result.LegalName, index)
	}

	if !reflect.DeepEqual(result.AllAvailableNames, expected.AllAvailableNames) {
		t.Errorf("Expected all available names %q, got %q instead (index %d)",
			expected.AllAvailableNames, result.AllAvailableNames, index)
	}
}

func checkErrLine(t *testing.T, result *csv.InvalidCSVLine, expected *csv.InvalidCSVLine, index int) {
	t.Helper()

	if result.Index != expected.Index {
		t.Errorf("Expected error line index %d, got %d instead (index %d)",
			expected.Index, result.Index, index)
	}

	if result.Line != expected.Line {
		t.Errorf("Expected error line %s, got %s instead (index %d)",
			expected.Line, result.Line, index)
	}

	if !errors.Is(result.Err, expected.Err) {
		t.Errorf("Expected error %q, got %q instead (index %d)",
			expected.Err, result.Err, index)
	}
}
