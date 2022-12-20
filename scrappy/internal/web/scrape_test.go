package web

import (
	"reflect"
	"testing"
)

func TestSpliceLink(t *testing.T) {
	testCases := []struct {
		name              string
		links             []string
		search            string
		expected          string
		expectedRemaining []string
	}{
		{
			name:   "get contact link",
			search: "contact",
			links: []string{
				"/wildlife",
				"/photos",
				"/shop",
				"/contact-us",
				"/philosophy",
			},
			expected: "/contact-us",
			expectedRemaining: []string{
				"/wildlife",
				"/photos",
				"/shop",
				"/philosophy",
			},
		},
		{
			name:   "get about link",
			search: "about",
			links: []string{
				"/wildlife",
				"/contact-us",
				"/shop",
				"/philosophy",
				"/about",
			},
			expected: "/about",
			expectedRemaining: []string{
				"/wildlife",
				"/contact-us",
				"/shop",
				"/philosophy",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			links := tc.links
			remaining, result := spliceLink(links, tc.search)

			if result != tc.expected {
				t.Errorf("Expected %q, got %q instead", tc.expected, result)
			}

			if !reflect.DeepEqual(remaining, tc.expectedRemaining) {
				t.Errorf("Expected %#v, got %#v instead", tc.expectedRemaining, remaining)
			}
		})
	}
}
