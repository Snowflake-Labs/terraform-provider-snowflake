// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package testhelpers

import (
	"os"
	"path/filepath"
	"testing"
)

func MustFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := Fixture(name)
	if err != nil {
		t.Error(err)
	}
	return b
}

func Fixture(name string) (string, error) {
	b, err := os.ReadFile(filepath.Join("testdata", name))
	return string(b), err
}
