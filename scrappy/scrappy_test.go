package scrappy_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"examples/scrappy"
)

func TestParseCSV(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		expected []scrappy.Website
	}{
		{
			name: "valid domains",
			body: `domain
				bostonzen.org
				mazautoglass.com
				melatee.com
				timent.com`,
			expected: []scrappy.Website{
				{Domain: "bostonzen.org"},
				{Domain: "mazautoglass.com"},
				{Domain: "melatee.com"},
				{Domain: "timent.com"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			result, err := scrappy.ParseCSV(reader)
			checkNoErr(t, err)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v instead", tc.expected, result)
			}
		})
	}
}

func TestParseCSV_failures(t *testing.T) {
	testCases := []struct {
		name            string
		body            string
		expectedErr     error
		expectedResults []scrappy.Website
	}{
		{
			name: "invalid header",
			body: `first_name, last_name, address
				Daniel, Smith, Someplace Nice 42
			`,
			expectedErr: scrappy.ErrInvalidCSVHeader,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)

			_, err := scrappy.ParseCSV(reader)
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
		t.Fatalf("Expected error %q, got %q instead", expected, err)
	}
}
