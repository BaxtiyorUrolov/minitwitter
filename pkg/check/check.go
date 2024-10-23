package check

import (
	"unicode"
)

func PhoneNumber(phone string) bool {
	for _, r := range phone {
		if r == '+' {
			continue
		} else if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		}
		if unicode.IsLower(char) {
			hasLowerCase = true
		}
	}

	return hasUpperCase && hasLowerCase
}
