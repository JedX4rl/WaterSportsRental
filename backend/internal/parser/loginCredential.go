package parser

import (
	"regexp"
)

func IsUserValid(password string) error {

	if len(password) < 8 {
		return ErrInvalidPassword
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial := regexp.MustCompile(`[#!@$%^&*-]`).MatchString

	if hasUpper(password) && hasLower(password) && hasDigit(password) && hasSpecial(password) {
		return nil
	}
	return ErrInvalidPassword
}
