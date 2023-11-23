package architest

import "regexp"

var (
	AcceptanceTestFileRegex *regexp.Regexp
	AcceptanceTestNameRegex *regexp.Regexp
	TestFileRegex           *regexp.Regexp
	TestNameRegex           *regexp.Regexp
)

func init() {
	AcceptanceTestFileRegex = regexp.MustCompile("^.*_acceptance_test.go$")
	AcceptanceTestNameRegex = regexp.MustCompile("^TestAcc_.+$")
	TestFileRegex = regexp.MustCompile("^.*_test.go$")
	TestNameRegex = regexp.MustCompile("^Test.*$")
}
