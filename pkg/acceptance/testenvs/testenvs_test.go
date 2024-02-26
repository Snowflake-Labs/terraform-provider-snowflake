package testenvs_test

import (
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func Test_GetOrSkipTest(t *testing.T) {
	// runGetOrSkipInGoroutineAndWaitForCompletion is needed because underneath we test t.Skipf, that leads to t.SkipNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runGetOrSkipInGoroutineAndWaitForCompletion := func(t *testing.T) string {
		t.Helper()
		var env string
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			env = testenvs.GetOrSkipTest(t, testenvs.User)
		}()
		wg.Wait()
		return env
	}

	t.Run("skip test if missing", func(t *testing.T) {
		t.Setenv(string(testenvs.User), "")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWaitForCompletion(tut)

		require.True(t, tut.Skipped())
		require.Empty(t, env)
	})

	t.Run("get env if exists", func(t *testing.T) {
		t.Setenv(string(testenvs.User), "user")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWaitForCompletion(tut)

		require.False(t, tut.Skipped())
		require.Equal(t, "user", env)
	})
}

func Test_Assertions(t *testing.T) {
	// runAssertionInGoroutineAndWaitForCompletion is needed because underneath we test require, that leads to t.FailNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runAssertionInGoroutineAndWaitForCompletion := func(assertion func()) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			assertion()
		}()
		wg.Wait()
	}

	t.Run("test if env does not exist", func(t *testing.T) {
		t.Setenv(string(testenvs.User), "")

		tut1 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvNotSet(tut1, string(testenvs.User)) })

		tut2 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvSet(tut2, string(testenvs.User)) })

		require.False(t, tut1.Failed())
		require.True(t, tut2.Failed())
	})

	t.Run("test if env exists", func(t *testing.T) {
		t.Setenv(string(testenvs.User), "user")

		tut1 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvNotSet(tut1, string(testenvs.User)) })

		tut2 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvSet(tut2, string(testenvs.User)) })

		require.True(t, tut1.Failed())
		require.False(t, tut2.Failed())
	})
}
