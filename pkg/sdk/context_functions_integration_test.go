package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CurrentSession(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	session, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, session)
}
