package web

import (
	"log"
	"net/url"

	"github.com/gocolly/colly/v2"
)

// GetLinks returns the child links of the parent `selector` from `rawURL`,
func GetLinks(rawUrl string, selector string) (links []string, err error) {
	domain, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	// Create a collector specifically for this domain
	c := NewCollector(domain)

	// Get all of the link hrefs from each nav element
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		navLinks := e.ChildAttrs("a", "href")
		links = append(links, navLinks...)
	})

	// Log each visited endpoint
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %q\n", r.URL.String())
	})

	// Do the thing! (visit domain and start scraping)
	err = c.Visit(domain.String())

	return links, err
}
