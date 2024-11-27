package experiments

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func Test_experiments(t *testing.T) {

	echo := func(content string) error {
		return exec.Command("echo", content).Run()
	}

	maskOnCi := func(line string) error {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			t.Log("masking")
			return echo(fmt.Sprintf("::add-mask::%s", line))
		}
		return nil
	}

	t.Run("dynamic masking", func(t *testing.T) {
		a := "something to mask"
		b := "something not to mask"
		c := "katarakta"

		err := maskOnCi(a)
		require.NoError(t, err)
		err = maskOnCi(c)
		require.NoError(t, err)

		require.NoError(t, echo(a))
		require.NoError(t, echo(a))
		require.NoError(t, echo(c))
		t.Log(a)
		t.Log(b)
		t.Log(c)
	})

	t.Run("masking from env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_ONE_LINER")

		value := os.Getenv("TEST_SF_TF_ONE_LINER")

		require.NoError(t, echo(value))
		t.Log(value)
	})
}
