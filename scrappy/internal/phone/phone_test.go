package phone

import (
	"errors"
	"reflect"
	"testing"
)

func TestMatchPhoneNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []Phone
	}{
		{
			name: "number with phone prefix",
			text: "Phone: +1 110-112-1111",
			expected: []Phone{
				{Number: "+1 110-112-1111", Confidence: PhoneRegexMatchWithPrefix},
			},
		},
		{
			name: "number with telephone prefix",
			text: "Telephone: 111-110-1011",
			expected: []Phone{
				{Number: "111-110-1011", Confidence: PhoneRegexMatchWithPrefix},
			},
		},
		{
			name: "number without phone prefix",
			text: " +1 110-112-1111",
			expected: []Phone{
				{Number: "+1 110-112-1111", Confidence: PhoneRegexMatch},
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
			expected: []Phone{
				{Number: "+1 110-112-1111", Confidence: PhoneRegexMatch},
				{Number: "+1 110-112-1112", Confidence: PhoneRegexMatchWithPrefix},
				{Number: "110-112-1113", Confidence: PhoneRegexMatch},
				{Number: "110.112.1115", Confidence: PhoneRegexMatchWithPrefix},
				{Number: "+1 110 112 1118", Confidence: PhoneRegexMatch},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			phoneNumbers := MatchPhoneNumbers(tc.text)

			if !reflect.DeepEqual(phoneNumbers, tc.expected) {
				t.Errorf("Expected %+v, got %+v instead", tc.expected, phoneNumbers)
			}
		})
	}
}

func TestValidatePhoneNumbers(t *testing.T) {
	phoneNums := []Phone{
		// Example valid US phone number from https://stdcxx.apache.org/doc/stdlibug/26-1.html
		{Number: "(541) 754-3010"},
		{Number: "(555) 555-5555"},
		{Number: "5417543012"},
		{Number: "541.754.3013"},
		{Number: "not really a phone number, lol"},
	}

	expectedValid := []Phone{
		{Number: "+1 541-754-3010"},
		{Number: "+1 541-754-3012"},
		{Number: "+1 541-754-3013"},
	}
	expectedFailed := []FailedValidation{
		{Index: 1, Number: "(555) 555-5555", Err: ErrInvalidNumber},
		{Index: 4, Number: "not really a phone number, lol", Err: ErrInvalidNumber},
	}

	validated, failed := ValidatePhoneNumbers(phoneNums)

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
		phone    Phone
		expected Phone
	}{
		{
			// Example valid US phone number from https://stdcxx.apache.org/doc/stdlibug/26-1.html
			name:     "valid US phone number",
			phone:    Phone{Number: "(541) 754-3010"},
			expected: Phone{Number: "+1 541-754-3010"},
		},
		{
			name:     "valid internationally formatted US number",
			phone:    Phone{Number: "+1 541 754-3010"},
			expected: Phone{Number: "+1 541-754-3010"},
		},
		{
			name:     "dot separated phone number",
			phone:    Phone{Number: "541.754.3010"},
			expected: Phone{Number: "+1 541-754-3010"},
		},
		{
			name:     "no separators phone number",
			phone:    Phone{Number: "5417543010"},
			expected: Phone{Number: "+1 541-754-3010"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePhoneNumber(&tc.phone)
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
		phone       Phone
		expectedErr error
	}{
		{
			name:        "invalid formatted number",
			phone:       Phone{Number: "123213451234531245"},
			expectedErr: ErrInvalidNumber,
		},
		{
			name:        "invalid correctly formatted number",
			phone:       Phone{Number: "+1 555-555-5555"},
			expectedErr: ErrInvalidNumber,
		},
		{
			name:        "invalid no separators phone number",
			phone:       Phone{Number: "5555555555"},
			expectedErr: ErrInvalidNumber,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidatePhoneNumber(&tc.phone)
			checkErrIs(t, err, tc.expectedErr)
		})
	}
}

func TestDedupPhoneNumbers(t *testing.T) {
	phoneNums := []Phone{
		{Number: "111.111.1111", Confidence: PhoneRegexMatch},
		{Number: "333.333.3333"},
		{Number: "222.222.2222"},
		{Number: "111.111.1111", Confidence: PhoneHrefTel},
		{Number: "222.222.2222", Confidence: PhoneRegexMatchWithPrefix},
		{Number: "333.333.3333"},
		{Number: "111.111.1111"},
	}

	expected := []Phone{
		{Number: "111.111.1111", Confidence: PhoneHrefTel},
		{Number: "333.333.3333"},
		{Number: "222.222.2222", Confidence: PhoneRegexMatchWithPrefix},
	}

	result := DedupPhoneNumbers(phoneNums)

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
