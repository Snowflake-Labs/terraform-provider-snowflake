package config_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HclProvider(t *testing.T) {
	// This test shows that all issues from the experiments package were resolved using formatters.
	t.Run("test default hcl provider", func(t *testing.T) {
		resourceJson := `{
            "resource": {
                "snowflake_share": {
                    "test": {
                        "attribute_int": 1,
                        "attribute_bool": true,
                        "attribute_string": "some string",
                        "string_template": "${resource.name.attribute}",
                        "string_list": ["a", "b", "a"],
                        "object_list": [
                            {
                                "int_field": 1,
                                "string_field": "first item"
                            },
                            {
                                "int_field": 2,
                                "string_field": "second item"
                            }
                        ],
                        "single_object": {
                            "prop1": 1,
                            "prop2": "two"
                        },
						"multiline_string": "SF_TF_TEST_MULTILINE_PLACEHOLDER_PRIVATE_KEY-----BEGIN PRIVATE KEY-----\nabc\ndef\nghj\n-----END PRIVATE KEY-----\nSF_TF_TEST_MULTILINE_PLACEHOLDER_PRIVATE_KEY",
                        "depends_on": [
                            "some_other_resource.some_name",
                            "other_resource.some_other_name",
							"data.some_datasource.some_fancy_datasource"
                        ]
                    }
                }
            }
        }`
		expectedResult := `resource "snowflake_share" "test" {
  attribute_int = 1
  attribute_bool = true
  attribute_string = "some string"
  string_template = "${resource.name.attribute}"
  string_list = ["a", "b", "a"]
  object_list {
    int_field = 1
    string_field = "first item"
  }
  object_list {
    int_field = 2
    string_field = "second item"
  }
  single_object {
    prop1 = 1
    prop2 = "two"
  }
  multiline_string = <<EOT
-----BEGIN PRIVATE KEY-----
abc
def
ghj
-----END PRIVATE KEY-----
EOT
  depends_on = [some_other_resource.some_name, other_resource.some_other_name, data.some_datasource.some_fancy_datasource]
}
`

		result, err := config.DefaultHclConfigProvider.HclFromJson([]byte(resourceJson))
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		fmt.Printf("%s", result)
	})
}
