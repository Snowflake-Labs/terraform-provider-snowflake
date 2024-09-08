package config

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	hclv1printer "github.com/hashicorp/hcl/hcl/printer"
	hclv1parser "github.com/hashicorp/hcl/json/parser"

	"github.com/stretchr/testify/require"
)

func Test_exploreHcl(t *testing.T) {

	// TODO: describe why V1 and not V2 is used
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

		return string(formatted[:]), nil
	}

	formatResult := func(input string) string {
		return fmt.Sprintf("%s", strings.ReplaceAll(input, "\n\n", "\n"))
	}

	examples := []string{
		`{
  "resource": {
    "some_resource": {
      "name": {
        "attribute_1": "123",
        "attribute_2": 1
      }
    }
  }
}`,
		`{
  "resource": {
    "some_resource": {
      "name": {
        "attribute_1": "some value",
        "attribute_2": "${resource.name.attribute}"
      }
    }
  }
}`,
	}

	for _, example := range examples {
		example := example
		t.Run("test HCL v1", func(t *testing.T) {
			result, err := convertJsonToHclStringV1(example)
			require.NoError(t, err)

			fmt.Printf("%s", result)
			fmt.Printf("%s", formatResult(result))
		})
	}
}
