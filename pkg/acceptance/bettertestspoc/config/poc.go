package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	hclv1printer "github.com/hashicorp/hcl/hcl/printer"
	hclv1parser "github.com/hashicorp/hcl/json/parser"

	"github.com/stretchr/testify/require"
)

type ResourceJson struct {
	Resource map[string]map[string]ResourceModel `json:"resource"`
}

func ResourceJsonFrom(model ResourceModel) ResourceJson {
	return ResourceJson{
		Resource: map[string]map[string]ResourceModel{
			fmt.Sprintf("%s", model.Resource()): {
				fmt.Sprintf("%s", model.ResourceName()): model,
			},
		},
	}
}

func FromModelPoc(t *testing.T, model ResourceModel) string {
	t.Helper()

	modelJson := ResourceJsonFrom(model)

	b, err := json.MarshalIndent(modelJson, "", "    ")
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", b)

	s, err := convertJsonToHclStringV1(b)
	require.NoError(t, err)

	formatResult := func(input string) string {
		return fmt.Sprintf("%s", strings.ReplaceAll(input, "\n\n", "\n"))
	}
	formatted := formatResult(s)

	// Based on https://developer.hashicorp.com/terraform/language/syntax/json#depends_on should be processed in a special way, but it isn't.
	fixDependsOn := func(s string) string {
		t.Log("Fixing depends_on in the generated config")
		dependsOnRegex := regexp.MustCompile(`("depends_on" = )(\["\w+\.\w+"(, "\w+\.\w+")*])`)
		submatches := dependsOnRegex.FindStringSubmatch(s)
		if len(submatches) < 2 {
			t.Log("No depends_on found, returning the input unchanged")
			return s
		} else {
			t.Logf("Submatches found:\n%s", submatches)
			withoutQuotes := strings.ReplaceAll(submatches[2], `"`, "")

			return dependsOnRegex.ReplaceAllString(s, fmt.Sprintf(`$1%s`, withoutQuotes))
		}
	}
	final := fixDependsOn(formatted)

	t.Logf("Generated config:\n%s", final)
	return final
}

func convertJsonToHclStringV1(jsonBytes []byte) (string, error) {
	parsed, err := hclv1parser.Parse(jsonBytes)
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
