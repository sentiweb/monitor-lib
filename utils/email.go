package utils

import (
	"errors"
	"regexp"
)

const (
	// HeaderFrom is the From email header
	HeaderFrom = "From"

	// HeaderSubject is the Subject email header
	HeaderSubject = "Subject"

	// HeaderTo is the To email header
	HeaderTo = "To"

	// MimeTextPlain is the MIME type for plain text
	MimeTextPlain = "text/plain"
)

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ErrBadEmail = errors.New("invalid email format")
)

func IsEmailValid(email string) bool {
	return emailRegexp.MatchString(email)
}
