package web

import (
	"errors"
	"strings"
	"testing"
)

func TestSitemapFromRobots(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name: "robots.txt with sitemap",
			body: `
				User-agent: *
				Disallow: /wp-admin/
				Allow: /wp-admin/admin-ajax.php

				Sitemap: https://cumberland-river.com/wp-sitemap.xml
			`,
			expected: "https://cumberland-river.com/wp-sitemap.xml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)
			url, err := sitemapFromRobots(reader)
			checkNoErr(t, err)

			if url != tc.expected {
				t.Errorf("Expected %q, got %q instead", tc.expected, url)
			}
		})
	}
}

func TestSitemapFromRobots_failure(t *testing.T) {
	testCases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name: "robots.txt with no sitemap",
			body: `
				User-agent: *
				Disallow: /wp-admin/
				Allow: /wp-admin/admin-ajax.php
			`,
			expectedErr: ErrNotFound,
		},
		{
			name: "robots.txt with invalid sitemap URL",
			body: `
				User-agent: *
				Disallow: /wp-admin/
				Allow: /wp-admin/admin-ajax.php

				Sitemap: nope nope nope invalid
			`,
			expectedErr: ErrInvalidURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.body)
			_, err := sitemapFromRobots(reader)

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
