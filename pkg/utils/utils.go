package utils

import (
	"crypto/sha512"
	"errors"
	"fmt"
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

func Probability(segment string, number int64, percentage int) bool {
	if percentage <= 0 {
		return false
	} else if percentage >= 100 {
		return true
	}
	data := []byte(fmt.Sprintf("%s%d", segment, number))
	hash := sha512.Sum512(data)

	threshold := (512 * percentage) / 200

	return int(hash[0]) < threshold
}
