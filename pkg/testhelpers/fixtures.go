package testhelpers

import (
	"io/ioutil"
	"path/filepath"
)

func MustFixture(name string) string {
	b, err := Fixture(name)
	if err != nil {
		panic(err)
	}
	return b
}

func Fixture(name string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	return string(b), err
}
