package web

import (
	"net/http"
)

// GetRobots returns the "robots.txt" file of a domain
func GetRobots(url string) (*http.Response, error) {
	// Get <domain>/robots.txt
	return NewClient(defaultTimeout).Get(url + "/robots.txt")
}
