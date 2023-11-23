package architest

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (f *File) AssertHasPackage(t *testing.T, expectedPackage string) {
	assert.Equalf(t, expectedPackage, f.packageName, "filename %s has package %s, expected package %s", f.fileName, f.packageName, expectedPackage)
}

func (method *Method) AssertAcceptanceTestNamedCorrectly(t *testing.T) {
	method.AssertNameMatches(t, AcceptanceTestNameRegex)
}

func (method *Method) AssertNameMatches(t *testing.T, regex *regexp.Regexp) {
	assert.Truef(t, regex.MatchString(method.name), "filename %s contains exported method %s which does not match %s", method.file.fileName, method.name, regex.String())
}

func (method *Method) AssertNameDoesNotMatch(t *testing.T, regex *regexp.Regexp) {
	assert.Falsef(t, regex.MatchString(method.name), "filename %s contains exported method %s which matches %s", method.file.fileName, method.name, regex.String())
}
