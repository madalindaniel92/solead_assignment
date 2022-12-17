package web

import (
	"net/url"
	"strings"
	"time"

	"examples/scrappy/internal/phone"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const hrefPrefix = "tel:"

func GetPhoneNums(rawUrl string) (phoneNums []phone.Phone, err error) {
	domain, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	// Create a collector specifically for this domain
	c := colly.NewCollector(colly.AllowedDomains(domain.Host))

	// Use a random delay to hopefully not get blocked by domain
	c.Limit(&colly.LimitRule{RandomDelay: 5 * time.Second})

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
		phoneNums = append(phoneNums, phone.MatchPhoneNumbers(textContent)...)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// Check if we have any links with a[href="tel:< phone number >"]
		if strings.HasPrefix(href, hrefPrefix) {
			tel := strings.TrimPrefix(href, hrefPrefix)
			phoneNums = append(phoneNums, *phone.NewFromHrefTel(tel))
		}
	})

	// Do the thing! (visit domain and start scraping)
	err = c.Visit(domain.String())

	return phoneNums, err
}
