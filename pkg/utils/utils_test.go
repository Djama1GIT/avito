package utils_test

import (
	"avito/pkg/utils"
	"errors"
	"testing"
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
