package phone

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// We'll handle US phone numbers for now
const defaultPhoneRegion = "US"

type PhoneNumberConfidence int

const (
	// Phone number extracted using a regex
	PhoneRegexMatch PhoneNumberConfidence = iota
	// Phone number extracted using a regex, having the word "phone" as a prefix
	PhoneRegexMatchWithPrefix
	// Phone number extracted from a[href] with 'tel:' prefix
	PhoneHrefTel
)

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

type FailedValidation struct {
	Index  int
	Number string
	Err    error
}

// ValidatePhoneNumbers validates and returns the valid phone numbers.
//
// Invalid numbers are returned in the second return value.
//
// Returned phone numbers are formatted using the international phone number scheme.
func ValidatePhoneNumbers(phoneNums []Phone) (valid []Phone, invalid []FailedValidation) {
	for index, phone := range phoneNums {
		result, err := ValidatePhoneNumber(&phone)
		if err != nil {
			invalid = append(invalid, FailedValidation{
				Index:  index,
				Number: phone.Number,
				Err:    err,
			})
			continue
		}

		valid = append(valid, *result)
	}

	return valid, invalid
}

func ValidatePhoneNumberString(number string) (*Phone, error) {
	phone := Phone{Number: number}
	return ValidatePhoneNumber(&phone)
}

// ValidatePhoneNumber validates phonenumbers, assuming they follow the
// North American Numbering Plan and are US numbers.
//
// The numbers will be formatted using the national US phone number scheme.
func ValidatePhoneNumber(phone *Phone) (*Phone, error) {
	result, err := phonenumbers.Parse(phone.Number, defaultPhoneRegion)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNumber, err)
	}

	if !phonenumbers.IsValidNumber(result) {
		return nil, ErrInvalidNumber
	}

	phone.Number = phonenumbers.Format(result, phonenumbers.INTERNATIONAL)
	return phone, nil
}

// DedupPhoneNumbers deduplicates phone numbers, keeping the option with the highest confidence.
func DedupPhoneNumbers(phoneNums []Phone) []Phone {
	seen := map[string]PhoneNumberConfidence{}
	results := []Phone{}

	for _, phone := range phoneNums {
		confidence, found := seen[phone.Number]
		if found {
			// We've already seen this number, so we compare the
			// current confidence with the old and keep the max
			if phone.Confidence > confidence {
				seen[phone.Number] = phone.Confidence
			}
		} else {
			// First time we see this number
			results = append(results, phone)
			seen[phone.Number] = phone.Confidence
		}
	}

	// Traverse results again and set confidence to the
	// largest value we've seen for each number
	for index := range results {
		result := &results[index] // we need to update through a reference
		result.Confidence = seen[result.Number]
	}

	return results
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewFromHrefTel returns a new Phone with confidence set to PhoneHrefTel.
func NewFromHrefTel(tel string) *Phone {
	return &Phone{
		Number:     strings.TrimSpace(tel),
		Confidence: PhoneHrefTel,
	}
}
