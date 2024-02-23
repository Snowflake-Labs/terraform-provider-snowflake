package testenvs_test

import (
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func Test_GetOrSkipTest(t *testing.T) {
	runGetOrSkipInGoroutineAndWait := func(tut *testing.T) string {
		var env string
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			env = testenvs.GetOrSkipTest(tut, testenvs.User)
		}()
		wg.Wait()
		return env
	}

	t.Run("skip test if missing", func(t *testing.T) {
		t.Setenv(string(testenvs.User), "")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWait(tut)

		require.True(t, tut.Skipped())
		require.Empty(t, env)
	})

	t.Run("get env if exists", func(t *testing.T) {
		t.Setenv(string(testenvs.User), "user")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWait(tut)

		require.False(t, tut.Skipped())
		require.Equal(t, "user", env)
	})
}
