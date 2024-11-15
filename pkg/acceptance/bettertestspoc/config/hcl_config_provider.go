package config

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	hclv1printer "github.com/hashicorp/hcl/hcl/printer"
	hclv1parser "github.com/hashicorp/hcl/json/parser"
)

type HclProvider interface {
	HclFromJson(json []byte) (string, error)
}

type hclV1ConfigProvider struct{}

func NewHclV1ConfigProvider() HclProvider {
	return &hclV1ConfigProvider{}
}

func (h *hclV1ConfigProvider) HclFromJson(json []byte) (string, error) {
	hcl, err := convertJsonToHclStringV1(json)
	if err != nil {
		return "", err
	}

	hcl, err = formatResult(hcl)
	if err != nil {
		return "", err
	}
	hcl, err = fixDependsOn(hcl)
	if err != nil {
		return "", err
	}

	return hcl, nil
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

func formatResult(input string) (string, error) {
	return fmt.Sprintf("%s", strings.ReplaceAll(input, "\n\n", "\n")), nil
}

// Based on https://developer.hashicorp.com/terraform/language/syntax/json#depends_on should be processed in a special way, but it isn't.
func fixDependsOn(s string) (string, error) {
	dependsOnRegex := regexp.MustCompile(`("depends_on" = )(\["\w+\.\w+"(, "\w+\.\w+")*])`)
	submatches := dependsOnRegex.FindStringSubmatch(s)
	if len(submatches) < 2 {
		return s, nil
	} else {
		withoutQuotes := strings.ReplaceAll(submatches[2], `"`, "")
		return dependsOnRegex.ReplaceAllString(s, fmt.Sprintf(`$1%s`, withoutQuotes)), nil
	}
}
