package web_test

import (
	"examples/scrappy/internal/web"
	"reflect"
	"testing"
)

func TestMatchPhoneNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []web.Phone
	}{
		{
			name: "number with phone prefix",
			text: "Phone: +1 110-112-1111",
			expected: []web.Phone{
				{Number: "+1 110-112-1111", Confidence: web.PhoneRegexMatchWithPrefix},
			},
		},
		{
			name: "number with telephone prefix",
			text: "Telephone: 111-110-1011",
			expected: []web.Phone{
				{Number: "111-110-1011", Confidence: web.PhoneRegexMatchWithPrefix},
			},
		},
		{
			name: "number without phone prefix",
			text: " +1 110-112-1111",
			expected: []web.Phone{
				{Number: "+1 110-112-1111", Confidence: web.PhoneRegexMatch},
			},
		},
		{
			name: "multiple phone numbers in the same text",
			text: `
				+1 110-112-1111
				Phone: +1 110-112-1112
				110-112-1113
				Telephone: 110.112.1115
				Call me maybe: +1 110 112 1118
			`,
			expected: []web.Phone{
				{Number: "+1 110-112-1111", Confidence: web.PhoneRegexMatch},
				{Number: "+1 110-112-1112", Confidence: web.PhoneRegexMatchWithPrefix},
				{Number: "110-112-1113", Confidence: web.PhoneRegexMatch},
				{Number: "110.112.1115", Confidence: web.PhoneRegexMatchWithPrefix},
				{Number: "+1 110 112 1118", Confidence: web.PhoneRegexMatch},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			phoneNumbers := web.MatchPhoneNumbers(tc.text)

			if !reflect.DeepEqual(phoneNumbers, tc.expected) {
				t.Errorf("Expected %+v, got %+v instead", tc.expected, phoneNumbers)
			}
		})
	}
}
