package utils

import (
	"errors"
	"regexp"
)

const (
	EmailNotValidError = "email is not valid"
	GroupNameNotValid  = "group name is not valid"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New(EmailNotValidError)
	}
	return nil
}

func ValidateNewGroup(name string) error {
	if len(name) > 300 {
		return errors.New(GroupNameNotValid)
	}

	return nil
}
