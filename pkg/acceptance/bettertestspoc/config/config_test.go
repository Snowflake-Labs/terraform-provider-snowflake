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
	convertJsonToHclV1 := func(json string) (string, error) {
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

		return fmt.Sprintf("%s", strings.ReplaceAll(string(formatted[:]), "\n\n", "\n")), nil
	}

	examples := []string{
		`{
  "resource": {
    "aws_instance": {
      "example": {
        "instance_type": "t2.micro",
        "ami": "ami-abc123"
      }
    }
  }
}`,
		`{
  "resource": {
    "aws_instance": {
      "example": {
        "instance_type": "t2.micro",
        "ami": "${resource.name.attribute}"
      }
    }
  }
}`,
	}

	for _, example := range examples {
		example := example
		t.Run("test HCL v1", func(t *testing.T) {
			parsed, err := convertJsonToHclV1(example)
			require.NoError(t, err)

			fmt.Printf("%s", parsed)
		})
	}
}
