package utils_test

import (
	"avito/pkg/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSlug_ValidSlug(t *testing.T) {
	tests := []struct {
		slug string
	}{
		{"valid-slug"},
		{"validslug"},
		{"AVITO_VOICE_MESSAGES"},
		{"AVITO_PERFORMANCE_VAS"},
		{"AVITO_DISCOUNT_30"},
		{"AVITO_DISCOUNT_50"},
	}

	for _, test := range tests {
		err := utils.ValidateSlug(test.slug)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

func TestValidateSlug_InvalidSlug(t *testing.T) {
	tests := []struct {
		slug        string
		expectedErr error
	}{
		{"invalid slug", errors.New("invalid slug")},
		{"", errors.New("invalid slug")},
		{"$invalid_slug", errors.New("invalid slug")},
		{" invalidslug", errors.New("invalid slug")},
		{"invalid-slug.", errors.New("invalid slug")},
		{"invalid_slug/", errors.New("invalid slug")},
		{"invalidslug-", errors.New("invalid slug")},
		{"-invalidslug", errors.New("invalid slug")},
		{"-invalidslug_", errors.New("invalid slug")},
	}

	for _, test := range tests {
		err := utils.ValidateSlug(test.slug)
		if err == nil {
			t.Error("Expected an error, but got nil")
		} else if err.Error() != test.expectedErr.Error() {
			t.Errorf("Expected error message '%s', but got '%s'", test.expectedErr.Error(), err.Error())
		}
	}
}

func TestValidateYearMonth_Valid(t *testing.T) {
	tests := []struct {
		yearMonth string
	}{
		{"2023-08"},
		{"1970-01"},
	}

	for _, test := range tests {
		err := utils.ValidateYearMonth(test.yearMonth)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

func TestValidateYearMonth_Invalid(t *testing.T) {
	tests := []struct {
		yearMonth   string
		expectedErr error
	}{
		{"2023-8", errors.New("invalid YearMonth")},
		{"2023/08", errors.New("invalid YearMonth")},
		{"2023.8", errors.New("invalid YearMonth")},
		{"23-08", errors.New("invalid YearMonth")},
		{"23/8", errors.New("invalid YearMonth")},
		{"23/08", errors.New("invalid YearMonth")},
		{"08/2023", errors.New("invalid YearMonth")},
		{"8/2023", errors.New("invalid YearMonth")},
		{"8/23", errors.New("invalid YearMonth")},
	}

	for _, test := range tests {
		err := utils.ValidateYearMonth(test.yearMonth)
		if err == nil {
			t.Error("Expected an error, but got nil")
		} else if err.Error() != test.expectedErr.Error() {
			t.Errorf("Expected error message '%s', but got '%s'", test.expectedErr.Error(), err.Error())
		}
	}
}

func TestProbability(t *testing.T) {
	tests := []struct {
		Slug       string
		UserId     int
		Percentage int
		Expected   bool
	}{
		{"example", 1, 50, true},
		{"example", 2, 50, true},
		{"example", 3, 50, true},
		{"example", 4, 50, true},
		{"example", 5, 50, true},
		{"example", 6, 50, false},
		{"example", 7, 50, true},
		{"example", 8, 50, false},
		{"AVITO_VOICE_MESSAGES", 275456710, 17, false},
		{"AVITO_PERFORMANCE_VAS", 12345678910, 85, true},
		{"AVITO_DISCOUNT_30", 10987654231, 46, false},
		{"AVITO_DISCOUNT_50", 88005553535, 60, true},
		{"", 0, -1, false},
		{"", 0, 101, true},
	}
	_ = tests
	for _, testCase := range tests {
		assert.Equal(t, testCase.Expected, utils.Probability(testCase.Slug, int64(testCase.UserId), testCase.Percentage))
	}
}
