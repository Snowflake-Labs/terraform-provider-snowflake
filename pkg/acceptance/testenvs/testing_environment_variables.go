package testenvs

import (
	"os"
	"testing"
)

const User = "TEST_SF_TF_USER"
const Password = "TEST_SF_TF_PASSWORD"
const Account = "TEST_SF_TF_ACCOUNT"
const Role = "TEST_SF_TF_ROLE"
const Host = "TEST_SF_TF_HOST"

// TODO: allow to be used only with the set above
// TODO: test
func GetOrSkipTest(t *testing.T, envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		t.Skipf("Skipping %s", t.Name())
	}
	return env
}
