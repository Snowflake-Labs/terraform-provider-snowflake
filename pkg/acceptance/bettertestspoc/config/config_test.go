package config

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/json/parser"
	"github.com/stretchr/testify/require"
)

func Test_hclV1(t *testing.T) {

	convertJsonToHcl := func(json string) (string, error) {
		parsed, err := parser.Parse([]byte(json))
		if err != nil {
			return "", err
		}

		var buffer bytes.Buffer
		err = printer.Fprint(&buffer, parsed)
		if err != nil {
			return "", err
		}

		formatted, err := printer.Format(buffer.Bytes())
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s", strings.ReplaceAll(string(formatted[:]), "\n\n", "\n")), nil
	}

	t.Run("basic example", func(t *testing.T) {
		example := `
{
  "resource": {
    "aws_instance": {
      "example": {
        "instance_type": "t2.micro",
        "ami": "ami-abc123"
      }
    }
  }
}
`

		parsed, err := convertJsonToHcl(example)
		require.NoError(t, err)

		fmt.Printf("%s", parsed)
	})
}
