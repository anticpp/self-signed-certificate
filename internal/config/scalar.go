package config

import (
	"regexp"
	"strconv"
)

var reStrQuote *regexp.Regexp
var reInt *regexp.Regexp
var reFloat *regexp.Regexp

func init() {
	// Go's regexp user RE2's engine that doesn't support backreferences.
	// See https://swtch.com/~rsc/regexp/regexp3.html

	// String
	// "dog", 'dog' , 'dog", "dog'
	reStrQuote = regexp.MustCompile(`^("|')(.*)("|')$`)

	// Int
	// 123, -123
	reInt = regexp.MustCompile(`^[+-]?[0-9]+$`)

	// Float
	// 123.12, 0.12, -123.12
	reFloat = regexp.MustCompile(`^[+-]?[0-9]+[.][0-9]+$`)
}

// Parse the scalar type.
//
// The `any` is returned the real data type.
// Quoted string will strip is't quotes.
// Int type will be parsed to int64
// Float type will be parsed to float64
func parseScalar(s string) any {
	var ok bool
	var matches []string

	matches = reStrQuote.FindStringSubmatch(s)
	if matches != nil {
		return matches[2]
	}

	if len(s) == 0 {
		return ""
	}

	ok = reInt.MatchString(s)
	if ok {
		rv, _ := strconv.ParseInt(s, 10, 64)
		return rv
	}

	ok = reFloat.MatchString(s)
	if ok {
		rv, _ := strconv.ParseFloat(s, 64)
		return rv
	}

	// Otherwise, it's a normal string
	return s
}
