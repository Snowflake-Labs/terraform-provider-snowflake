package testhelpers

import (
	"io/ioutil"
	"path/filepath"
)

func Fixture(name string) string {
	b, err := FixtureE(name)
	if err != nil {
		panic(err)
	}
	return b
}

func FixtureE(name string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	return string(b), err
}
