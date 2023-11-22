package archtests

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertPackage(t *testing.T, f *File, expectedPackage string) {
	assert.Equalf(t, expectedPackage, f.packageName, "filename %s has package %s, expected package %s", f.fileName, f.packageName, expectedPackage)
}

func assertAcceptanceTestNamedCorrectly(t *testing.T, f *File, methodName string) {
	assertMethodNameMatches(t, f, methodName, acceptanceTestNameRegex)
}

func assertMethodNameMatches(t *testing.T, f *File, methodName string, regex *regexp.Regexp) {
	assert.Truef(t, regex.MatchString(methodName), "filename %s contains exported method %s which does not match %s", f.fileName, methodName, regex.String())
}

func assertMethodNameDoesNotMatch(t *testing.T, f *File, methodName string, regex *regexp.Regexp) {
	assert.Falsef(t, regex.MatchString(methodName), "filename %s contains exported method %s which matches %s", f.fileName, methodName, regex.String())
}
