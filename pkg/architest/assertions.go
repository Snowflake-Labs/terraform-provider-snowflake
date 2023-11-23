package architest

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (f *File) AssertHasPackage(t *testing.T, expectedPackage string) {
	assert.Equalf(t, expectedPackage, f.packageName, "filename %s has package %s, expected package %s", f.FileName(), f.PackageName(), expectedPackage)
}

func (method *Method) AssertAcceptanceTestNamedCorrectly(t *testing.T) {
	method.AssertNameMatches(t, AcceptanceTestNameRegex)
}

func (method *Method) AssertTestNamedCorrectly(t *testing.T) {
	method.AssertNameMatches(t, TestNameRegex)
}

func (method *Method) AssertNameMatches(t *testing.T, regex *regexp.Regexp) {
	assert.Truef(t, regex.MatchString(method.Name()), "filename %s contains exported method %s which does not match %s", method.FileName(), method.Name(), regex.String())
}

func (method *Method) AssertNameDoesNotMatch(t *testing.T, regex *regexp.Regexp) {
	assert.Falsef(t, regex.MatchString(method.Name()), "filename %s contains exported method %s which matches %s", method.FileName(), method.Name(), regex.String())
}
