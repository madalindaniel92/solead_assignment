// Package web handles http requests and data scraping.
package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

const defaultTimeout time.Duration = 10 * time.Second

// NewClient returns a new HTTP client with the given timeout.
//
// If timeout is 0, it will use the `defaultTimeout`
func NewClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &http.Client{Timeout: timeout}
}

// CheckURL send an http HEAD request to the url to check if it is reachable.
func CheckURL(url string) (status int, err error) {
	response, err := NewClient(defaultTimeout).Head(url)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}

// domainJob represents a job for each worker processing a different domain.
type domainJob struct {
	index int
	url   string
}

// CheckUrlResult represents the result of each worker running CheckURL
type CheckUrlResult struct {
	job    domainJob
	Status int
	Err    error
}

// URL returns the corresponding url of this job
func (c *CheckUrlResult) URL() string {
	return c.job.url
}

type checkUrlCallback func(c *CheckUrlResult)

// CheckURLs will check urls through http head requests using `numWorkers` goroutines.
func CheckURLs(urls []string, numWorkers int, handleResult checkUrlCallback) []CheckUrlResult {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	var wg sync.WaitGroup

	// Slice in which results are collected
	results := make([]CheckUrlResult, len(urls))

	// Channel on which jobs are enqueued
	jobCh := make(chan domainJob, len(urls))

	// Channel on which results will be received
	resultCh := make(chan CheckUrlResult, len(urls))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Process each check url job
			for job := range jobCh {
				status, err := CheckURL(job.url)
				resultCh <- CheckUrlResult{job: job, Status: status, Err: err}
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
		results[result.job.index] = result

		if handleResult != nil {
			handleResult(&result)
		}
	}

	return results
}

// FilterSuccessfulDomains returns the URLs for the domains that responded with HEAD 200.
func FilterSuccessfulDomains(results []CheckUrlResult) []string {
	successful := make([]string, 0, len(results))

	for _, result := range results {
		if result.Status == http.StatusOK {
			successful = append(successful, result.URL())
		}
	}

	return successful
}

// NewCollector returns a new colly Collector with default settings applied.
func NewCollector(domain *url.URL) *colly.Collector {
	return colly.NewCollector(colly.AllowedDomains(allowedDomains(domain)...))
}

func allowedDomains(url *url.URL) []string {
	host := url.Host
	allowed := []string{url.Host}

	// If host does not have "www." subdomain, add it to allowed hosts.
	// This will permit sites to redirect to their "www." subdomain when scraping
	if !strings.HasPrefix("www.", host) {
		allowed = append(allowed, fmt.Sprintf("www.%s", host))
	}

	return allowed
}
