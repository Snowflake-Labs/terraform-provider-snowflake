package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientPing(t *testing.T) {
	client, err := NewDefaultClient()
	require.NoError(t, err)
	err = client.Ping()
	require.NoError(t, err)
}
