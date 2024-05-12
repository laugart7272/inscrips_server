package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUserName = regexp.MustCompile(`^[a-zA-Z\\s]`).MatchString
	isValidCWTitle  = regexp.MustCompile(`^[a-z0-9A-Z\\s]`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d charaters", minLength, maxLength)
	}

	return nil
}

func ValidateUserName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUserName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidateLastUserName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUserName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 2, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not valid email address")
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value < 1 {
		return fmt.Errorf("is not valid email id")
	}

	return nil
}

func ValidateSecretCode(value string) error {
	if err := ValidateString(value, 20, 128); err != nil {
		return err
	}

	return nil
}

func ValidateCientificWorkTitle(value string) error {
	if err := ValidateString(value, 10, 100); err != nil {
		return err
	}

	if !isValidCWTitle(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidateCientificWorkResume(value string) error {
	if err := ValidateString(value, 10, 100); err != nil {
		return err
	}

	if !isValidCWTitle(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}
