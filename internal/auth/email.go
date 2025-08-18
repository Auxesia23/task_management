package auth

import (
	"net"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func doesDomainExist(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	mx, err := net.LookupMX(domain)
	if err != nil {
		return false
	}
	return len(mx) > 0
}

func ValidateEmail(email string) bool {
	// Step 1: Regex format check
	if !isEmailValid(email) {
		return false
	}

	// Step 2: DNS MX record check
	if !doesDomainExist(email) {
		return false
	}

	return true
}
