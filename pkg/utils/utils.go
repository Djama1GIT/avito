package utils

import (
	"errors"
	"regexp"
)

func ValidateSlug(slug string) error {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+(?:[-_][a-zA-Z0-9]+)*$", slug)
	if !match {
		return errors.New("Invalid slug")
	}
	return nil
}
