package shortgen

import (
	"fmt"
	"regexp"
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	shortLength = 10
	base        = 63
)

func EncodeID(id uint64) (string, error) {
	if id == 0 {
		return "", fmt.Errorf("id is 0")
	}

	b := make([]byte, shortLength)
	for i := shortLength - 1; i >= 0; i-- {
		b[i] = charset[id%base]
		id /= base
	}

	if id > 0 {
		return "", fmt.Errorf("id is too large for 10 characters")
	}

	return string(b), nil
}

var shortPattern = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9_]{%d}$`, shortLength))

func Validate(short string) error {
	if !shortPattern.MatchString(short) {
		return fmt.Errorf("short must be exactly %d characters of [a-zA-Z0-9_]", shortLength)
	}
	return nil
}
