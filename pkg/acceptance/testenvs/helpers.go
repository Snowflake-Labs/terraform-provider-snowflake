package testenvs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertEnvNotSet(t *testing.T, envName string) {
	require.Empty(t, os.Getenv(envName))
}

func AssertEnvSet(t *testing.T, envName string) {
	require.NotEmpty(t, os.Getenv(envName))
}
