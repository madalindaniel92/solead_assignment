package web_test

import (
	"errors"
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

func TestValidatePhoneNumbers(t *testing.T) {
	phoneNums := []web.Phone{
		// Example valid US phone number from https://stdcxx.apache.org/doc/stdlibug/26-1.html
		{Number: "(541) 754-3010"},
		{Number: "(555) 555-5555"},
		{Number: "5417543012"},
		{Number: "541.754.3013"},
		{Number: "not really a phone number, lol"},
	}

	expectedValid := []web.Phone{
		{Number: "+1 541-754-3010"},
		{Number: "+1 541-754-3012"},
		{Number: "+1 541-754-3013"},
	}
	expectedFailed := []web.FailedValidation{
		{Index: 1, Number: "(555) 555-5555", Err: web.ErrInvalidPhoneNumber},
		{Index: 4, Number: "not really a phone number, lol", Err: web.ErrInvalidPhoneNumber},
	}

	validated, failed := web.ValidatePhoneNumbers(phoneNums)

	// Check successfully validated
	if !reflect.DeepEqual(validated, expectedValid) {
		t.Errorf("Expected %+v. got %+v instead", expectedValid, validated)
	}

	// Check validation failures
	// (can't use `DeepEqual` because of checkErrIs check)
	if len(failed) != len(expectedFailed) {
		t.Fatalf("Expected %d failed validations. got %d instead",
			len(failed), len(expectedFailed))
	}

	for index, expected := range expectedFailed {
		actual := failed[index]

		if actual.Index != expected.Index || actual.Number != expected.Number {
			t.Errorf("Expected %+v. got %+v instead", expected, actual)
		}

		checkErrIs(t, actual.Err, expected.Err)
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	testCases := []struct {
		name     string
		phone    web.Phone
		expected web.Phone
	}{
		{
			// Example valid US phone number from https://stdcxx.apache.org/doc/stdlibug/26-1.html
			name:     "valid US phone number",
			phone:    web.Phone{Number: "(541) 754-3010"},
			expected: web.Phone{Number: "+1 541-754-3010"},
		},
		{
			name:     "valid internationally formatted US number",
			phone:    web.Phone{Number: "+1 541 754-3010"},
			expected: web.Phone{Number: "+1 541-754-3010"},
		},
		{
			name:     "dot separated phone number",
			phone:    web.Phone{Number: "541.754.3010"},
			expected: web.Phone{Number: "+1 541-754-3010"},
		},
		{
			name:     "no separators phone number",
			phone:    web.Phone{Number: "5417543010"},
			expected: web.Phone{Number: "+1 541-754-3010"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := web.ValidatePhoneNumber(&tc.phone)
			checkNoErr(t, err)

			if *result != tc.expected {
				t.Errorf("Expected %+v, got %+v instead", tc.expected, result)
			}
		})
	}
}

func TestValidatePhoneNumber_failure(t *testing.T) {
	testCases := []struct {
		name        string
		phone       web.Phone
		expectedErr error
	}{
		{
			name:        "invalid formatted number",
			phone:       web.Phone{Number: "123213451234531245"},
			expectedErr: web.ErrInvalidPhoneNumber,
		},
		{
			name:        "invalid correctly formatted number",
			phone:       web.Phone{Number: "+1 555-555-5555"},
			expectedErr: web.ErrInvalidPhoneNumber,
		},
		{
			name:        "invalid no separators phone number",
			phone:       web.Phone{Number: "5555555555"},
			expectedErr: web.ErrInvalidPhoneNumber,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := web.ValidatePhoneNumber(&tc.phone)
			checkErrIs(t, err, tc.expectedErr)
		})
	}
}

func TestDedupPhoneNumbers(t *testing.T) {
	phoneNums := []web.Phone{
		{Number: "111.111.1111", Confidence: web.PhoneRegexMatch},
		{Number: "333.333.3333"},
		{Number: "222.222.2222"},
		{Number: "111.111.1111", Confidence: web.PhoneHrefTel},
		{Number: "222.222.2222", Confidence: web.PhoneRegexMatchWithPrefix},
		{Number: "333.333.3333"},
		{Number: "111.111.1111"},
	}

	expected := []web.Phone{
		{Number: "111.111.1111", Confidence: web.PhoneHrefTel},
		{Number: "333.333.3333"},
		{Number: "222.222.2222", Confidence: web.PhoneRegexMatchWithPrefix},
	}

	result := web.DedupPhoneNumbers(phoneNums)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v. got %+v instead", expected, result)
	}
}

// Helpers

func checkNoErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error: %q", err)
	}
}

func checkErrIs(t *testing.T, err error, expected error) {
	t.Helper()

	if !errors.Is(err, expected) {
		t.Fatalf("Expected error %q, got %q (%T) instead", expected, err, err)
	}
}
