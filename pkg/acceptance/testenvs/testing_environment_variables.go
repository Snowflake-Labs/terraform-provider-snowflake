package testenvs

import (
	"fmt"
	"os"
	"testing"
)

type env string

const (
	User     env = "TEST_SF_TF_USER"
	Password env = "TEST_SF_TF_PASSWORD"
	Account  env = "TEST_SF_TF_ACCOUNT"
	Role     env = "TEST_SF_TF_ROLE"
	Host     env = "TEST_SF_TF_HOST"
)

// TODO: test
func GetOrSkipTest(t *testing.T, envName Env) string {
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if env == "" {
		t.Skipf("Skipping %s, env %v missing", t.Name(), envName)
	}
	return env
}

type Env interface {
	xxxProtected()
}

func (e env) xxxProtected() {}
