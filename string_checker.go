package validator

import "strings"

type SuffixChecker struct{}

type PrefixChecker struct{}


type StringChecker interface{
	IsValid(string,substring string)
}

// IsValid checks if the given 'substring' is not a suffix of the 'string'.
// It returns 'true' if the 'substring' is not a suffix, otherwise 'false'.
func (*SuffixChecker)IsValid (givenString,substring string) bool {
	return !strings.HasSuffix(givenString,substring)
}

// IsValid checks if the given 'substring' is not a prefix of the 'string'.
// It returns 'true' if the 'substring' is NOT a prefix, otherwise 'false'.
func (*PrefixChecker)IsValid (givenString,substring string) bool {
	return !strings.HasPrefix(givenString,substring)
}