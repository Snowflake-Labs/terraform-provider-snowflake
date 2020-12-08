package testhelpers

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func MustFixture(t *testing.T, name string) string {
	b, err := Fixture(name)
	if err != nil {
		t.Error(err)
	}
	return b
}

func Fixture(name string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	return string(b), err
}
