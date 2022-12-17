package web

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

// GetRobots returns the "robots.txt" file of a domain
func GetRobots(rawUrl string) (*http.Response, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	// Get <domain>/robots.txt
	parsedUrl.Path = "/robots.txt"

	return NewClient(defaultTimeout).Get(parsedUrl.String())
}

// GetSitemapLinks requests the robots.txt file of a website,
// then traverses the sitemap to get the required links
func GetSitemapLinks(url string) (links []string, err error) {
	response, err := GetRobots(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	sitemapUrl, err := sitemapFromRobots(response.Body)
	if err != nil {
		return nil, err
	}

	return CollectSitemapLinks(sitemapUrl)
}

// CollectSitemapLinks parses a sitemap extracting links.
//
// Example code from:
//
//	https://github.com/gocolly/colly/blob/master/_examples/shopify_sitemap/shopify_sitemap.go
func CollectSitemapLinks(sitemapUrl string) (links []string, err error) {
	domain, err := url.Parse(sitemapUrl)
	if err != nil {
		return nil, err
	}

	// Create a collector specifically for this domain
	c := colly.NewCollector(colly.AllowedDomains(domain.Host))

	// Extract locations from sitemap
	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		links = append(links, e.Text)
	})

	// Visit each sitemap in index
	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		c.Visit(e.Text)
	})

	// Log each visited endpoint
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %q\n", r.URL.String())
	})

	// Do the thing! (visit domain and start scraping)
	err = c.Visit(domain.String())

	return links, err
}

const sitemapPrefix = "Sitemap:"

func sitemapFromRobots(in io.Reader) (sitemapUrl string, err error) {
	scanner := bufio.NewScanner(in)

	// Traverse each line of input looking for the sitemap url
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// If line starts with the sitemap prefix,
		// parse and return the urn
		if strings.HasPrefix(line, sitemapPrefix) {
			// Trim the prefix and parse the URL
			rawUrl := strings.TrimPrefix(line, sitemapPrefix)
			parsedUrl, err := url.Parse(strings.TrimSpace(rawUrl))
			if err != nil {
				return "", sitemapErr(rawUrl, err)
			}

			// We need the host to be present (absolute URL)
			if parsedUrl.Host == "" {
				err = fmt.Errorf("missing host")
				return "", sitemapErr(rawUrl, err)
			}

			// Return sitemap url if it has been successfully parsed
			return parsedUrl.String(), nil
		}
	}

	// Check if we have any errors reading lines
	if scanner.Err() != nil {
		return "", fmt.Errorf("failed to parse sitemap %w", scanner.Err())
	}

	// No line had sitemap prefix
	return "", fmt.Errorf("sitemap %w", ErrNotFound)
}

func sitemapErr(rawUrl string, err error) error {
	return fmt.Errorf("%w: failed to parse sitemap %q - %q",
		ErrInvalidURL, rawUrl, err)
}
