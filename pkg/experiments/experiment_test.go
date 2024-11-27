package experiments

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_experiments(t *testing.T) {
	ghActionsValue := testenvs.GetOrSkipTest(t, testenvs.GithubActions)

	echoWithOutput := func(t *testing.T, content string) string {
		cmd := exec.Command("echo", content)
		t.Logf(cmd.String())
		output, err := cmd.Output()
		require.NoError(t, err)
		return strings.TrimSpace(string(output))
	}

	maskOnCi := func(t *testing.T, line string) string {
		if ghActionsValue == "true" {
			t.Logf("trying to mask using `%s`", line)
			return echoWithOutput(t, line)
		}
		return ""
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
			output := maskOnCi(t, fmt.Sprintf(option, value))
			assert.Contains(t, string(output), value)

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

		output := echoWithOutput(t, value)
		assert.Equal(t, "***", output)

		output = echoWithOutput(t, "very secret info")
		assert.Equal(t, "***", output)

		output = echoWithOutput(t, "secret info")
		assert.Equal(t, "secret info", output)
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

		output := echoWithOutput(t, value)
		assert.Equal(t, "***\n***\n***\n***", output)

		output = echoWithOutput(t, "different space-separated really, really, really secret infos")
		assert.Equal(t, "*** *** *** ***\n", output)

		output = echoWithOutput(t, "different space-separatedreally, really,reallysecret infos")
		assert.Equal(t, "************\n", output)
	})
}
