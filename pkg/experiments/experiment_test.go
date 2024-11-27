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
	ghActionsValue := testenvs.GetOrSkipTest(t, testenvs.GithubActions)

	echoWithOutput := func(content string) ([]byte, error) {
		cmd := exec.Command("echo", content)
		t.Logf(cmd.String())
		return cmd.Output()
	}

	maskOnCi := func(line string) ([]byte, error) {
		if ghActionsValue == "true" {
			t.Logf("trying to mask using `%s`", line)
			return echoWithOutput(line)
		}
		return []byte{}, nil
	}

	// This test show that programmatically we are not able to mask using the GH workflow command.
	t.Run("dynamic masking", func(t *testing.T) {
		value := "something to mask"

		for idx, option := range []string{
			`::add-mask::%s`,
			`"::add-mask::%s"`,
			`'::add-mask::%s'`,
			`::add-mask::"%s"`,
			`::add-mask::'%s'`,
			`::add-mask:: %s`,
			`::add-mask:: "%s"`,
			`::add-mask:: '%s'`,
			`"::add-mask:: %s"`,
			`'::add-mask:: %s'`,
		} {
			output, err := maskOnCi(fmt.Sprintf(option, value))
			require.NoError(t, err)
			require.Equal(t, value, string(output))

			t.Logf("option %d: %s", idx+1, value)
		}
	})

	// This test assumes that TEST_SF_TF_ONE_LINER is set to `very secret info`
	t.Run("masking from env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_ONE_LINER")

		value := os.Getenv("TEST_SF_TF_ONE_LINER")

		t.Log(value)
		t.Log("very secret info")
		t.Log("secret info")

		output, err := echoWithOutput(value)
		require.NoError(t, err)
		require.Equal(t, "***", string(output))

		output, err = echoWithOutput("very secret info")
		require.NoError(t, err)
		require.Equal(t, "***", string(output))

		output, err = echoWithOutput("secret info")
		require.NoError(t, err)
		require.Equal(t, "secret info", string(output))
	})

	// This test assumes that TEST_SF_TF_MULTI_LINE is set to:
	//  different space-separated
	//  really, really,
	//  really
	//  secret infos
	t.Run("masking from multiple env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_MULTI_LINE")

		value := os.Getenv("TEST_SF_TF_MULTI_LINE")

		t.Log(value)
		t.Log("different space-separated really, really, really secret infos")
		t.Log("different space-separatedreally, really,reallysecret infos")

		output, err := echoWithOutput(value)
		require.NoError(t, err)
		require.Equal(t, "*** *** *** ***", string(output))

		output, err = echoWithOutput("different space-separated really, really, really secret infos")
		require.NoError(t, err)
		require.Equal(t, "*** *** *** ***", string(output))

		output, err = echoWithOutput("different space-separatedreally, really,reallysecret infos")
		require.NoError(t, err)
		require.Equal(t, "************", string(output))
	})
}
