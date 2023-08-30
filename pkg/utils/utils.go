package utils

import (
	"errors"
	"regexp"
)

func ValidateSlug(slug string) error {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+(?:[-_][a-zA-Z0-9]+)*$", slug)
	if !match {
		return errors.New("invalid slug")
	}
	return nil
}

func ValidateYearMonth(yearMonth string) error {
	match, _ := regexp.MatchString("^\\d{4}-\\d{2}$", yearMonth)
	if !match {
		return errors.New("invalid YearMonth")
	}
	return nil
}
