package web

import (
	"examples/scrappy/internal/phone"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// Maximum number of pages to scrape for each domain
const maxPagesScrapedPerDomain = 10

// ScrapeInfo represents the information gathered for a specific domain.
type ScrapeInfo struct {
	PhoneNumbers []phone.Phone
	LinksVisited []string
}

// EnoughInfo returns true once we have collected enough information for a domain.
func (s *ScrapeInfo) EnoughInfo() bool {
	return len(s.PhoneNumbers) > 0
}

// ExceededPageLimit returns true if we have exceeded the maximum number
// of pages to be scraped.
func (s *ScrapeInfo) ExceededPageLimit() bool {
	return len(s.LinksVisited) > maxPagesScrapedPerDomain
}

// SanitizePhoneNumbers will validate and deduplicate phone numbers.
//
// The phone number format is also normalized as part of the process.
func (s *ScrapeInfo) SanitizePhoneNumbers() {
	s.PhoneNumbers, _ = phone.ValidatePhoneNumbers(s.PhoneNumbers)
	s.PhoneNumbers = phone.DedupPhoneNumbers(s.PhoneNumbers)
}

// ScrapeJobResult represents the result of each worker running ScrapeDomain
type ScrapeJobResult struct {
	Url  string
	Info ScrapeInfo
	Err  error
}

func ScrapeDomain(domain string) (*ScrapeInfo, error) {
	// Check domain URL first
	domainUrl, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}

	// State
	info := ScrapeInfo{}
	links := []string{}
	seen := map[string]bool{}

	// Create a collector specifically for this domain
	c := NewCollector(domainUrl)

	// Use a random delay to hopefully not get blocked by domain
	c.Limit(&colly.LimitRule{RandomDelay: 5 * time.Second})

	// Scrape body text content, after culling script and style tags
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Scrape div text content, ignoring script, style, and random img tags
		//
		// Answer to my prayers provided by this kind soul from Stack Overflow
		//
		// https://stackoverflow.com/questions/44441665/how-to-extract-only-text-from-html-in-golang#answer-44444296
		e.DOM.Find("script,style").Each(func(i int, el *goquery.Selection) {
			el.Remove()
		})

		textContent := e.DOM.Text()
		info.PhoneNumbers = append(info.PhoneNumbers, phone.MatchPhoneNumbers(textContent)...)
	})

	// Get all of the link hrefs from each nav element
	c.OnHTML("nav", func(e *colly.HTMLElement) {
		navLinks := e.ChildAttrs("a", "href")
		for _, link := range navLinks {
			_, found := seen[link]
			if !found {
				seen[link] = true
				links = append(links, link)
			}
		}
	})

	// Collect phone numbers from a[href="tel:"]
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// Check if we have any links with a[href="tel:< phone number >"]
		if strings.HasPrefix(href, hrefPrefix) {
			tel := strings.TrimPrefix(href, hrefPrefix)
			info.PhoneNumbers = append(info.PhoneNumbers, *phone.NewFromHrefTel(tel))
		}
	})

	// After we scraped each page, check we see if we gathered enough info.
	c.OnScraped(func(r *colly.Response) {
		info.LinksVisited = append(info.LinksVisited, r.Request.URL.String())

		if !info.EnoughInfo() && !info.ExceededPageLimit() {
			var nextLink string

			// Try contact page, if it is available
			links, nextLink = spliceLink(links, "contact")
			if nextLink != "" {
				c.Visit(nextLink)
				return
			}

			// Try about page, if it is available
			links, nextLink = spliceLink(links, "about")
			if nextLink != "" {
				c.Visit(nextLink)
				return
			}

			// Try any other page, if available
			if len(links) > 0 {
				links, nextLink = links[1:], links[0]
				c.Visit(nextLink)
				return
			}
		}
	})

	// Log each visited endpoint
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %q\n", r.URL.String())
	})

	// Do the thing! (visit domain and start scraping)
	err = c.Visit(domainUrl.String())

	// Wait for collector jobs to return, in case we choose to use async
	c.Wait()

	// Sanitize gathered information
	info.SanitizePhoneNumbers()

	return &info, err
}

// spliceLink returns the link containing the provided string,
// and removes it from the links slice
func spliceLink(links []string, needle string) (remaining []string, found string) {
	for index, link := range links {
		if strings.Contains(link, needle) {
			// https://github.com/golang/go/wiki/SliceTricks#delete
			links = append(links[:index], links[index+1:]...)
			return links, link
		}
	}

	return links, ""
}

type handleScrapeResult func(s *ScrapeJobResult)

// ScrapeDomains will scrape domains for information using `numWorkers` goroutines.
// Each ScrapeResult is passed to the handleScrapeResult function.
func ScrapeDomains(urls []string, numWorkers int, handleResult handleScrapeResult) {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	var wg sync.WaitGroup

	// Channel on which jobs are enqueued
	jobCh := make(chan domainJob, len(urls))

	// Channel on which results will be received
	resultCh := make(chan ScrapeJobResult, len(urls))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Process each check url job
			for job := range jobCh {
				result, err := ScrapeDomain(job.url)
				resultCh <- ScrapeJobResult{Url: job.url, Info: *result, Err: err}
			}
		}()
	}

	// Enqueue jobs
	for index, url := range urls {
		jobCh <- domainJob{index: index, url: url}
	}
	close(jobCh)

	// Once all workers complete their jobs, close the result channel
	// to signal the top level goroutine no more results will be received
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		if handleResult != nil {
			handleResult(&result)
		}
	}
}
