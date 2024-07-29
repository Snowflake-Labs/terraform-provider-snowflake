package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func Test_Encoding_And_Parsing_Of_ResourceIdentifier(t *testing.T) {
	testCases := []struct {
		Input                 []string
		Expected              string
		ExpectedAfterDecoding []string
	}{
		{Input: []string{sdk.NewSchemaObjectIdentifier("a", "b", "c").FullyQualifiedName(), "info"}, Expected: `"a"."b"."c"|info`},
		{Input: []string{}, Expected: ``},
		{Input: []string{"", "", ""}, Expected: `||`},
		{Input: []string{"a", "b", "c"}, Expected: `a|b|c`},
		// If one of the parts contains a separator sign (pipe in this case),
		// we can end up with more parts than we started with.
		{Input: []string{"a", "b", "c|d"}, Expected: `a|b|c|d`, ExpectedAfterDecoding: []string{"a", "b", "c", "d"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Encoding and parsing %s resource identifier`, testCase.Input), func(t *testing.T) {
			encodedIdentifier := EncodeResourceIdentifier(testCase.Input...)
			assert.Equal(t, testCase.Expected, encodedIdentifier)

			parsedIdentifier := ParseResourceIdentifier(encodedIdentifier)
			if testCase.ExpectedAfterDecoding != nil {
				assert.Equal(t, testCase.ExpectedAfterDecoding, parsedIdentifier)
			} else {
				assert.Equal(t, testCase.Input, parsedIdentifier)
			}
		})
	}
}
