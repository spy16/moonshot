package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/mail"
)

const (
	CharsetNums     = "0123456789"
	CharsetLower    = "abcdefghijklmnopqrstuvwxyz"
	CharsetUpper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetLowerNum = CharsetLower + CharsetNums
	CharsetAlphaNum = CharsetLower + CharsetUpper + CharsetNums
)

// RandStr returns a random string of length 'n'. Characters are picked from
// the charset if provided or alphanumeric characters are used.
func RandStr(n int, charset ...string) string {
	chars := CharsetAlphaNum
	if len(charset) >= 1 {
		chars = charset[0]
	}

	s := make([]byte, n)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}

// GravatarURL returns a valid Gravatar URL for the given email id. Size parameter
// is passed to the URL to make thumbnail of given size.
func GravatarURL(email string, size int) string {
	if size <= 0 || size >= 2048 {
		size = 128
	}
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}

// IsValidEmail returns true if the given email is valid.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
