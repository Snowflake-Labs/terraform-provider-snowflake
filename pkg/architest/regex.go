package architest

import "regexp"

var (
	AcceptanceTestFileRegex  = regexp.MustCompile("^.*_acceptance_test.go$")
	AcceptanceTestNameRegex  = regexp.MustCompile("^TestAcc_.+$")
	IntegrationTestFileRegex = regexp.MustCompile("^.*_integration_test.go$")
	IntegrationTestNameRegex = regexp.MustCompile("^TestInt_.+$")
	TestFileRegex            = regexp.MustCompile("^.*_test.go$")
	TestNameRegex            = regexp.MustCompile("^Test.*$")
)
