package experiments_test

import (
	"bytes"
	"fmt"
	"testing"

	hclv1printer "github.com/hashicorp/hcl/hcl/printer"
	hclv1parser "github.com/hashicorp/hcl/json/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_exploreHclV1 shows conversion from .tf.json to .tf using HCL v1 lib.
// The main takes:
// - block types (e.g. resource) stay quoted (while they should be unquoted)
// - attribute names stay quoted (while they should be unquoted)
// - objects and object lists have equal sign (what should not be there)
// - references in depends_on stay quoted (while they should be unquoted)
// - there are two newlines between every attribute
//
// Because of the above, the current implementation using json generation and later conversion to HCL using hcl v1
// introduces formatters tackling the issues above. This is not the perfect solution but it's fast and will work in the meantime.
// Check: config.DefaultHclConfigProvider.
//
// References:
// - https://developer.hashicorp.com/terraform/language/syntax/json
// - https://github.com/hashicorp/hcl/blob/56a9aee5207dbaed7f061cd926b96fc159d26ea0/json/spec.md
// - https://developer.hashicorp.com/terraform/language/resources/syntax
// - https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on
// TODO [SNOW-1501905]: explore HCL v2 in more detail (especially struct tags generation; probably with migration to plugin framework because of schema models); ref: https://github.com/hashicorp/hcl/blob/bee2dc2e75f7528ad85777b7a013c13796426bd6/gohcl/encode_test.go#L48
func Test_exploreHclV1(t *testing.T) {
	convertJsonToHclStringV1 := func(json string) (string, error) {
		parsed, err := hclv1parser.Parse([]byte(json))
		if err != nil {
			return "", err
		}

		var buffer bytes.Buffer
		err = hclv1printer.Fprint(&buffer, parsed)
		if err != nil {
			return "", err
		}

		formatted, err := hclv1printer.Format(buffer.Bytes())
		if err != nil {
			return "", err
		}

		return string(formatted), nil
	}

	t.Run("test HCL v1", func(t *testing.T) {
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
                        "depends_on": [
                            "some_other_resource.some_name",
                            "other_resource.some_other_name",
							"data.some_datasource.some_fancy_datasource"
                        ]
                    }
                }
            }
        }`
		expectedResult := `"resource" "snowflake_share" "test" {
  "attribute_int" = 1

  "attribute_bool" = true

  "attribute_string" = "some string"

  "string_template" = "${resource.name.attribute}"

  "string_list" = ["a", "b", "a"]

  "object_list" = {
    "int_field" = 1

    "string_field" = "first item"
  }

  "object_list" = {
    "int_field" = 2

    "string_field" = "second item"
  }

  "single_object" = {
    "prop1" = 1

    "prop2" = "two"
  }

  "depends_on" = ["some_other_resource.some_name", "other_resource.some_other_name", "data.some_datasource.some_fancy_datasource"]
}
`

		result, err := convertJsonToHclStringV1(resourceJson)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		fmt.Printf("%s", result)
	})
}
