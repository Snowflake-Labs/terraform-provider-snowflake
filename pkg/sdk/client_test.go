package sdk

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientPing(t *testing.T) {
	// Secrets are required to run this test, so it is skipped by default.  To run it, set SKIP_SDK_TEST=false
	if os.Getenv("SKIP_SDK_TEST") != "false" {
		t.Skip("SKIP_SDK_TEST")
	}
	client, err := NewDefaultClient()
	require.NoError(t, err)
	err = client.Ping()
	require.NoError(t, err)
}
