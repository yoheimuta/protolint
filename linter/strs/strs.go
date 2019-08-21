// Package strs contains common string manipulation functionality.
package strs

import (
	"strings"
	"unicode/utf8"
)

// IsUpperCamelCase returns true if s is not empty and is camel case with an initial capital.
func IsUpperCamelCase(s string) bool {
	if !isCapitalized(s) {
		return false
	}
	return isCamelCase(s)
}

// IsUpperSnakeCase returns true if s only contains uppercase letters,
// digits, and/or underscores. s MUST NOT begin or end with an underscore.
func IsUpperSnakeCase(s string) bool {
	if s == "" || s[0] == '_' || s[len(s)-1] == '_' {
		return false
	}
	for _, r := range s {
		if !(isUpper(r) || isDigit(r) || r == '_') {
			return false
		}
	}
	return true
}

// IsLowerSnakeCase returns true if s only contains lowercase letters,
// digits, and/or underscores. s MUST NOT begin or end with an underscore.
func IsLowerSnakeCase(s string) bool {
	if s == "" || s[0] == '_' || s[len(s)-1] == '_' {
		return false
	}
	for _, r := range s {
		if !(isLower(r) || isDigit(r) || r == '_') {
			return false
		}
	}
	return true
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

// isSnake returns true if s only contains letters, digits, and/or underscores.
// s MUST NOT begin or end with an underscore.
func isSnake(s string) bool {
	if s == "" || s[0] == '_' || s[len(s)-1] == '_' {
		return false
	}
	for _, r := range s {
		if !(isLetter(r) || isDigit(r) || r == '_') {
			return false
		}
	}
	return true
}

// IsLowerCase returns true if s only contain characters in the range a-z0-9.
func IsLowerCase(s string) bool {
	for _, r := range s {
		if !(isLower(r) || isDigit(r)) {
			return false
		}
	}
	return true
}

// toSnake converts s to snake_case.
func toSnake(s string) string {
	output := ""
	s = strings.TrimSpace(s)
	priorLower := false
	for _, c := range s {
		if priorLower && isUpper(c) {
			output += "_"
		}
		output += string(c)
		priorLower = isLower(c)
	}
	return output
}

// SplitCamelCaseWord splits a CamelCase word into its parts.
//
// If s is empty, returns nil.
// If s is not CamelCase, returns nil.
func SplitCamelCaseWord(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.TrimSpace(s)
	if !isCamelCase(s) {
		return nil
	}
	return SplitSnakeCaseWord(toSnake(s))
}

// SplitSnakeCaseWord splits a snake_case word into its parts.
//
// If s is empty, returns nil.
// If s is not snake_case, returns nil.
func SplitSnakeCaseWord(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.TrimSpace(s)
	if !isSnake(s) {
		return nil
	}
	return strings.Split(s, "_")
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
