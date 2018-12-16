// Package strs contains common string manipulation functionality.
package strs

import "unicode/utf8"

// IsUpperCamelCase returns true if s is not empty and is camel case with an initial capital.
func IsUpperCamelCase(s string) bool {
	if !isCapitalized(s) {
		return false
	}
	return isCamelCase(s)
}

// isCapitalized returns true if is not empty and the first letter is
// an uppercase character.
func isCapitalized(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return isUpper(r)
}

// isCamelCase returns false if s is empty or contains any character that is
// not between 'A' and 'Z', 'a' and 'z', '0' and '9', or in extraRunes.
// It does not care about lowercase or uppercase.
func isCamelCase(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !(isLetter(c) || isDigit(c)) {
			return false
		}
	}
	return true
}

func isLetter(r rune) bool {
	return isUpper(r) || isLower(r)
}

func isLower(r rune) bool {
	return 'a' <= r && r <= 'z'
}

func isUpper(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
