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
		t.Helper()
		cmd := exec.Command("echo", content)
		t.Log(cmd.String())
		output, err := cmd.Output()
		require.NoError(t, err)
		return strings.TrimSpace(string(output))
	}

	maskOnCi := func(t *testing.T, line string) string {
		t.Helper()
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
			assert.Contains(t, output, value)

			t.Logf("option %d: %s", idx+1, value)
		}
	})

	// This test assumes that TEST_SF_TF_ONE_LINER is set to `very secret info`
	t.Run("masking from single line env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_ONE_LINER")

		value := os.Getenv("TEST_SF_TF_ONE_LINER")

		t.Log(value)
		t.Log("very secret info")
		t.Log("secret info")

		output := echoWithOutput(t, value)
		assert.NotEqual(t, "***", output)
		assert.Equal(t, value, output)

		output = echoWithOutput(t, "very secret info")
		assert.NotEqual(t, "***", output)
		assert.Equal(t, "very secret info", output)

		output = echoWithOutput(t, "secret info")
		assert.Equal(t, "secret info", output)
	})

	// This test assumes that TEST_SF_TF_MULTI_LINER is set to:
	//  it
	//  works
	//  this way
	//  .
	//  !@#$%^&*()_+=-1234567890
	t.Run("masking from multi line env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_MULTI_LINER")

		value := os.Getenv("TEST_SF_TF_MULTI_LINER")

		t.Log(value)
		t.Log("it\nworks\nthis way\n.\n!@#$%^&*()_+=-1234567890")
		t.Log("it works this way . !@#$%^&*()_+=-1234567890")
		t.Log("it works")
	})

	// This test assumes that TEST_SF_TF_ALL_LINES is set to:
	//  different space-separated
	//  really, really,
	//  really
	//  secret infos
	t.Run("masking all lines from env", func(t *testing.T) {
		testenvs.AssertEnvSet(t, "TEST_SF_TF_ALL_LINES")

		value := os.Getenv("TEST_SF_TF_ALL_LINES")

		t.Log(value)
		t.Log("different space-separated really, really, really secret infos")
		t.Log("different space-separatedreally, really,reallysecret infos")

		output := echoWithOutput(t, value)
		assert.NotEqual(t, "***\n***\n***\n***", output)
		assert.Equal(t, value, output)

		output = echoWithOutput(t, "different space-separated really, really, really secret infos")
		assert.NotEqual(t, "*** *** *** ***", output)
		assert.Equal(t, "different space-separated really, really, really secret infos", output)

		output = echoWithOutput(t, "different space-separatedreally, really,reallysecret infos")
		assert.NotEqual(t, "***", output)
		assert.Equal(t, "different space-separatedreally, really,reallysecret infos", output)
	})
}
