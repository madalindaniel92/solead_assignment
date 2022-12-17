package web

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type PhoneNumberConfidence int

const (
	// Phone number extracted using a regex
	PhoneRegexMatch PhoneNumberConfidence = iota
	// Phone number extracted using a regex, having the word "phone" as a prefix
	PhoneRegexMatchWithPrefix
	// Phone number extracted from a[href] with 'tel:' prefix
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
	case PhoneRegexMatchWithPrefix:
		return "regex match with 'phone' prefix"
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
		phoneNums = append(phoneNums, MatchPhoneNumbers(textContent)...)
	})

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

// Test text to see if it says "phone" somewhere close to the regexp match
var phonePrefixRegex = regexp.MustCompile(`(?i)\b(phone|telephone)\b`)

// Regex to match US phone numbers - seems legit!
//
// https://stackoverflow.com/questions/16699007/regular-expression-to-match-standard-10-digit-phone-number#answer-16699507
var usPhoneNumberRegex = regexp.MustCompile(`(\+\d{1,2}\s)?[\s.-]*\(?\d{3}\)?[\s.-]*\d{3}[\s.-]*\d{4}`)

// MatchPhoneNumbers uses a regex to scrape the text for viable phone numbers.
//
// If we have a prefix like "phone" or "telephone" before the number, we'll set
// confidence to PhoneRegexMatchWithPrefix, otherwise it is PhoneRegexMatch
func MatchPhoneNumbers(text string) []Phone {
	var phoneNums []Phone

	matchIndexes := usPhoneNumberRegex.FindAllStringIndex(text, -1)
	for _, matchIndex := range matchIndexes {
		startIndex, endIndex := matchIndex[0], matchIndex[1]
		number := text[startIndex:endIndex]

		// Grab 12 characters before the match, to see if they match "phone"
		prefix := text[max(0, startIndex-12):startIndex]
		confidence := PhoneRegexMatch

		// If before the phone number we have a prefix like "phone" or "telephone"
		// our confidence that this is in fact a phone number increases
		if phonePrefixRegex.MatchString(prefix) {
			confidence = PhoneRegexMatchWithPrefix
		}

		phoneNums = append(phoneNums, Phone{
			Number:     strings.TrimSpace(number),
			Confidence: confidence,
		})
	}

	return phoneNums
}
