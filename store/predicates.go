package store

import (
	"bytes"
	"regexp"
)

type Predicate func(key []byte) bool

func PrefixMatcher(prefix string) Predicate {
	return func(k []byte) bool { return bytes.HasPrefix(k, []byte(prefix)) }
}

func RegexMatcher(pattern string) (Predicate, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return func(k []byte) bool { return r.Match(k) }, nil
}
