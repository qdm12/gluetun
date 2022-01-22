package pprof

import (
	"regexp"
	"time"

	gomock "github.com/golang/mock/gomock"
)

func boolPtr(b bool) *bool                       { return &b }
func stringPtr(s string) *string                 { return &s }
func durationPtr(d time.Duration) *time.Duration { return &d }

var _ gomock.Matcher = (*regexMatcher)(nil)

type regexMatcher struct {
	regexp *regexp.Regexp
}

func (r *regexMatcher) Matches(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return false
	}
	return r.regexp.MatchString(s)
}

func (r *regexMatcher) String() string {
	return "regular expression " + r.regexp.String()
}

func newRegexMatcher(regex string) *regexMatcher {
	return &regexMatcher{
		regexp: regexp.MustCompile(regex),
	}
}
