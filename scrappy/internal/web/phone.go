package web

import (
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type PhoneNumberConfidence int

const (
	PhoneRegexMatch PhoneNumberConfidence = iota
	PhoneHrefTel
)

const hrefPrefix = "tel:"

type Phone struct {
	Number string
	// Depending on how we scraped the phone number,
	// we can have more or less confidence that it is valid.
	//
	// Example: a[href] of type `tel` should be a phone number,
	// 			while a regex match is less certain
	Confidence PhoneNumberConfidence
}

// Implements fmt.Stringer
func (c PhoneNumberConfidence) String() string {
	switch c {
	case PhoneHrefTel:
		return "a[href=\"tel:< phone number >\"]"
	case PhoneRegexMatch:
		return "regex match"
	default:
		return "unknown"
	}
}

func GetPhoneNums(rawUrl string) (phoneNums []Phone, err error) {
	domain, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	// Create a collector specifically for this domain
	c := colly.NewCollector(colly.AllowedDomains(domain.Host))

	// Use a random delay to hopefully not get blocked by domain
	c.Limit(&colly.LimitRule{RandomDelay: 5 * time.Second})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// Check if we have any links with a[href="tel:< phone number >"]
		if strings.HasPrefix(href, hrefPrefix) {
			tel := strings.TrimPrefix(href, hrefPrefix)

			phoneNums = append(phoneNums, Phone{
				Number:     strings.TrimSpace(tel),
				Confidence: PhoneHrefTel,
			})
		}
	})

	// Do the thing! (visit domain and start scraping)
	err = c.Visit(domain.String())

	return phoneNums, err
}
