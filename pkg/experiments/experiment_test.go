package experiments

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func Test_experiments(t *testing.T) {

	echo := func(content string) error {
		cmd := exec.Command("echo", content)
		t.Logf(cmd.String())
		return cmd.Run()
	}

	maskOnCi := func(line string) error {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			t.Logf("masking `%s`", line)
			t.Setenv("TEST_SF_TF_MASKING_TEST", line)
			return echo(`"::add-mask::$TEST_SF_TF_MASKING_TEST"`)
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
		require.NoError(t, echo(b))
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

	t.Run("masking from multiple env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_MULTI_LINE")

		value := os.Getenv("TEST_SF_TF_MULTI_LINE")

		require.NoError(t, echo(value))
		t.Log(value)
		t.Log("different space-separated really, really, really secret infos")
	})
}
